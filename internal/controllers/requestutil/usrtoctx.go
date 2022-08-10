package requestutil

import "context"

// AssignAuthenticatedUserIDToCxt
func AssignAuthenticatedUserIDToCxt(ctx context.Context, userID uint) context.Context {

	ctx = context.WithValue(ctx, contextKeyAuthBasicUserID, userID)
	return ctx
}

// GetAuthenticatedUserIDFromCtx
func GetAuthenticatedUserIDFromCtx(ctx context.Context) uint {

	authBasicUserIDCtx := ctx.Value(contextKeyAuthBasicUserID)
	authBasicUserID, ok := authBasicUserIDCtx.(uint)
	if !ok {
		return 0
	}
	return authBasicUserID
}

/*
// AssignAuthBasicCredentialsToCxt
func AssignAuthBasicCredentialsToCxt(ctx context.Context, username, password string) context.Context {

	ctx = context.WithValue(ctx, contextKeyAuthBasicUsername, username)
	ctx = context.WithValue(ctx, contextKeyAuthBasicPassword, password)
	return ctx
}

// GetAuthBasicCredentialsFromCtx
func GetAuthBasicCredentialsFromCtx(ctx context.Context) (*dto.UserCredentials, error) {

	authBasicUsernameCtx := ctx.Value(contextKeyAuthBasicUsername)
	authBasicUsername, ok := authBasicUsernameCtx.(string)
	if !ok {
		return nil, errors.New("error asserting authorization basic user")
	}

	authBasicPasswordCtx := ctx.Value(contextKeyAuthBasicPassword)
	authBasicPassword, ok := authBasicPasswordCtx.(string)
	if !ok {
		return nil, errors.New("error asserting authorization basic password")
	}

	return &dto.UserCredentials{
			Username: authBasicUsername,
			Password: authBasicPassword,
		},
		nil
}
*/
