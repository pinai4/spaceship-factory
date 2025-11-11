package service

// PasswordHasher defines methods for hashing and validating passwords.
type PasswordHasher interface {
	// Hash takes a plain-text password and returns its bcrypt hash.
	Hash(password string) (string, error)

	// Compare checks if the provided plain-text password matches the hashed password.
	Compare(hashedPassword, password string) error
}
