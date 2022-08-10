package requestutil

// ContextKey is used for context.Context value. The value requires a key that is not primitive type.
type authContextKey string

const (
	contextKeyRequestID         authContextKey = "requestID"
	contextKeyRequestBody       authContextKey = "requestBody"
	contextKeyAuthBasicUserID   authContextKey = "authBasic_UserID"
	contextKeyAuthBasicUsername authContextKey = "authBasic_User"
	contextKeyAuthBasicPassword authContextKey = "authBasic_Password"
)
