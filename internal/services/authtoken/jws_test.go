package authtoken

import (
	"testing"

	"github.com/maxiancillotti/access-control/internal/domain"

	"github.com/stretchr/testify/assert"
)

func TestVerifyToken(t *testing.T) {

	t.Run("Success", func(t *testing.T) {

		// 10 years of validity.
		token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE5NTkyNzc0NzEsImlzcyI6IkF1dGhUb2tlbkFQSSIsInBlcm1pc3Npb25zIjp7InJlc291cmNlIjpbInBlcm1pc3Npb24iXX0sInN1YiI6InVzZXJTdWNjZXNzIn0.gCSXyLP8eInotFu7WvYg-Y3Gn4piNtR0LhYPgYm6RBI"

		jwtToken, err := testAuthTokenServices.jws.verifyToken([]byte(token), []byte(testAuthTokenServices.config.builder.jwtSigningSecretKey))
		assert.Nil(t, err)
		assert.NotNil(t, jwtToken)
	})

	t.Run("Incorrect secretKey error", func(t *testing.T) {

		// 10 years of validity.
		token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE5NTkyNzc0NzEsImlzcyI6IkF1dGhUb2tlbkFQSSIsInBlcm1pc3Npb25zIjp7InJlc291cmNlIjpbInBlcm1pc3Npb24iXX0sInN1YiI6InVzZXJTdWNjZXNzIn0.gCSXyLP8eInotFu7WvYg-Y3Gn4piNtR0LhYPgYm6RBI"

		jwtToken, err := testAuthTokenServices.jws.verifyToken([]byte(token), []byte("this is an incorrect secret key"))
		assert.NotNil(t, err)
		assert.Nil(t, jwtToken)

		// Probable error: failed to verify jws signature: failed to verify message: failed to match hmac signature
		// Cannot test the err value because that way we would depend on the jwt library.
		t.Log("err value:", err)
	})

	t.Run("Invalid token error", func(t *testing.T) {

		// Same token as before but final character "1" added
		token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE5NTkyNzc0NzEsImlzcyI6IkF1dGhUb2tlbkFQSSIsInBlcm1pc3Npb25zIjp7InJlc291cmNlIjpbInBlcm1pc3Npb24iXX0sInN1YiI6InVzZXJTdWNjZXNzIn0.gCSXyLP8eInotFu7WvYg-Y3Gn4piNtR0LhYPgYm6RBI1"

		jwtToken, err := testAuthTokenServices.jws.verifyToken([]byte(token), []byte("this is an incorrect secret key"))
		assert.NotNil(t, err)
		assert.Nil(t, jwtToken)

		// Probable error: failed to verify jws signature: failed to verify message: failed to match hmac signature
		// Cannot test the err value because that way we would depend on the jwt library.
		t.Log("err value:", err)

	})

	t.Run("Token expired error", func(t *testing.T) {

		// Expired token
		token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2NDMzMDk5NzUsImlzcyI6IkF1dGhUb2tlbkFQSSIsInBlcm1pc3Npb25zIjp7Ii9jdXN0b21lcnMiOlsiR0VUIiwiUE9TVCJdfSwic3ViIjoiQVBJVXNlciJ9.waVrdOzWk46tNu9S-ObPyBY4d6OFLr6uqLBz53sv9uc"

		jwtToken, err := testAuthTokenServices.jws.verifyToken([]byte(token), []byte(testAuthTokenServices.config.builder.jwtSigningSecretKey))
		assert.NotNil(t, err)
		assert.Nil(t, jwtToken)

		// Probable error: exp not satisfied
		// Cannot test the err value because that way we would depend on the jwt library.
		t.Log("err value:", err)
	})
}

func TestSignToken(t *testing.T) {

	t.Run("Success", func(t *testing.T) {
		// Initialization
		var userID uint = 1

		//usrPerm := make(domain.UserPermissions)
		//usrPerm["resource"] = append(usrPerm["resource"], "permission")

		usrPerm := domain.UserPermissions{
			domain.RESTpermissionCategory: domain.RESTPermissionsPathsMethods{
				"/resource": {"POST"},
			},
		}

		token, err := testAuthTokenServices.getPayload(userID, usrPerm)

		if err != nil || token == nil {
			t.Fatal("cannot initialize payload for the signing call")
		}

		// Call
		signedToken, err := testAuthTokenServices.jws.signToken(token, []byte(testAuthTokenServicesConfig.builder.jwtSigningSecretKey))

		// TEST

		// Unit
		assert.Nil(t, err)
		assert.NotNil(t, signedToken)

		//t.Log("Signed token value:", signedToken)

		// Integration
		jwtToken, err := testAuthTokenServices.jws.verifyToken(signedToken, []byte(testAuthTokenServices.config.builder.jwtSigningSecretKey))
		assert.Nil(t, err)
		assert.NotNil(t, jwtToken)
	})

	// Error case cannot be generated because it would be from the external
	// library and for reasons that cannot be reproduced.
	// Also really unlikely.
}
