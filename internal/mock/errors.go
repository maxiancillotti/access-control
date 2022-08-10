package mock

import "github.com/pkg/errors"

var (
	errMockInvalidInputID     = errors.New("invalid_input")
	errMockEmptyResult        = errors.New("empty_result")
	errMockInternal           = errors.New("internal_error")
	errMockInvalidCredentials = errors.New("invalid_credentials_error")
	errMockInvalidToken       = errors.New("invalid_token_error")
	errMockPermissions        = errors.New("permissions_error")
	errMockUnprocessable      = errors.New("unprocessable_error")
)
