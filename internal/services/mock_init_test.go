package services

import "github.com/maxiancillotti/access-control/internal/mock"

var (
	usrDALMock                      = &mock.UsersDALMock{}
	testUsersServices UsersServices = NewUsersServices(usrDALMock)

	usrRestPermDALMock                                     = &mock.UsersRESTPermissionsDALMock{}
	testUsersRESTPermServices UsersRESTPermissionsServices = NewUsersRESTPermissionsServices(usrRestPermDALMock, testUsersServices)

	usersSvcMock          = &usersInteractorMock{}
	usersRestPermSvcMock  = &usersRESTPermissionsInteractorMock{}
	testUsersAuthServices = newUsersAuthServices(usersSvcMock, usersRestPermSvcMock)

	admDALMock                        = &mock.AdminsDALMock{}
	testAdminsServices AdminsServices = NewAdminsServices(admDALMock)

	adminsSvcMock           = &adminsInteractorMock{}
	testAdminssAuthServices = newAdminsAuthServices(adminsSvcMock)

	adminsAuthSvcMock = &adminsAuthInteractorMock{}
	usersAuthSvcMock  = &usersAuthInteractorMock{}
	tokenSvcMock      = &authTokenMock{}

	testAuthServices AuthServices = &authInteractor{
		adminsAuthSvc: adminsAuthSvcMock,
		//userAuthSvc: testUsersAuthServices, // NOT MOCKED
		userAuthSvc: usersAuthSvcMock,
		tokenSvc:    tokenSvcMock,
	}

	rscsDALMock                        = &mock.ResourcesDALMock{}
	testResourcesSvc ResourcesServices = &resourcesInteractor{rscsDALMock}

	mthdsDALMock                           = &mock.HttpmethodsDALMock{}
	testHttpMethodsSvc HttpMethodsServices = &httpMethodsInteractor{mthdsDALMock}
)
