package services

import (
	"testing"

	"github.com/maxiancillotti/access-control/internal/dto"
	"github.com/maxiancillotti/access-control/internal/services/internal"
	"github.com/maxiancillotti/access-control/internal/services/svcerr"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestValidateAdminCredentials(t *testing.T) {

	type testCase struct {
		name                 string
		userCredInput        *dto.AdminCredentials
		expectedSvcErrOutput *svcerr.ServiceError
	}

	table := make([]testCase, 0)

	table = append(table, testCase{
		name: "Success",
		userCredInput: &dto.AdminCredentials{
			Username: "usernameSuccess",
			Password: "APIUserPassword", // Will be check against hashed pw
		},
		expectedSvcErrOutput: nil,
	})

	table = append(table, testCase{
		name: "Error. Invalid input ID. Invalid credentials: username.",
		userCredInput: &dto.AdminCredentials{
			Username: "usernameErrInvalidInputID",
			Password: "pw",
		},
		expectedSvcErrOutput: svcerr.New(
			errors.Wrap(errors.New("error msg invalid input"), internal.ErrMsgInvalidUsername),
			internal.ErrorCategoryInvalidCredentials,
		),
	})

	table = append(table, testCase{
		name: "Error. Invalid input ID. Invalid credentials: username.",
		userCredInput: &dto.AdminCredentials{
			Username: "usernameErrInternal",
			Password: "pw",
		},
		expectedSvcErrOutput: svcerr.New(
			errors.New("error msg internal"),
			internal.ErrorCategoryInternal,
		),
	})

	table = append(table, testCase{
		name: "Error. User disabled.",
		userCredInput: &dto.AdminCredentials{
			Username: "usernameSuccessStateDisabled",
			Password: "pw",
		},
		expectedSvcErrOutput: svcerr.New(
			errors.New(internal.ErrMsgUserDisabled),
			//internal.ErrorCategoryUserDisabled,
			internal.ErrorCategoryInvalidCredentials,
		),
	})

	table = append(table, testCase{
		name: "Error. Invalid credentials: password.",
		userCredInput: &dto.AdminCredentials{
			Username: "usernameSuccessInvalidPassword",
			Password: "pw",
		},
		expectedSvcErrOutput: svcerr.New(
			errors.New(internal.ErrMsgInvalidPassword),
			internal.ErrorCategoryInvalidCredentials,
		),
	})

	for _, test := range table {

		t.Run(test.name, func(t *testing.T) {

			userID, svcErr := testAdminssAuthServices.validateAdminCredentials(test.userCredInput)

			if test.expectedSvcErrOutput == nil {
				assert.Nil(t, svcErr)
				assert.Greater(t, int(userID), 0)
			} else {

				assert.NotNil(t, svcErr)
				assert.Zero(t, int(userID))

				assert.Equal(t, test.expectedSvcErrOutput.Category(), svcErr.Category())

				if test.userCredInput.Username == "usernameSuccessInvalidPassword" {
					assert.Contains(t, svcErr.ErrorValue().Error(), test.expectedSvcErrOutput.ErrorValue().Error())
				} else {
					assert.Equal(t, test.expectedSvcErrOutput.ErrorValue().Error(), svcErr.ErrorValue().Error())
				}

			}
		})
	}
}
