package dal

import "github.com/maxiancillotti/access-control/internal/domain"

// Methods implementing this interface must return error ErrorDataAccessEmptyResult
// when appropriate. See on "authservices.go".
type UserDAL interface {
	// Inserts a User and returns the created ID
	Insert(*domain.User) (uint, error)

	// Updates a password for the given user
	UpdatePassword(*domain.User) error

	// Updates enabled state for the given user
	UpdateEnabledState(*domain.User) error

	// Deletes a User logically by its ID
	Delete(uint) error

	// Retrieves a User by its ID
	Select(uint) (*domain.User, error)

	// Retrieves a User by its Username
	SelectByUsername(string) (*domain.User, error)

	// Checks if a UserID exists
	Exists(uint) (bool, error)

	// Checks if a username exists even if it was soft-deleted
	UsernameExists(string) (bool, error)
}
