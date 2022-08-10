package dal

import "github.com/maxiancillotti/access-control/internal/domain"

// Methods implementing this interface must return error ErrorDataAccessEmptyResult
// when appropriate. See on "authservices.go".
type AdminsDAL interface {
	// Inserts an Admin and returns the created ID
	Insert(*domain.Admin) (uint, error)

	// Updates a password for the given admin
	UpdatePassword(*domain.Admin) error

	// Updates enabled state for the given admin
	UpdateEnabledState(*domain.Admin) error

	// Deletes an Admin logically by its ID
	Delete(uint) error

	// Retrieves an Admin by its ID
	Select(uint) (*domain.Admin, error)

	// Retrieves an Admin by its Username
	SelectByUsername(string) (*domain.Admin, error)

	// Checks if a UserID exists
	Exists(uint) (bool, error)

	// Checks if a username exists even if it was soft-deleted
	UsernameExists(string) (bool, error)
}
