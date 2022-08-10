package services

// See the method comments to know which AuthService method can return each error.
type ServiceErrorChecker interface {

	// CRUD error
	// Can return the full error message safely. Does not containg private info.
	ErrorIsInvalidInputIdentifier(err error) bool

	// Retrieval error.
	// Can return the full error message safely. Does not containg private info.
	ErrorIsEmptyResult(err error) bool

	// Authenticate error.
	// CreateToken error.
	// CRUD error.
	ErrorIsInternal(err error) bool

	// Authenticate error.
	ErrorIsInvalidCredentials(err error) bool

	// Authorize error.
	ErrorIsInvalidToken(err error) bool

	// Authorize error.
	ErrorIsNotEnoughPermissions(err error) bool

	// Authorize error.
	ErrorIsSemanticallyUnprocesable(err error) bool
}
