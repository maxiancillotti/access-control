package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

// Call before passing httpRouter to an http server.
// Always add routes from more to less specific.
// URLs must be lowercase so CaselessMatcher wrapper works.

// Authentication vía "Authorization" header with custom type "X-Admin" is necessary.
func SetRoutesAdmins() {

	// Updates Enabled State given the desired one and the admin's ID.
	// Receives body with field dto.Admin.EnabledState.
	/*
		{
			"enabled_state": bool
		}

	*/
	httpRouter.HandleFunc("/api/admins/{id:[0-9]+}/enabled-state",
		loggingMDW.LogInAndOut(
			panicRecoverMDW.PanicRecover(
				presenterMDW.CheckPresentationHeaders(
					authController.AuthorizationXAdmin(
						adminsController.PATCHEnabledState))))).
		Methods(http.MethodPatch)

	// Updates password given admin's ID and returns it raw so the
	// wielder knows what value to use when authenticating.
	httpRouter.HandleFunc("/api/admins/{id:[0-9]+}/password",
		loggingMDW.LogInAndOut(
			panicRecoverMDW.PanicRecover(
				presenterMDW.CheckPresentationHeaders(
					authController.AuthorizationXAdmin(
						adminsController.PATCHPassword))))).
		Methods(http.MethodPatch)

	// Deletes an admin given its ID
	httpRouter.HandleFunc("/api/admins/{id:[0-9]+}",
		loggingMDW.LogInAndOut(
			panicRecoverMDW.PanicRecover(
				presenterMDW.CheckPresentationHeaders(
					authController.AuthorizationXAdmin(
						adminsController.DELETE))))).
		Methods(http.MethodDelete)

	// Creates a new admin
	// Receives body with field "username".
	/*
		{
			"username": string
		}

	*/
	httpRouter.HandleFunc("/api/admins",
		loggingMDW.LogInAndOut(
			panicRecoverMDW.PanicRecover(
				presenterMDW.CheckPresentationHeaders(
					authController.AuthorizationXAdmin(
						adminsController.POST))))).
		Methods(http.MethodPost)

	// Retrieves a user indicating its username.
	// Receives body with field "username".
	/*
		{
			"username": string
		}

	*/
	httpRouter.HandleFunc("/api/admins",
		loggingMDW.LogInAndOut(
			panicRecoverMDW.PanicRecover(
				presenterMDW.CheckPresentationHeaders(
					authController.AuthorizationXAdmin(
						adminsController.GETByUsername))))).
		Methods(http.MethodGet)
}

func SetRoutesUsersToken() {

	// Authorices access to a resource for a given token.
	// Receives dto.AuthorizationRequest body with fields:
	/*
		{
		    "token": string,
		    "resource_requested": string,
		    "method_requested": string
		}
	*/
	httpRouter.HandleFunc("/api/users/token/authorize",
		loggingMDW.LogInAndOut(
			panicRecoverMDW.PanicRecover(
				presenterMDW.CheckPresentationHeaders(
					authController.PostAuthTokenAuthorize)))).
		Methods(http.MethodPost)

	// Creates a new access authorization token for a user
	// authenticated vía Authorization Basic.
	httpRouter.HandleFunc("/api/users/token",
		loggingMDW.LogInAndOut(
			panicRecoverMDW.PanicRecover(
				presenterMDW.CheckPresentationHeaders(
					authController.AuthorizationBasic(
						authController.PostAuthToken))))).
		Methods(http.MethodPost)
}

// Authentication vía "Authorization" header with custom type "X-Admin" is necessary.
func SetRoutesUsers() {

	// Retrieves all REST Permissions associated to the given userID, with all its descriptions.
	httpRouter.HandleFunc("/api/users/{userID:[0-9]+}/rest-permissions-with-descriptions",
		loggingMDW.LogInAndOut(
			panicRecoverMDW.PanicRecover(
				presenterMDW.CheckPresentationHeaders(
					authController.AuthorizationXAdmin(
						usersRESTPermissionsController.GETCollectionWithDescriptionsByUserID))))).
		Methods(http.MethodGet)

	// Retrieves all REST Permissions associated to the given userID.
	httpRouter.HandleFunc("/api/users/{userID:[0-9]+}/rest-permissions",
		loggingMDW.LogInAndOut(
			panicRecoverMDW.PanicRecover(
				presenterMDW.CheckPresentationHeaders(
					authController.AuthorizationXAdmin(
						usersRESTPermissionsController.GETCollectionByUserID))))).
		Methods(http.MethodGet)

	// Updates Enabled State given the desired one and the user's ID.
	// Receives body with field dto.User.EnabledState.
	/*
		{
			"enabled_state": bool
		}

	*/
	httpRouter.HandleFunc("/api/users/{id:[0-9]+}/enabled-state",
		loggingMDW.LogInAndOut(
			panicRecoverMDW.PanicRecover(
				presenterMDW.CheckPresentationHeaders(
					authController.AuthorizationXAdmin(
						usersController.PATCHEnabledState))))).
		Methods(http.MethodPatch)

	// Updates password given user's ID and returns it raw so the
	// wielder knows what value to use when authenticating.
	httpRouter.HandleFunc("/api/users/{id:[0-9]+}/password",
		loggingMDW.LogInAndOut(
			panicRecoverMDW.PanicRecover(
				presenterMDW.CheckPresentationHeaders(
					authController.AuthorizationXAdmin(
						usersController.PATCHPassword))))).
		Methods(http.MethodPatch)

	// Deletes a user given its ID
	httpRouter.HandleFunc("/api/users/{id:[0-9]+}",
		loggingMDW.LogInAndOut(
			panicRecoverMDW.PanicRecover(
				presenterMDW.CheckPresentationHeaders(
					authController.AuthorizationXAdmin(
						usersController.DELETE))))).
		Methods(http.MethodDelete)

	// Creates a new user
	// Receives body with field "username".
	/*
		{
			"username": string
		}

	*/
	httpRouter.HandleFunc("/api/users",
		loggingMDW.LogInAndOut(
			panicRecoverMDW.PanicRecover(
				presenterMDW.CheckPresentationHeaders(
					authController.AuthorizationXAdmin(
						usersController.POST))))).
		Methods(http.MethodPost)

	// Retrieves a user indicating its username.
	// Receives body with field "username".
	/*
		{
			"username": string
		}

	*/
	httpRouter.HandleFunc("/api/users",
		loggingMDW.LogInAndOut(
			panicRecoverMDW.PanicRecover(
				presenterMDW.CheckPresentationHeaders(
					authController.AuthorizationXAdmin(
						usersController.GETByUsername))))).
		Methods(http.MethodGet)
}

// Authentication vía "Authorization" header with custom type "X-Admin" is necessary.
func SetRoutesUsersRESTPermissions() {

	// Deletes a REST Permission associated to a UserID.
	// Receives dto.UserRESTPermission body with fields:
	/*
		{
			"user_id": uint,
			"permission": {
				"resource_id": uint,
				"method_id": uint
			}
		}

	*/
	httpRouter.HandleFunc("/api/users-rest-permissions",
		loggingMDW.LogInAndOut(
			panicRecoverMDW.PanicRecover(
				presenterMDW.CheckPresentationHeaders(
					authController.AuthorizationXAdmin(
						usersRESTPermissionsController.DELETE))))).
		Methods(http.MethodDelete)

	// Creates a REST Permission associated to a UserID.
	// Receives dto.UserRESTPermission body with fields:
	/*
		{
			"user_id": uint,
			"permission": {
				"resource_id": uint,
				"method_id": uint
			}
		}

	*/
	httpRouter.HandleFunc("/api/users-rest-permissions",
		loggingMDW.LogInAndOut(
			panicRecoverMDW.PanicRecover(
				presenterMDW.CheckPresentationHeaders(
					authController.AuthorizationXAdmin(
						usersRESTPermissionsController.POST))))).
		Methods(http.MethodPost)
}

// Authentication vía "Authorization" header with custom type "X-Admin" is necessary.
func SetRoutesResources() {

	// Deletes a resource given its ID.
	httpRouter.HandleFunc("/api/resources/{id:[0-9]+}",
		loggingMDW.LogInAndOut(
			panicRecoverMDW.PanicRecover(
				presenterMDW.CheckPresentationHeaders(
					authController.AuthorizationXAdmin(
						resourcesController.DELETE))))).
		Methods(http.MethodDelete)

	// Retrieves a resource given its ID.
	httpRouter.HandleFunc("/api/resources/{id:[0-9]+}",
		loggingMDW.LogInAndOut(
			panicRecoverMDW.PanicRecover(
				presenterMDW.CheckPresentationHeaders(
					authController.AuthorizationXAdmin(
						resourcesController.GET))))).
		Methods(http.MethodGet)

	// Retrieves a resource given its path.
	// Receives body with a dto.Resource.Path field:
	/*
		{
			"path": string
		}

	*/
	httpRouter.HandleFunc("/api/resources",
		loggingMDW.LogInAndOut(
			panicRecoverMDW.PanicRecover(
				presenterMDW.CheckPresentationHeaders(
					authController.AuthorizationXAdmin(
						resourcesController.GETByPath))))).
		Methods(http.MethodGet).
		MatcherFunc(func(r *http.Request, rm *mux.RouteMatch) bool {
			return r.ContentLength > 0
		})

	// Retrieves a collection of all resources.
	// TO-DO: Filtered search.
	httpRouter.HandleFunc("/api/resources",
		loggingMDW.LogInAndOut(
			panicRecoverMDW.PanicRecover(
				presenterMDW.CheckPresentationHeaders(
					authController.AuthorizationXAdmin(
						resourcesController.GETCollection))))).
		Methods(http.MethodGet).
		MatcherFunc(func(r *http.Request, rm *mux.RouteMatch) bool {
			return r.ContentLength <= 0
		})

	// Creates a resource given its Path.
	// Receives body with a dto.Resource.Path field:
	/*
		{
			"path": string
		}

	*/
	httpRouter.HandleFunc("/api/resources",
		loggingMDW.LogInAndOut(
			panicRecoverMDW.PanicRecover(
				presenterMDW.CheckPresentationHeaders(
					authController.AuthorizationXAdmin(
						resourcesController.POST))))).
		Methods(http.MethodPost)
}

// Authentication vía "Authorization" header with custom type "X-Admin" is necessary.
func SetRoutesHttpMethods() {

	// Retrieves an http method given its name.
	// Receives a body with the field dto.HttpMethod.Name.
	/*
		{
			"name": string
		}
	*/
	httpRouter.HandleFunc("/api/http-methods",
		loggingMDW.LogInAndOut(
			panicRecoverMDW.PanicRecover(
				presenterMDW.CheckPresentationHeaders(
					authController.AuthorizationXAdmin(
						httpMethodsController.GETByName))))).
		Methods(http.MethodGet).
		MatcherFunc(func(r *http.Request, rm *mux.RouteMatch) bool {
			return r.ContentLength > 0
		})

	// Retrieves a collection of all http methods.
	httpRouter.HandleFunc("/api/http-methods",
		loggingMDW.LogInAndOut(
			panicRecoverMDW.PanicRecover(
				presenterMDW.CheckPresentationHeaders(
					authController.AuthorizationXAdmin(
						httpMethodsController.GETCollection))))).
		Methods(http.MethodGet).
		MatcherFunc(func(r *http.Request, rm *mux.RouteMatch) bool {
			return r.ContentLength <= 0
		})
}

func SetRoutesBase() {
	// Health
	httpRouter.HandleFunc("/", func(resp http.ResponseWriter, req *http.Request) {
		fmt.Fprintln(resp, "Up and running...")
	}).Methods(http.MethodGet)
}

// Wrap httpRouter before starting the server
func CaselessMatcher(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.URL.Path = strings.ToLower(r.URL.Path)
		next.ServeHTTP(w, r)
	})
}
