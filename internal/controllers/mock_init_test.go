package controllers

import (
	"net/http"

	"github.com/maxiancillotti/access-control/internal/mock"
	"github.com/maxiancillotti/access-control/internal/mock/writerloggermock"
)

var (
	pstrMock mock.PresenterMock
	lgrMock  = mock.LoggerMock
	// TODO: MOCK LOGGER TOO SO THE INTERNAL ERRORS CAN BE SEEN AND TESTED
	// And obviously, so the test doesn't depend on the logger too, duh.
	wrlgrMock = writerloggermock.NewWriterLoggerMock(&pstrMock, lgrMock)

	errChkMock mock.ErrorCheckerMock

	authSvcMock        mock.AuthServiceMock
	testAuthController AuthController = NewAuthController(&authSvcMock, &errChkMock, &pstrMock, wrlgrMock)

	usrRESTPermSvcMock         userRESTPermSvcMock
	testUserRESTPermController UsersRESTPermissionsController = NewUsersRESTPermissionsController(&usrRESTPermSvcMock, &errChkMock, &pstrMock, wrlgrMock)

	usrSvcMock         usersServiceMock
	testUserController UsersController = NewUsersController(&usrSvcMock, &errChkMock, &pstrMock, wrlgrMock)

	resourcesSvcMock        resourcesServiceMock
	testResourcesController ResourcesController = NewResourcesController(&resourcesSvcMock, &errChkMock, &pstrMock, wrlgrMock)

	httpMethodsSvcMock        httpMethodsServiceMock
	testHttpMethodsController HttpMethodsController = NewHttpMethodsController(&httpMethodsSvcMock, &errChkMock, &pstrMock, wrlgrMock)

	admSvcMock          adminsServiceMock
	testAdminController AdminsController = NewAdminsController(&admSvcMock, &errChkMock, &pstrMock, wrlgrMock)
)

func handlerMockforMDW(rw http.ResponseWriter, req *http.Request) {
	rw.WriteHeader(http.StatusOK)
}

func handlerMockforMDWStatus500(rw http.ResponseWriter, req *http.Request) {
	rw.WriteHeader(http.StatusInternalServerError)
}
