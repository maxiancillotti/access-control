package authtoken

import (
	"testing"

	"github.com/maxiancillotti/access-control/internal/domain"
	"github.com/maxiancillotti/access-control/internal/services/internal"

	"github.com/stretchr/testify/assert"
)

func TestGenerateToken(t *testing.T) {

	t.Run("Success", func(t *testing.T) {
		var userID uint = 1

		//usrPerm := make(domain.UserPermissions)
		//usrPerm["resource"] = append(usrPerm["resource"], "permission")

		usrPerm := domain.UserPermissions{
			domain.RESTpermissionCategory: domain.RESTPermissionsPathsMethods{
				"/resource": {"POST"},
			},
		}

		token, err := testAuthTokenServices.GenerateToken(userID, usrPerm)
		assert.Nil(t, err)
		assert.NotNil(t, token)

		t.Log("Token created:", string(token))
	})

	t.Run("Signing Error", func(t *testing.T) {
		var userID uint = 1

		//usrPerm := make(domain.UserPermissions)
		//usrPerm["resource"] = append(usrPerm["resource"], "permission")

		usrPerm := domain.UserPermissions{
			domain.RESTpermissionCategory: domain.RESTPermissionsPathsMethods{
				"/resource": {"POST"},
			},
		}

		// Error trigger
		testAuthTokenServices.config.builder.jwtSigningSecretKey = ""

		// CALL
		token, svcErr := testAuthTokenServices.GenerateToken(userID, usrPerm)

		// TEST
		assert.NotNil(t, svcErr)
		assert.Nil(t, token)

		//assert.ErrorIs(t, svcErr.ErrorValue(), internal.ErrMsgSignFailed)
		assert.Contains(t, svcErr.ErrorValue().Error(), internal.ErrMsgSignFailed.Error())

		assert.Equal(t, internal.ErrorCategoryInternal, svcErr.Category())

		// Correcting previous error trigger
		testAuthTokenServices.config.builder.jwtSigningSecretKey = testSigningSecretKey
	})

	t.Run("Token encryption error", func(t *testing.T) {
		var userID uint = 1

		//usrPerm := make(domain.UserPermissions)
		//usrPerm["resource"] = append(usrPerm["resource"], "permission")

		usrPerm := domain.UserPermissions{
			domain.RESTpermissionCategory: domain.RESTPermissionsPathsMethods{
				"/resource": {"POST"},
			},
		}

		// Error trigger
		testAuthTokenServices.config.builder.jwtEncryptionSecretKey = ""

		// CALL
		token, svcErr := testAuthTokenServices.GenerateToken(userID, usrPerm)

		// TEST
		assert.NotNil(t, svcErr)
		assert.Nil(t, token)

		//assert.ErrorIs(t, svcErr.ErrorValue(), internal.ErrMsgFailedToEncryptPayload)
		assert.Contains(t, svcErr.ErrorValue().Error(), internal.ErrMsgFailedToEncryptPayload.Error())

		assert.Equal(t, internal.ErrorCategoryInternal, svcErr.Category())

		// Correcting previous error trigger
		testAuthTokenServices.config.builder.jwtEncryptionSecretKey = testEncryptionSecretyKey
	})

}

func TestValidateToken(t *testing.T) {

	//(validationData dto.ValidationRequest) (svcErr error)

	t.Run("Success", func(t *testing.T) {

		// validationData := dto.ValidationRequest{
		// 	// Token generated on previous test.
		// 	// 10 years of validity.
		// 	// Generic resource and method names, like the field names.
		// 	Token:             "eyJhbGciOiJkaXIiLCJlbmMiOiJBMTI4Q0JDLUhTMjU2In0..6QLqr6pXAP-gfclqIal4pw.S7R6RLB6qhH8u2N4RPgRTm0j4GccFVcQJKCX6Q9FaGXdVPRzjCe5Ndbl-vvM_AC-0S29DTbLUshqG15r-9LKLL8KropVQUDQaNiVTWwPOGOEWyRZ4GMLgvRcbCI9HrDtdN_QBq0ZHkheQmbCW65ec1HtawQCUS-Db9xRLNgjC4deR_6zVDSOxZb32ShrMfssOjvcvRspz5zt_5aRLiALlw1TRN6ytLO2mBr2elTvEdtiq-ftt-5ICNlQbhlyJvlKcp-E0nf6HVusaJAKV7cSXWYS_cyyCgHFE_XO-_Pzcgs.OUNiW0sMaTRSppXEHr_t2A",
		// 	ResourceRequested: "resource",
		// 	MethodRequested:   "permission",
		// }

		token := "eyJhbGciOiJkaXIiLCJlbmMiOiJBMTI4Q0JDLUhTMjU2In0..6QLqr6pXAP-gfclqIal4pw.S7R6RLB6qhH8u2N4RPgRTm0j4GccFVcQJKCX6Q9FaGXdVPRzjCe5Ndbl-vvM_AC-0S29DTbLUshqG15r-9LKLL8KropVQUDQaNiVTWwPOGOEWyRZ4GMLgvRcbCI9HrDtdN_QBq0ZHkheQmbCW65ec1HtawQCUS-Db9xRLNgjC4deR_6zVDSOxZb32ShrMfssOjvcvRspz5zt_5aRLiALlw1TRN6ytLO2mBr2elTvEdtiq-ftt-5ICNlQbhlyJvlKcp-E0nf6HVusaJAKV7cSXWYS_cyyCgHFE_XO-_Pzcgs.OUNiW0sMaTRSppXEHr_t2A"

		permissions, err := testAuthTokenServices.ValidateToken(token)
		assert.Nil(t, err)
		assert.NotNil(t, permissions)
	})

	t.Run("Decrypt error", func(t *testing.T) {

		// validationData := dto.ValidationRequest{
		// 	// Error trigger: same token as before with "1" at the beginning.
		// 	Token:             "1eyJhbGciOiJkaXIiLCJlbmMiOiJBMTI4Q0JDLUhTMjU2In0..6QLqr6pXAP-gfclqIal4pw.S7R6RLB6qhH8u2N4RPgRTm0j4GccFVcQJKCX6Q9FaGXdVPRzjCe5Ndbl-vvM_AC-0S29DTbLUshqG15r-9LKLL8KropVQUDQaNiVTWwPOGOEWyRZ4GMLgvRcbCI9HrDtdN_QBq0ZHkheQmbCW65ec1HtawQCUS-Db9xRLNgjC4deR_6zVDSOxZb32ShrMfssOjvcvRspz5zt_5aRLiALlw1TRN6ytLO2mBr2elTvEdtiq-ftt-5ICNlQbhlyJvlKcp-E0nf6HVusaJAKV7cSXWYS_cyyCgHFE_XO-_Pzcgs.OUNiW0sMaTRSppXEHr_t2A",
		// 	ResourceRequested: "resource",
		// 	MethodRequested:   "permission",
		// }

		token := "1eyJhbGciOiJkaXIiLCJlbmMiOiJBMTI4Q0JDLUhTMjU2In0..6QLqr6pXAP-gfclqIal4pw.S7R6RLB6qhH8u2N4RPgRTm0j4GccFVcQJKCX6Q9FaGXdVPRzjCe5Ndbl-vvM_AC-0S29DTbLUshqG15r-9LKLL8KropVQUDQaNiVTWwPOGOEWyRZ4GMLgvRcbCI9HrDtdN_QBq0ZHkheQmbCW65ec1HtawQCUS-Db9xRLNgjC4deR_6zVDSOxZb32ShrMfssOjvcvRspz5zt_5aRLiALlw1TRN6ytLO2mBr2elTvEdtiq-ftt-5ICNlQbhlyJvlKcp-E0nf6HVusaJAKV7cSXWYS_cyyCgHFE_XO-_Pzcgs.OUNiW0sMaTRSppXEHr_t2A"

		// CALL
		permissions, svcErr := testAuthTokenServices.ValidateToken(token)

		// TEST
		assert.Nil(t, permissions)
		assert.NotNil(t, svcErr)

		//assert.ErrorIs(t, svcErr.ErrorValue(), internal.ErrMsgFailedToDecryptToken)
		assert.Contains(t, svcErr.ErrorValue().Error(), internal.ErrMsgFailedToDecryptToken.Error())

		assert.Equal(t, internal.ErrorCategoryInvalidToken, svcErr.Category())
	})

	t.Run("Verify error", func(t *testing.T) {
		// validationData := dto.ValidationRequest{
		// 	// Error trigger: old expired token.
		// 	Token:             "eyJhbGciOiJkaXIiLCJlbmMiOiJBMTI4Q0JDLUhTMjU2In0..s46rXL3R8ceNgG4r0lcMzw.7aVyoFReUBMt7tm6ZiPVejP7-JQrcVjZ09N_8_arnjEmXzgnLic2KZ9CWU6OZbMNxqvDvT-8iFOT7RBOGIU5totPdRJcPqCDLo0lVdk0_FM_CAw2iylrcs9xNMXgKRn3-ZeFS0GhS9KQLB2z533AmqVhZUOPPTPCaBgzM4lqy26uq-4qyypPSLnxp2s6Vj7er5IpeX6GREmixYZV0r4-t9YgvmLwMCRBKYtprAmDIxW7ONm55ls1n3fCOqVftsHjUNwzUF0ZtWxBzFvuaqqhrPPOsoqQilgAFZfatOqW2IM.S39VyLZa9GQ8bbAPHoDmug",
		// 	ResourceRequested: "/customers",
		// 	MethodRequested:   "POST",
		// }

		token := "eyJhbGciOiJkaXIiLCJlbmMiOiJBMTI4Q0JDLUhTMjU2In0..s46rXL3R8ceNgG4r0lcMzw.7aVyoFReUBMt7tm6ZiPVejP7-JQrcVjZ09N_8_arnjEmXzgnLic2KZ9CWU6OZbMNxqvDvT-8iFOT7RBOGIU5totPdRJcPqCDLo0lVdk0_FM_CAw2iylrcs9xNMXgKRn3-ZeFS0GhS9KQLB2z533AmqVhZUOPPTPCaBgzM4lqy26uq-4qyypPSLnxp2s6Vj7er5IpeX6GREmixYZV0r4-t9YgvmLwMCRBKYtprAmDIxW7ONm55ls1n3fCOqVftsHjUNwzUF0ZtWxBzFvuaqqhrPPOsoqQilgAFZfatOqW2IM.S39VyLZa9GQ8bbAPHoDmug"

		// CALL
		permissions, svcErr := testAuthTokenServices.ValidateToken(token)

		// TEST
		assert.Nil(t, permissions)
		assert.NotNil(t, svcErr)

		//assert.ErrorIs(t, svcErr.ErrorValue(), internal.ErrMsgInvalidToken)
		assert.Contains(t, svcErr.ErrorValue().Error(), internal.ErrMsgInvalidToken.Error())

		assert.Equal(t, internal.ErrorCategoryInvalidToken, svcErr.Category())
	})

	t.Run("No permissions key error", func(t *testing.T) {
		// validationData := dto.ValidationRequest{
		// 	// Error trigger: token with no permissions key.
		// 	Token:             "eyJhbGciOiJkaXIiLCJlbmMiOiJBMTI4Q0JDLUhTMjU2In0..hT5regCKsiMWohBPvbDZ_w.NR4fbIIacqgofPCiRF3kQyYbjCbRkz9GE5LmTQWfoh3TVU_UfTgrvpvTq6D9blCGKVukTv0FCi-zI0y3MFZNXyHhfQQo4n0v6kbqtuUAeOCPjiKWraayTsvB9gcoOwTP0uhKV-mrhNMvD98CePvTHGZ28W0LYzB-uM8H13ElEmu3kz9X7OZnFQjzcVS0d2vJ6fLLPkyzmXyp58iQBJ1Esg.HdQj6GElULrOwIX9BWyuwA",
		// 	ResourceRequested: "/customers",
		// 	MethodRequested:   "POST",
		// }

		token := "eyJhbGciOiJkaXIiLCJlbmMiOiJBMTI4Q0JDLUhTMjU2In0..hT5regCKsiMWohBPvbDZ_w.NR4fbIIacqgofPCiRF3kQyYbjCbRkz9GE5LmTQWfoh3TVU_UfTgrvpvTq6D9blCGKVukTv0FCi-zI0y3MFZNXyHhfQQo4n0v6kbqtuUAeOCPjiKWraayTsvB9gcoOwTP0uhKV-mrhNMvD98CePvTHGZ28W0LYzB-uM8H13ElEmu3kz9X7OZnFQjzcVS0d2vJ6fLLPkyzmXyp58iQBJ1Esg.HdQj6GElULrOwIX9BWyuwA"

		// CALL
		permissions, svcErr := testAuthTokenServices.ValidateToken(token)

		// TEST
		assert.Nil(t, permissions)
		assert.NotNil(t, svcErr)

		//assert.ErrorIs(t, svcErr.ErrorValue(), internal.ErrMsgTokenDoesntClaimAnyPermissions)
		assert.Contains(t, svcErr.ErrorValue().Error(), internal.ErrMsgTokenDoesntClaimAnyPermissions.Error())

		assert.Equal(t, internal.ErrorCategorySemanticallyUnprocesable, svcErr.Category())
	})

	t.Run("Validation error", func(t *testing.T) {
		// validationData := dto.ValidationRequest{
		// 	// Error trigger: token with empty value for permissions key.
		// 	Token:             "eyJhbGciOiJkaXIiLCJlbmMiOiJBMTI4Q0JDLUhTMjU2In0..hT5regCKsiMWohBPvbDZ_w.NR4fbIIacqgofPCiRF3kQyYbjCbRkz9GE5LmTQWfoh3TVU_UfTgrvpvTq6D9blCGKVukTv0FCi-zI0y3MFZNXyHhfQQo4n0v6kbqtuUAeOCPjiKWraayTsvB9gcoOwTP0uhKV-mrhNMvD98CePvTHGZ28W0LYzB-uM8H13ElEmu3kz9X7OZnFQjzcVS0d2vJ6fLLPkyzmXyp58iQBJ1Esg.HdQj6GElULrOwIX9BWyuwA",
		// 	ResourceRequested: "/customers",
		// 	MethodRequested:   "POST",
		// }

		token := "eyJhbGciOiJkaXIiLCJlbmMiOiJBMTI4Q0JDLUhTMjU2In0..hT5regCKsiMWohBPvbDZ_w.NR4fbIIacqgofPCiRF3kQyYbjCbRkz9GE5LmTQWfoh3TVU_UfTgrvpvTq6D9blCGKVukTv0FCi-zI0y3MFZNXyHhfQQo4n0v6kbqtuUAeOCPjiKWraayTsvB9gcoOwTP0uhKV-mrhNMvD98CePvTHGZ28W0LYzB-uM8H13ElEmu3kz9X7OZnFQjzcVS0d2vJ6fLLPkyzmXyp58iQBJ1Esg.HdQj6GElULrOwIX9BWyuwA"

		// CALL
		permissions, svcErr := testAuthTokenServices.ValidateToken(token)

		t.Log("err:", svcErr)

		// TEST
		assert.Nil(t, permissions)
		assert.NotNil(t, svcErr)
	})
}
