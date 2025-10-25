package closer

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/pinai4/spaceship-factory/platform/pkg/logger"
)

// shutdownTimeout timeout that is used in handleSignals method
const shutdownTimeout = 5 * time.Second

type Logger interface {
	Info(ctx context.Context, msg string, fields ...logger.Field)
	Error(ctx context.Context, msg string, fields ...logger.Field)
}

// Closer manages the application's graceful shutdown process
type Closer struct {
	mu     sync.Mutex                    // Protection against race conditions when adding functions
	once   sync.Once                     // Ensures CloseAll is called only once
	done   chan struct{}                 // Channel for completion notification
	funcs  []func(context.Context) error // Registered shutdown functions
	logger Logger                        // Logger in use
}

// New creates a new instance of Closer
// If signals are passed, Closer starts listening to them and calls CloseAll upon receiving one.
func New(logger Logger, signals ...os.Signal) *Closer {
	c := &Closer{
		done:   make(chan struct{}),
		logger: logger,
	}

	if len(signals) > 0 {
		go c.handleSignals(signals...)
	}

	return c
}

// handleSignals processes system signals and calls CloseAll with a fresh shutdown context
func (c *Closer) handleSignals(signals ...os.Signal) {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, signals...)
	defer signal.Stop(ch)

	select {
	case <-ch:
		c.logger.Info(context.Background(), "🛑 System signal received, starting graceful shutdown...")

		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), shutdownTimeout)
		defer shutdownCancel()

		if err := c.CloseAll(shutdownCtx); err != nil {
			c.logger.Error(context.Background(), "❌ Error while closing resources: %v", logger.Error(err))
		}

	case <-c.done:
		// CloseAll has already been called manually, just exit
	}
}

// AddNamed adds a shutdown function with a dependency name for logging
func (c *Closer) AddNamed(name string, f func(context.Context) error) {
	c.Add(func(ctx context.Context) error {
		start := time.Now()
		c.logger.Info(ctx, fmt.Sprintf("🧩 Closing %s...", name))

		err := f(ctx)

		duration := time.Since(start)
		if err != nil {
			c.logger.Error(ctx, fmt.Sprintf("❌ Error closing %s: %v (took %s)", name, err, duration))
		} else {
			c.logger.Info(ctx, fmt.Sprintf("✅ %s closed successfully in %s", name, duration))
		}
		return err
	})
}

// Add adds one or more shutdown functions
func (c *Closer) Add(f ...func(context.Context) error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.funcs = append(c.funcs, f...)
}

// CloseAll calls all registered shutdown functions.
// Returns the first occurred error, if any.
func (c *Closer) CloseAll(ctx context.Context) error {
	var result error

	c.once.Do(func() {
		defer close(c.done)

		c.mu.Lock()
		funcs := c.funcs
		c.funcs = nil // free memory
		c.mu.Unlock()

		if len(funcs) == 0 {
			c.logger.Info(ctx, "ℹ️ No functions to close.")
			return
		}

		c.logger.Info(ctx, "🚦 Starting graceful shutdown process...")

		errCh := make(chan error, len(funcs))
		var wg sync.WaitGroup

		// Execute in reverse order of addition
		for i := len(funcs) - 1; i >= 0; i-- {
			f := funcs[i]
			wg.Add(1)
			go func(f func(context.Context) error) {
				defer wg.Done()

				// Panic protection
				defer func() {
					if r := recover(); r != nil {
						errCh <- errors.New("panic recovered in closer")
						c.logger.Error(ctx, "⚠️ Panic in shutdown function", logger.String("panic", fmt.Sprintf("%v", r)))
					}
				}()

				if err := f(ctx); err != nil {
					errCh <- err
				}
			}(f)
		}

		// Close the error channel once all functions finish
		go func() {
			wg.Wait()
			close(errCh)
		}()

		// Read errors or context cancellation
		for {
			select {
			case <-ctx.Done():
				c.logger.Info(ctx, "⚠️ Context canceled during shutdown", logger.Error(ctx.Err()))
				if result == nil {
					result = ctx.Err()
				}
				return
			case err, ok := <-errCh:
				if !ok {
					c.logger.Info(ctx, "✅ All resources closed successfully")
					return
				}
				c.logger.Error(ctx, "❌ Error while closing", logger.Error(err))
				if result == nil {
					result = err
				}
			}
		}
	})

	return result
}
