package services

import (
	"testing"

	"github.com/maxiancillotti/access-control/internal/domain"
	"github.com/maxiancillotti/access-control/internal/dto"
	"github.com/maxiancillotti/access-control/internal/services/internal"
	"github.com/maxiancillotti/access-control/internal/services/svcerr"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestValidateUserCredentials(t *testing.T) {

	type testCase struct {
		name                 string
		userCredInput        *dto.UserCredentials
		expectedSvcErrOutput *svcerr.ServiceError
	}

	table := make([]testCase, 0)

	table = append(table, testCase{
		name: "Success",
		userCredInput: &dto.UserCredentials{
			Username: "usernameSuccess",
			Password: "APIUserPassword", // Will be check against hashed pw
		},
		expectedSvcErrOutput: nil,
	})

	table = append(table, testCase{
		name: "Error. Invalid input ID. Invalid credentials: username.",
		userCredInput: &dto.UserCredentials{
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
		userCredInput: &dto.UserCredentials{
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
		userCredInput: &dto.UserCredentials{
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
		userCredInput: &dto.UserCredentials{
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

			userID, svcErr := testUsersAuthServices.validateUserCredentials(test.userCredInput)

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

func TestGetUserPermissionsByUserID(t *testing.T) {

	type testCase struct {
		name                  string
		userIDInput           uint
		expectedUsrPermOutput domain.UserPermissions
		expectedSvcErrOutput  *svcerr.ServiceError
	}

	table := make([]testCase, 0)

	table = append(table, testCase{
		name:        "Success",
		userIDInput: 1,
		expectedUsrPermOutput: domain.UserPermissions{
			domain.RESTpermissionCategory: domain.RESTPermissionsPathsMethods{
				"/customers": {"POST"},
			},
		},
		expectedSvcErrOutput: nil,
	})

	table = append(table, testCase{
		name:                  "Error. Empty result, user doesn't have any permissions. Internal.",
		userIDInput:           2,
		expectedUsrPermOutput: nil,
		expectedSvcErrOutput: svcerr.New(
			errors.Wrap(errors.New("err msg empty result"), internal.ErrMsgUserDoesntHaveAnyPermissions),
			internal.ErrorCategoryInternal,
		),
	})

	table = append(table, testCase{
		name:                  "Error. Internal.",
		userIDInput:           3,
		expectedUsrPermOutput: nil,
		expectedSvcErrOutput: svcerr.New(
			errors.New("err msg internal"),
			internal.ErrorCategoryInternal,
		),
	})

	table = append(table, testCase{
		name:                  "Error. Invalid input id, error retrieving permissions. Internal.",
		userIDInput:           4,
		expectedUsrPermOutput: nil,
		expectedSvcErrOutput: svcerr.New(
			errors.Wrap(errors.New("err msg invalid input id"), internal.ErrMsgRetrievingPermissions),
			internal.ErrorCategoryInternal,
		),
	})

	for _, test := range table {

		t.Run(test.name, func(t *testing.T) {

			usrPerm, svcErr := testUsersAuthServices.getUserPermissionsByUserID(test.userIDInput)

			if test.expectedSvcErrOutput == nil {
				assert.Nil(t, svcErr)
				assert.NotNil(t, usrPerm)

				assert.Equal(t,
					usrPerm[domain.RESTpermissionCategory]["/customers"][0],
					test.expectedUsrPermOutput[domain.RESTpermissionCategory]["/customers"][0],
				)

			} else {
				assert.NotNil(t, svcErr)
				assert.Nil(t, usrPerm)

				assert.Equal(t, test.expectedSvcErrOutput.ErrorValue().Error(), svcErr.ErrorValue().Error())
				assert.Equal(t, test.expectedSvcErrOutput.Category(), svcErr.Category())
			}
		})
	}
}

// TO-DO: OTHER CASES

// func TestValidatePermissions(t *testing.T) {

// 	t.Run("Success", func(t *testing.T) {

// 		resourceRequested := "/customers"
// 		methodRequested := "POST"

// 		// var permissions interface{} = domain.UserPermissions{
// 		// 	domain.RESTpermissionCategory: domain.RESTPermissionsPathsMethods{
// 		// 		"/customers": {"POST"},
// 		// 	},
// 		// }

// 		var permissions = `
// 		{
// 			"REST" : {
// 				"/customers": ["POST"]
// 			}
// 		}`

// 		err := testUsersAuthServices.validatePermissions(permissions, resourceRequested, methodRequested, domain.RESTpermissionCategory)

// 		assert.Nil(t, err)
// 	})
// }
