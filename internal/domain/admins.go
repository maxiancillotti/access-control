package domain

type Admin struct {
	ID           uint
	Username     string
	PasswordHash string // Salted and hashed password
	PasswordSalt string
	Enabled      bool
}
