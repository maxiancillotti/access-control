package main

import (
	"github.com/maxiancillotti/access-control/app"
	"github.com/maxiancillotti/access-control/app/config"

	"github.com/maxiancillotti/access-control/internal/controllers"
	"github.com/maxiancillotti/access-control/internal/controllers/middlewares"
	"github.com/maxiancillotti/access-control/internal/controllers/responseutil"
	"github.com/maxiancillotti/access-control/internal/dataaccess"
	"github.com/maxiancillotti/access-control/internal/presenters"

	"github.com/maxiancillotti/access-control/internal/services"
	"github.com/maxiancillotti/access-control/internal/services/authtoken"
	"github.com/maxiancillotti/access-control/internal/services/errchk"

	"github.com/maxiancillotti/logger"

	"github.com/gorilla/mux"
)

const (
	configFileDirName        = "access-control-api-rest"
	enableStrictSlashRouting = true
)

var (
	// CONFIG

	configData = config.GetConfig(configFileDirName)

	// DAL

	dbconn = app.GetDBConn(&configData.Database)

	adminsDAL          = dataaccess.NewAdminsDAL(dbconn)
	userDAL            = dataaccess.NewUserDAL(dbconn)
	userPermissionsDAL = dataaccess.NewUserRESTPermissionsDAL(dbconn)
	resourcesDAL       = dataaccess.NewResourcesDAL(dbconn)
	httpMethodsDAL     = dataaccess.NewHttpMethodsDAL(dbconn)

	// SERVICES

	adminsSvc               = services.NewAdminsServices(adminsDAL)
	usersSvc                = services.NewUsersServices(userDAL)
	usersRESTPermissionsSvc = services.NewUsersRESTPermissionsServices(userPermissionsDAL, usersSvc)
	resourcesSvc            = services.NewResourcesServices(resourcesDAL)
	httpMethodsSvc          = services.NewHttpMethodsServices(httpMethodsDAL)

	authTokenConfig   = app.GetServiceConfig(&configData.Service)
	authTokenServices = authtoken.NewJWTServices(authTokenConfig)
	authServices      = services.NewAuthServices(usersSvc, adminsSvc, usersRESTPermissionsSvc, authTokenServices)
	svcErrChecker     = errchk.NewServiceErrorChecker()

	// HANDLERS

	globalLogger = logger.New()

	jsonPresenter = presenters.NewJSONPresenter()
	presenterMDW  = presenters.NewJSONPresenterMDW(writerLogger)

	writerLogger = responseutil.NewWriterLogger(jsonPresenter, globalLogger)

	authController = controllers.NewAuthController(authServices, svcErrChecker, jsonPresenter, writerLogger)

	adminsController = controllers.NewAdminsController(adminsSvc, svcErrChecker, jsonPresenter, writerLogger)

	usersController                = controllers.NewUsersController(usersSvc, svcErrChecker, jsonPresenter, writerLogger)
	usersRESTPermissionsController = controllers.NewUsersRESTPermissionsController(usersRESTPermissionsSvc, svcErrChecker, jsonPresenter, writerLogger)
	resourcesController            = controllers.NewResourcesController(resourcesSvc, svcErrChecker, jsonPresenter, writerLogger)
	httpMethodsController          = controllers.NewHttpMethodsController(httpMethodsSvc, svcErrChecker, jsonPresenter, writerLogger)

	panicRecoverMDW = middlewares.NewPanicRecoverMiddleware(writerLogger)
	loggingMDW      = middlewares.NewLoggingMiddleware(globalLogger)

	// ROUTES

	httpRouter = mux.NewRouter().StrictSlash(enableStrictSlashRouting)
)

func main() {
	Start()
}

func Start() {
	defer dbconn.Close()

	SetRoutesAdmins()

	SetRoutesUsersToken()

	SetRoutesUsers()
	SetRoutesUsersRESTPermissions()
	SetRoutesResources()
	SetRoutesHttpMethods()

	SetRoutesBase()

	app.StartHttpServer(&configData.HttpServer, CaselessMatcher(httpRouter))
}
