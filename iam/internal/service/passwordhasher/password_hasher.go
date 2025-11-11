package passwordhasher

import (
	"golang.org/x/crypto/bcrypt"

	def "github.com/pinai4/spaceship-factory/iam/internal/service"
)

var _ def.PasswordHasher = (*passwordHasher)(nil)

type passwordHasher struct {
	cost int
}

// NewPasswordHasher creates a new instance of passwordHasher with the given cost factor.
// Recommended bcrypt cost is between 10 and 14 for production.
func NewPasswordHasher(cost int) *passwordHasher {
	if cost < bcrypt.MinCost || cost > bcrypt.MaxCost {
		cost = bcrypt.DefaultCost
	}
	return &passwordHasher{cost: cost}
}

func (b *passwordHasher) Hash(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), b.cost)
	if err != nil {
		return "", err
	}
	return string(hashed), nil
}

// Compare compares a bcrypt-hashed password with its possible plain-text equivalent.
// Returns nil if they match, or an error otherwise.
func (b *passwordHasher) Compare(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
