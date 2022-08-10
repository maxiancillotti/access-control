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

func TestUserRESTPermCreate(t *testing.T) {

	type testCase struct {
		name                      string
		permissionUserIDInput     uint
		permissionResourceIDInput uint
		permissionMethodIDInput   uint
		expectedErrOutput         *svcerr.ServiceError
	}

	table := make([]testCase, 0)

	table = append(table, testCase{
		name:                      "Success",
		permissionUserIDInput:     2,
		permissionResourceIDInput: 1,
		permissionMethodIDInput:   1,
		expectedErrOutput:         nil,
	})

	table = append(table, testCase{
		name:                      "Error. Permission insert failed. Internal.",
		permissionUserIDInput:     5,
		permissionResourceIDInput: 1,
		permissionMethodIDInput:   1,
		expectedErrOutput: svcerr.New(
			errors.Wrapf(errors.New("error insert user 5"), internal.ErrMsgFmtInsertFailed, "permission"),
			internal.ErrorCategoryInternal,
		),
	})

	table = append(table, testCase{
		name:                      "Error. Permission already exists. Invalid input ID.",
		permissionUserIDInput:     1,
		permissionResourceIDInput: 1,
		permissionMethodIDInput:   1,
		expectedErrOutput: svcerr.New(
			errors.Errorf(internal.ErrMsgFmtAlreadyExists, "permission"),
			internal.ErrorCategoryInvalidInputID,
		),
	})

	table = append(table, testCase{
		name:                      "Error. Failed to check if permission exists. Internal.",
		permissionUserIDInput:     3,
		permissionResourceIDInput: 1,
		permissionMethodIDInput:   1,
		expectedErrOutput: svcerr.New(
			errors.Wrapf(errors.New("err user 3"), internal.ErrMsgFmtFailedToCheckIfExists, "permission"),
			internal.ErrorCategoryInternal,
		),
	})

	table = append(table, testCase{
		name:                      "Error. Failed to check if permission's relationships exists. Internal.",
		permissionUserIDInput:     6,
		permissionResourceIDInput: 1,
		permissionMethodIDInput:   1,
		expectedErrOutput: svcerr.New(
			errors.Wrapf(errors.New("err user 6"), internal.ErrMsgFmtFailedToCheckIfExists, "permission's relationships"),
			internal.ErrorCategoryInternal,
		),
	})

	for _, test := range table {

		t.Run(test.name, func(t *testing.T) {

			urp := dto.UserRESTPermission{
				UserID: &test.permissionUserIDInput,
				Permission: &dto.Permission{
					ResourceID: &test.permissionResourceIDInput,
					MethodID:   &test.permissionMethodIDInput,
				},
			}

			err := testUsersRESTPermServices.Create(urp)

			if test.expectedErrOutput == nil {
				assert.Nil(t, err)
			} else {
				assert.NotNil(t, err)

				svcErr := err.(*svcerr.ServiceError)

				assert.Equal(t, test.expectedErrOutput.ErrorValue().Error(), svcErr.ErrorValue().Error())
				assert.Equal(t, test.expectedErrOutput.Category(), svcErr.Category())
			}
		})
	}
}

func TestUserRESTPermDelete(t *testing.T) {

	type testCase struct {
		name                      string
		permissionUserIDInput     uint
		permissionResourceIDInput uint
		permissionMethodIDInput   uint
		expectedErrOutput         *svcerr.ServiceError
	}

	table := make([]testCase, 0)

	table = append(table, testCase{
		name:                      "Success",
		permissionUserIDInput:     1,
		permissionResourceIDInput: 1,
		permissionMethodIDInput:   1,
		expectedErrOutput:         nil,
	})

	table = append(table, testCase{
		name:                      "Error. Permission does no exist. Invalid input id.",
		permissionUserIDInput:     2,
		permissionResourceIDInput: 1,
		permissionMethodIDInput:   1,
		expectedErrOutput: svcerr.New(
			errors.Errorf(internal.ErrMsgFmtDoesNotExist, "permission"),
			internal.ErrorCategoryInvalidInputID,
		),
	})

	// table = append(table, testCase{
	// 	name:            "Error. Failed to check if permission already exists. Internal.",
	// 	permissionInput: dto.UserRESTPermission{UserID: 3, ResourceID: 1, MethodID: 1},
	// 	expectedErrOutput: svcerr.New(
	// 		errors.Wrapf(errors.New("err user 3"), internal.ErrMsgFmtFailedToCheckIfExists, "permission"),
	// 		internal.ErrorCategoryInternal,
	// 	),
	// })

	table = append(table, testCase{
		name:                      "Error. Permission delete failed. Internal.",
		permissionUserIDInput:     4,
		permissionResourceIDInput: 1,
		permissionMethodIDInput:   1,
		expectedErrOutput: svcerr.New(
			errors.Wrapf(errors.New("delete err"), internal.ErrMsgFmtDeleteFailed, "permission"),
			internal.ErrorCategoryInternal,
		),
	})

	for _, test := range table {

		t.Run(test.name, func(t *testing.T) {

			urp := dto.UserRESTPermission{
				UserID: &test.permissionUserIDInput,
				Permission: &dto.Permission{
					ResourceID: &test.permissionResourceIDInput,
					MethodID:   &test.permissionMethodIDInput,
				},
			}

			err := testUsersRESTPermServices.Delete(urp)

			if test.expectedErrOutput == nil {
				assert.Nil(t, err)
			} else {
				assert.NotNil(t, err)

				svcErr := err.(*svcerr.ServiceError)

				assert.Equal(t, test.expectedErrOutput.ErrorValue().Error(), svcErr.ErrorValue().Error())
				assert.Equal(t, test.expectedErrOutput.Category(), svcErr.Category())
			}

		})
	}
}

func TestUserRESTPermExistsOrErr(t *testing.T) {

	type testCase struct {
		name              string
		permissionInput   *domain.UserRESTPermission
		expectedErrOutput *svcerr.ServiceError
	}

	table := make([]testCase, 0)

	table = append(table, testCase{
		name:              "Success",
		permissionInput:   &domain.UserRESTPermission{UserID: 1, ResourceID: 1, MethodID: 1},
		expectedErrOutput: nil,
	})

	table = append(table, testCase{
		name:            "Error. Permission does no exist. Invalid input id.",
		permissionInput: &domain.UserRESTPermission{UserID: 2, ResourceID: 1, MethodID: 1},
		expectedErrOutput: svcerr.New(
			errors.Errorf(internal.ErrMsgFmtDoesNotExist, "permission"),
			internal.ErrorCategoryInvalidInputID,
		),
	})

	table = append(table, testCase{
		name:            "Error. Failed to check if permission already exists. Internal.",
		permissionInput: &domain.UserRESTPermission{UserID: 3, ResourceID: 1, MethodID: 1},
		expectedErrOutput: svcerr.New(
			errors.Wrapf(errors.New("err user 3"), internal.ErrMsgFmtFailedToCheckIfExists, "permission"),
			internal.ErrorCategoryInternal,
		),
	})

	for _, test := range table {

		t.Run(test.name, func(t *testing.T) {
			testUsersRESTPermSvc := testUsersRESTPermServices.(*usersRESTPermissionsInteractor)
			svcErr := testUsersRESTPermSvc.existsOrErr(test.permissionInput)

			if test.expectedErrOutput == nil {
				assert.Nil(t, svcErr)
			} else {
				assert.NotNil(t, svcErr)

				assert.Equal(t, test.expectedErrOutput.ErrorValue().Error(), svcErr.ErrorValue().Error())
				assert.Equal(t, test.expectedErrOutput.Category(), svcErr.Category())
			}

		})
	}
}

func TestUserRESTPermRelationshipsExistsOrErr(t *testing.T) {

	type testCase struct {
		name              string
		permissionInput   *domain.UserRESTPermission
		expectedErrOutput *svcerr.ServiceError
	}

	table := make([]testCase, 0)

	table = append(table, testCase{
		name:              "Success",
		permissionInput:   &domain.UserRESTPermission{UserID: 2, ResourceID: 1, MethodID: 1},
		expectedErrOutput: nil,
	})

	table = append(table, testCase{
		name:            "Error. Failed to check if permission's relationships exists. Internal.",
		permissionInput: &domain.UserRESTPermission{UserID: 6, ResourceID: 1, MethodID: 1},
		expectedErrOutput: svcerr.New(
			errors.Wrapf(errors.New("err user 6"), internal.ErrMsgFmtFailedToCheckIfExists, "permission's relationships"),
			internal.ErrorCategoryInternal,
		),
	})

	table = append(table, testCase{
		name:            "Error. UserID does no exist. Invalid input id.",
		permissionInput: &domain.UserRESTPermission{UserID: 7, ResourceID: 1, MethodID: 1},
		expectedErrOutput: svcerr.New(
			errors.Errorf(internal.ErrMsgFmtDoesNotExist, "UserID"),
			internal.ErrorCategoryInvalidInputID,
		),
	})

	table = append(table, testCase{
		name:            "Error. ResourceID does no exist. Invalid input id.",
		permissionInput: &domain.UserRESTPermission{UserID: 8, ResourceID: 1, MethodID: 1},
		expectedErrOutput: svcerr.New(
			errors.Errorf(internal.ErrMsgFmtDoesNotExist, "ResourceID"),
			internal.ErrorCategoryInvalidInputID,
		),
	})

	table = append(table, testCase{
		name:            "Error. HttpMethodID does no exist. Invalid input id.",
		permissionInput: &domain.UserRESTPermission{UserID: 9, ResourceID: 1, MethodID: 1},
		expectedErrOutput: svcerr.New(
			errors.Errorf(internal.ErrMsgFmtDoesNotExist, "HttpMethodID"),
			internal.ErrorCategoryInvalidInputID,
		),
	})

	table = append(table, testCase{
		name:            "Error. Unexpected value 0 for 'exists' when err == nil.",
		permissionInput: &domain.UserRESTPermission{UserID: 10, ResourceID: 1, MethodID: 1},
		expectedErrOutput: svcerr.New(
			errors.Wrapf(errors.Errorf("unexpected value received from DAL for 'exists' = %d", 0),
				internal.ErrMsgFmtFailedToCheckIfExists, "permission's relationships"),
			internal.ErrorCategoryInternal,
		),
	})

	for _, test := range table {

		t.Run(test.name, func(t *testing.T) {
			testUsersRESTPermSvc := testUsersRESTPermServices.(*usersRESTPermissionsInteractor)
			svcErr := testUsersRESTPermSvc.relationshipsExistsOrErr(test.permissionInput)

			if test.expectedErrOutput == nil {
				assert.Nil(t, svcErr)
			} else {
				assert.NotNil(t, svcErr)

				assert.Equal(t, test.expectedErrOutput.ErrorValue().Error(), svcErr.ErrorValue().Error())
				assert.Equal(t, test.expectedErrOutput.Category(), svcErr.Category())
			}

		})
	}
}

func TestUserRESTPermRetrieveAllByUserID(t *testing.T) {

	type testCase struct {
		name                      string
		userIDInput               uint
		expectedPermissionsOutput *dto.UserRESTPermissionsCollection
		expectedErrOutput         *svcerr.ServiceError
	}

	table := make([]testCase, 0)

	table = append(table, testCase{
		name:        "Success",
		userIDInput: 1,
		expectedPermissionsOutput: &dto.UserRESTPermissionsCollection{
			UserID: 1,
			Permissions: []dto.PermissionsIDs{
				{ResourceID: 1, MethodsIDs: []uint{1, 2, 3}},
			},
		},
		expectedErrOutput: nil,
	})

	table = append(table, testCase{
		name:                      "Error. User does no exist. Invalid input id.",
		userIDInput:               2,
		expectedPermissionsOutput: nil,
		expectedErrOutput: svcerr.New(
			errors.Errorf(internal.ErrMsgFmtDoesNotExist, "user"),
			internal.ErrorCategoryInvalidInputID,
		),
	})

	table = append(table, testCase{
		name:                      "Error. Permission retrieval failed. Internal.",
		userIDInput:               4,
		expectedPermissionsOutput: nil,
		expectedErrOutput: svcerr.New(
			errors.Wrapf(errors.New("SelectAllByUserID err userID 4"), internal.ErrMsgFmtRetrievalFailed, "permissions"),
			internal.ErrorCategoryInternal,
		),
	})

	table = append(table, testCase{
		name:                      "Error. Empty permission retrieval. Empty result.",
		userIDInput:               5,
		expectedPermissionsOutput: nil,
		expectedErrOutput: svcerr.New(
			errors.Errorf(internal.ErrMsgFmtEmptyResultForTheGivenID, "permissions", "userID"),
			internal.ErrorCategoryEmptyResult,
		),
	})

	for _, test := range table {

		t.Run(test.name, func(t *testing.T) {
			userPermissions, err := testUsersRESTPermServices.RetrieveAllByUserID(test.userIDInput)

			if test.expectedErrOutput == nil {
				assert.Nil(t, err)
				assert.NotNil(t, userPermissions)

				assert.Equal(t, test.expectedPermissionsOutput.UserID, userPermissions.UserID)
				assert.Equal(t, test.expectedPermissionsOutput.Permissions[0].ResourceID, userPermissions.Permissions[0].ResourceID)
				assert.Equal(t, test.expectedPermissionsOutput.Permissions[0].MethodsIDs[0], userPermissions.Permissions[0].MethodsIDs[0])

			} else {
				assert.NotNil(t, err)
				assert.Nil(t, userPermissions)

				svcErr := err.(*svcerr.ServiceError)

				assert.Equal(t, test.expectedErrOutput.ErrorValue().Error(), svcErr.ErrorValue().Error())
				assert.Equal(t, test.expectedErrOutput.Category(), svcErr.Category())
			}

		})
	}
}

func TestUserRESTPermRetrieveAllWithDescriptionsByUserID(t *testing.T) {

	type testCase struct {
		name                      string
		userIDInput               uint
		expectedPermissionsOutput *dto.UserRESTPermissionsDescriptionsCollection
		expectedErrOutput         *svcerr.ServiceError
	}

	table := make([]testCase, 0)

	case1InputResourceID := uint(10)
	case1InputResourcePath := "/customers"
	case1InputMethodID := uint(15)
	case1InputMethodName := "POST"

	table = append(table, testCase{
		name:        "Success",
		userIDInput: 1,
		expectedPermissionsOutput: &dto.UserRESTPermissionsDescriptionsCollection{
			UserID: 1,
			Permissions: []dto.PermissionsWithDescriptions{
				{
					Resource: dto.Resource{
						ID:   &case1InputResourceID,
						Path: &case1InputResourcePath,
					},
					Methods: []dto.HttpMethod{
						{
							ID:   &case1InputMethodID,
							Name: &case1InputMethodName,
						},
					},
				},
			},
		},
		expectedErrOutput: nil,
	})

	table = append(table, testCase{
		name:                      "Error. User does no exist. Invalid input id.",
		userIDInput:               2,
		expectedPermissionsOutput: nil,
		expectedErrOutput: svcerr.New(
			errors.Errorf(internal.ErrMsgFmtDoesNotExist, "user"),
			internal.ErrorCategoryInvalidInputID,
		),
	})

	table = append(table, testCase{
		name:                      "Error. Permission retrieval failed. Internal.",
		userIDInput:               4,
		expectedPermissionsOutput: nil,
		expectedErrOutput: svcerr.New(
			errors.Wrapf(errors.New("SelectAllWithDescriptionsByUserID err userID 4"), internal.ErrMsgFmtRetrievalFailed, "permissions"),
			internal.ErrorCategoryInternal,
		),
	})

	table = append(table, testCase{
		name:                      "Error. Empty permission retrieval. Empty result.",
		userIDInput:               5,
		expectedPermissionsOutput: nil,
		expectedErrOutput: svcerr.New(
			errors.Errorf(internal.ErrMsgFmtEmptyResultForTheGivenID, "permissions", "userID"),
			internal.ErrorCategoryEmptyResult,
		),
	})

	for _, test := range table {

		t.Run(test.name, func(t *testing.T) {
			userPermissions, err := testUsersRESTPermServices.RetrieveAllWithDescriptionsByUserID(test.userIDInput)

			if test.expectedErrOutput == nil {
				assert.Nil(t, err)
				assert.NotNil(t, userPermissions)

				assert.Equal(t, test.expectedPermissionsOutput.UserID, userPermissions.UserID)
				assert.Equal(t, test.expectedPermissionsOutput.Permissions[0].Resource.ID, userPermissions.Permissions[0].Resource.ID)
				assert.Equal(t, test.expectedPermissionsOutput.Permissions[0].Methods[0].ID, userPermissions.Permissions[0].Methods[0].ID)

			} else {
				assert.NotNil(t, err)
				assert.Nil(t, userPermissions)

				svcErr := err.(*svcerr.ServiceError)

				assert.Equal(t, test.expectedErrOutput.ErrorValue().Error(), svcErr.ErrorValue().Error())
				assert.Equal(t, test.expectedErrOutput.Category(), svcErr.Category())
			}

		})
	}
}

func TestUserRESTPermRetrieveAllPathMethodsByUserID(t *testing.T) {

	type testCase struct {
		name                      string
		userIDInput               uint
		expectedPathMethodsOutput domain.RESTPermissionsPathsMethods
		expectedErrOutput         *svcerr.ServiceError
	}

	table := make([]testCase, 0)

	table = append(table, testCase{
		name:                      "Success",
		userIDInput:               1,
		expectedPathMethodsOutput: domain.RESTPermissionsPathsMethods{"/customers": {"POST", "GET"}},
		expectedErrOutput:         nil,
	})

	table = append(table, testCase{
		name:                      "Error. User does no exist. Invalid input id.",
		userIDInput:               2,
		expectedPathMethodsOutput: nil,
		expectedErrOutput: svcerr.New(
			errors.Errorf(internal.ErrMsgFmtDoesNotExist, "user"),
			internal.ErrorCategoryInvalidInputID,
		),
	})

	table = append(table, testCase{
		name:                      "Error. Permission retrieval failed. Internal.",
		userIDInput:               4,
		expectedPathMethodsOutput: nil,
		expectedErrOutput: svcerr.New(
			errors.Wrapf(errors.New("SelectAllPathMethodsByUserID err userID 4"), internal.ErrMsgFmtRetrievalFailed, "permissions"),
			internal.ErrorCategoryInternal,
		),
	})

	table = append(table, testCase{
		name:                      "Error. Empty permission retrieval. Empty result.",
		userIDInput:               5,
		expectedPathMethodsOutput: nil,
		expectedErrOutput: svcerr.New(
			errors.Errorf(internal.ErrMsgFmtEmptyResultForTheGivenID, "permissions", "userID"),
			internal.ErrorCategoryEmptyResult,
		),
	})

	for _, test := range table {

		t.Run(test.name, func(t *testing.T) {
			permissions, err := testUsersRESTPermServices.retrieveAllPathMethodsByUserID(test.userIDInput)

			if test.expectedErrOutput == nil {
				assert.Nil(t, err)
				assert.NotNil(t, permissions)

				assert.Equal(t, test.expectedPathMethodsOutput["/customers"][0], permissions["/customers"][0])

			} else {
				assert.NotNil(t, err)
				assert.Nil(t, permissions)

				svcErr := err.(*svcerr.ServiceError)

				assert.Equal(t, test.expectedErrOutput.ErrorValue().Error(), svcErr.ErrorValue().Error())
				assert.Equal(t, test.expectedErrOutput.Category(), svcErr.Category())
			}

		})
	}
}
