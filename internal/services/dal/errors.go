package dal

import "errors"

var (

	// Errors received from data access layer.
	// Apply DIP applied to this err value:
	// DAL should return this err directly, not wrapped, when detecting that
	// an internal error of the data storage interface has this meaning.
	// Then, the service can detect this particular error without depending
	// on an err value originated at the DAL.
	// if errors.Cause(err).Error() == ErrorDataAccessEmptyResult.Error() {
	// USAGE OF ERRORS.IS:
	// When using custom error type, errors.Is doesn't work, always returns false.
	// More info: https://stackoverflow.com/questions/62441960/error-wrap-unwrap-type-checking-with-errors-is
	// So don't wrap this error.
	ErrorDataAccessEmptyResult = errors.New("no result found with the given parameters")
)
