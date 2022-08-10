package controllers

import (
	"github.com/maxiancillotti/access-control/internal/domain"
	"github.com/maxiancillotti/access-control/internal/dto"
	"github.com/maxiancillotti/access-control/internal/services"
	"github.com/pkg/errors"
)

type userRESTPermSvcMock struct {
	// Implementing this interface embedding it so the compiler can recognize unexported methods
	services.UsersRESTPermissionsServices
}

func (urp *userRESTPermSvcMock) Create(pDTO dto.UserRESTPermission) error {

	switch *pDTO.UserID {
	case 1:
		return nil
	case 0:
		return errors.New("invalid_input")
	case 2:
		return errors.New("internal_error")
	}
	panic("invalid userID received in mocked Create method")
}

func (urp *userRESTPermSvcMock) Delete(pDTO dto.UserRESTPermission) error {
	switch *pDTO.UserID {
	case 1:
		return nil
	case 0:
		return errors.New("invalid_input")
	case 2:
		return errors.New("internal_error")
	}
	panic("invalid userID received in mocked Delete method")
}

func (urp *userRESTPermSvcMock) RetrieveAllByUserID(userID uint) (*dto.UserRESTPermissionsCollection, error) {
	switch userID {
	case 1:
		return &dto.UserRESTPermissionsCollection{
			UserID: 1,
			Permissions: []dto.PermissionsIDs{
				{
					ResourceID: 10,
					MethodsIDs: []uint{1, 2, 3},
				},
			},
		}, nil
	case 0:
		return nil, errors.New("empty_result")
	case 2:
		return nil, errors.New("internal_error")
	}
	panic("invalid userID received in mocked RetrieveAllByUserID method")
}

func (urp *userRESTPermSvcMock) RetrieveAllWithDescriptionsByUserID(userID uint) (*dto.UserRESTPermissionsDescriptionsCollection, error) {
	switch userID {
	case 1:
		uintTwo := uint(2)
		uintTen := uint(10)
		resourcePathSuccess := "path"
		methodNameSuccess := "method"

		return &dto.UserRESTPermissionsDescriptionsCollection{
			UserID: 1,
			Permissions: []dto.PermissionsWithDescriptions{
				{
					Resource: dto.Resource{
						ID:   &uintTen,
						Path: &resourcePathSuccess,
					},
					Methods: []dto.HttpMethod{
						{
							ID:   &uintTwo,
							Name: &methodNameSuccess,
						},
					},
				},
			},
		}, nil

	case 0:
		return nil, errors.New("empty_result")
	case 2:
		return nil, errors.New("internal_error")
	}
	panic("invalid userID received in mocked RetrieveAllByUserID method")
}

// Not implemented
func (urp *userRESTPermSvcMock) retrieveAllPathMethodsByUserID(userID uint) (domain.RESTPermissionsPathsMethods, error) {
	return nil, nil
}
