//go:build integration

package e2e

import (
	"github.com/pinai4/spaceship-factory/platform/pkg/logger"
)

func (s *E2ESuite) TearDownSuite() {
	s.cleanup()
}

func (s *E2ESuite) cleanup() {
	s.cleanupInventoryV1Client()
	s.cleanupTestEnvironment()
	s.suiteCtxCancel()
}

// cleanupTestEnvironment — releases all test environment resources
func (s *E2ESuite) cleanupTestEnvironment() {
	if s.env == nil {
		return
	}

	s.log.Info(s.suiteCtx, "🧹 Cleaning up the test environment...")

	if s.env.App != nil {
		if err := s.env.App.Terminate(s.suiteCtx); err != nil {
			s.log.Error(s.suiteCtx, "failed to stop application container", logger.Error(err))
		} else {
			s.log.Info(s.suiteCtx, "🛑 Application container stopped")
		}
	}

	if s.env.Auth != nil {
		if err := s.env.Auth.Terminate(s.suiteCtx); err != nil {
			s.log.Error(s.suiteCtx, "failed to stop Auth(mock) container", logger.Error(err))
		} else {
			s.log.Info(s.suiteCtx, "🛑 Auth(mock) container stopped")
		}
	}

	if s.env.Mongo != nil {
		if err := s.env.Mongo.Terminate(s.suiteCtx); err != nil {
			s.log.Error(s.suiteCtx, "failed to stop MongoDB container", logger.Error(err))
		} else {
			s.log.Info(s.suiteCtx, "🛑 MongoDB container stopped")
		}
	}

	if s.env.Network != nil {
		if err := s.env.Network.Remove(s.suiteCtx); err != nil {
			s.log.Error(s.suiteCtx, "failed to remove network", logger.Error(err))
		} else {
			s.log.Info(s.suiteCtx, "🛑 Network removed")
		}
	}

	s.log.Info(s.suiteCtx, "✅ Test environment successfully cleaned up")
}

func (s *E2ESuite) cleanupInventoryV1Client() {
	if s.grpcClientConn != nil {
		_ = s.grpcClientConn.Close()
	}
	s.log.Info(s.suiteCtx, "🛑 grpcClientConn has been closed")
}
