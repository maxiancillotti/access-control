package authtoken

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDecryptToken(t *testing.T) {

	t.Run("Success", func(t *testing.T) {

		// 10 years of validity.
		token := "eyJhbGciOiJkaXIiLCJlbmMiOiJBMTI4Q0JDLUhTMjU2In0..6QLqr6pXAP-gfclqIal4pw.S7R6RLB6qhH8u2N4RPgRTm0j4GccFVcQJKCX6Q9FaGXdVPRzjCe5Ndbl-vvM_AC-0S29DTbLUshqG15r-9LKLL8KropVQUDQaNiVTWwPOGOEWyRZ4GMLgvRcbCI9HrDtdN_QBq0ZHkheQmbCW65ec1HtawQCUS-Db9xRLNgjC4deR_6zVDSOxZb32ShrMfssOjvcvRspz5zt_5aRLiALlw1TRN6ytLO2mBr2elTvEdtiq-ftt-5ICNlQbhlyJvlKcp-E0nf6HVusaJAKV7cSXWYS_cyyCgHFE_XO-_Pzcgs.OUNiW0sMaTRSppXEHr_t2A"

		decryptedToken, err := testAuthTokenServices.jwe.decryptToken([]byte(token), []byte(testAuthTokenServices.config.builder.jwtEncryptionSecretKey))
		assert.Nil(t, err)
		assert.NotNil(t, decryptedToken)

		t.Log("Decrypted Token:", string(decryptedToken))
	})

	t.Run("Success with expired token", func(t *testing.T) {

		// Expired token
		token := "eyJhbGciOiJkaXIiLCJlbmMiOiJBMTI4Q0JDLUhTMjU2In0..s46rXL3R8ceNgG4r0lcMzw.7aVyoFReUBMt7tm6ZiPVejP7-JQrcVjZ09N_8_arnjEmXzgnLic2KZ9CWU6OZbMNxqvDvT-8iFOT7RBOGIU5totPdRJcPqCDLo0lVdk0_FM_CAw2iylrcs9xNMXgKRn3-ZeFS0GhS9KQLB2z533AmqVhZUOPPTPCaBgzM4lqy26uq-4qyypPSLnxp2s6Vj7er5IpeX6GREmixYZV0r4-t9YgvmLwMCRBKYtprAmDIxW7ONm55ls1n3fCOqVftsHjUNwzUF0ZtWxBzFvuaqqhrPPOsoqQilgAFZfatOqW2IM.S39VyLZa9GQ8bbAPHoDmug"

		decryptedToken, err := testAuthTokenServices.jwe.decryptToken([]byte(token), []byte(testAuthTokenServices.config.builder.jwtEncryptionSecretKey))
		assert.Nil(t, err)
		assert.NotNil(t, decryptedToken)

		t.Log("Decrypted Token:", string(decryptedToken))
	})

	t.Run("Error invalid token", func(t *testing.T) {

		// Same token as the first one but final character chaged for "1"
		token := "eyJhbGciOiJkaXIiLCJlbmMiOiJBMTI4Q0JDLUhTMjU2In0..6QLqr6pXAP-gfclqIal4pw.S7R6RLB6qhH8u2N4RPgRTm0j4GccFVcQJKCX6Q9FaGXdVPRzjCe5Ndbl-vvM_AC-0S29DTbLUshqG15r-9LKLL8KropVQUDQaNiVTWwPOGOEWyRZ4GMLgvRcbCI9HrDtdN_QBq0ZHkheQmbCW65ec1HtawQCUS-Db9xRLNgjC4deR_6zVDSOxZb32ShrMfssOjvcvRspz5zt_5aRLiALlw1TRN6ytLO2mBr2elTvEdtiq-ftt-5ICNlQbhlyJvlKcp-E0nf6HVusaJAKV7cSXWYS_cyyCgHFE_XO-_Pzcgs.OUNiW0sMaTRSppXEHr_t21"

		decryptedToken, err := testAuthTokenServices.jwe.decryptToken([]byte(token), []byte(testAuthTokenServices.config.builder.jwtEncryptionSecretKey))
		assert.NotNil(t, err)
		assert.Nil(t, decryptedToken)

		// Probable error: failed to decrypt message: failed to find matching recipient to decrypt key (last error = failed to decrypt: failed to decrypt payload: aead.Open failed: invalid ciphertext (tag mismatch))
		// Cannot test the err value because that way we would depend on the jwt library.
		t.Log("err value:", err)
	})

	t.Run("Error invalid secretKey size", func(t *testing.T) {

		// Same token as the first one
		token := "eyJhbGciOiJkaXIiLCJlbmMiOiJBMTI4Q0JDLUhTMjU2In0..6QLqr6pXAP-gfclqIal4pw.S7R6RLB6qhH8u2N4RPgRTm0j4GccFVcQJKCX6Q9FaGXdVPRzjCe5Ndbl-vvM_AC-0S29DTbLUshqG15r-9LKLL8KropVQUDQaNiVTWwPOGOEWyRZ4GMLgvRcbCI9HrDtdN_QBq0ZHkheQmbCW65ec1HtawQCUS-Db9xRLNgjC4deR_6zVDSOxZb32ShrMfssOjvcvRspz5zt_5aRLiALlw1TRN6ytLO2mBr2elTvEdtiq-ftt-5ICNlQbhlyJvlKcp-E0nf6HVusaJAKV7cSXWYS_cyyCgHFE_XO-_Pzcgs.OUNiW0sMaTRSppXEHr_t2A"
		secretKey := "invalid size"

		decryptedToken, err := testAuthTokenServices.jwe.decryptToken([]byte(token), []byte(secretKey))
		assert.NotNil(t, err)
		assert.Nil(t, decryptedToken)

		// Probable error: failed to decrypt message: failed to find matching recipient to decrypt key (last error = failed to decrypt: failed to decrypt payload: failed to fetch AEAD data: cipher: failed to create AES cipher for CBC: failed to execute block cipher function: crypto/aes: invalid key size 15)
		// Cannot test the err value because that way we would depend on the jwt library.
		t.Log("err value:", err)
	})

	t.Run("Error invalid secret key value", func(t *testing.T) {

		// Same token as the first one
		token := "eyJhbGciOiJkaXIiLCJlbmMiOiJBMTI4Q0JDLUhTMjU2In0..6QLqr6pXAP-gfclqIal4pw.S7R6RLB6qhH8u2N4RPgRTm0j4GccFVcQJKCX6Q9FaGXdVPRzjCe5Ndbl-vvM_AC-0S29DTbLUshqG15r-9LKLL8KropVQUDQaNiVTWwPOGOEWyRZ4GMLgvRcbCI9HrDtdN_QBq0ZHkheQmbCW65ec1HtawQCUS-Db9xRLNgjC4deR_6zVDSOxZb32ShrMfssOjvcvRspz5zt_5aRLiALlw1TRN6ytLO2mBr2elTvEdtiq-ftt-5ICNlQbhlyJvlKcp-E0nf6HVusaJAKV7cSXWYS_cyyCgHFE_XO-_Pzcgs.OUNiW0sMaTRSppXEHr_t2A"

		// Original key for testing but last character changed, so the lenght is ok
		secretKey := "%C*F-JaNdRgUkXp2s5v8x/A?D(G+KbP1"

		decryptedToken, err := testAuthTokenServices.jwe.decryptToken([]byte(token), []byte(secretKey))
		assert.NotNil(t, err)
		assert.Nil(t, decryptedToken)

		// Probable error: failed to decrypt message: failed to find matching recipient to decrypt key (last error = failed to decrypt: failed to decrypt payload: aead.Open failed: failed to generate plaintext from decrypted blocks: invalid padding)
		// Cannot test the err value because that way we would depend on the jwt library.
		t.Log("err value:", err)
	})
}

func TestEncryptToken(t *testing.T) {

	t.Run("Success", func(t *testing.T) {

		signedToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE5NTkyNzc0NzEsImlzcyI6IkF1dGhUb2tlbkFQSSIsInBlcm1pc3Npb25zIjp7InJlc291cmNlIjpbInBlcm1pc3Npb24iXX0sInN1YiI6InVzZXJTdWNjZXNzIn0.gCSXyLP8eInotFu7WvYg-Y3Gn4piNtR0LhYPgYm6RBI"

		encryptedToken, err := testAuthTokenServices.jwe.encryptToken([]byte(signedToken), []byte(testAuthTokenServices.config.builder.jwtEncryptionSecretKey))
		assert.Nil(t, err)
		assert.NotNil(t, encryptedToken)

		t.Log("Encrypted Token:", string(encryptedToken))
	})

	t.Run("Error invalid secret key size", func(t *testing.T) {

		signedToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE5NTkyNzc0NzEsImlzcyI6IkF1dGhUb2tlbkFQSSIsInBlcm1pc3Npb25zIjp7InJlc291cmNlIjpbInBlcm1pc3Npb24iXX0sInN1YiI6InVzZXJTdWNjZXNzIn0.gCSXyLP8eInotFu7WvYg-Y3Gn4piNtR0LhYPgYm6RBI"
		secretKey := "invalid size"

		encryptedToken, err := testAuthTokenServices.jwe.encryptToken([]byte(signedToken), []byte(secretKey))
		assert.NotNil(t, err)
		assert.Nil(t, encryptedToken)

		// Probable error: failed to encrypt payload: failed to encrypt payload: failed to crypt content: failed to fetch AEAD: cipher: failed to create AES cipher for CBC: failed to execute block cipher function: crypto/aes: invalid key size 6
		// Cannot test the err value because that way we would depend on the jwt library.
		t.Log("err value:", err)
	})
}
