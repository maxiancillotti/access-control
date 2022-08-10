package mock

type ErrorCheckerMock struct{}

func (c *ErrorCheckerMock) ErrorIsInvalidInputIdentifier(err error) bool {
	return err.Error() == errMockInvalidInputID.Error()
}

func (c *ErrorCheckerMock) ErrorIsEmptyResult(err error) bool {
	return err.Error() == errMockEmptyResult.Error()
}

// Token creation error
func (c *ErrorCheckerMock) ErrorIsInternal(err error) bool {
	return err.Error() == errMockInternal.Error()
}

// Token creation error
func (c *ErrorCheckerMock) ErrorIsInvalidCredentials(err error) bool {
	return err.Error() == errMockInvalidCredentials.Error()
}

// Token validation error
func (c *ErrorCheckerMock) ErrorIsInvalidToken(err error) bool {
	return err.Error() == errMockInvalidToken.Error()
}

// Token validation error
func (c *ErrorCheckerMock) ErrorIsNotEnoughPermissions(err error) bool {
	return err.Error() == errMockPermissions.Error()
}

// Token validation error
func (c *ErrorCheckerMock) ErrorIsSemanticallyUnprocesable(err error) bool {
	return err.Error() == errMockUnprocessable.Error()
}
