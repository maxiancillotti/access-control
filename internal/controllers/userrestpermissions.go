package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/maxiancillotti/access-control/internal/controllers/internal"
	"github.com/maxiancillotti/access-control/internal/controllers/presenter"
	"github.com/maxiancillotti/access-control/internal/controllers/requestutil"
	"github.com/maxiancillotti/access-control/internal/controllers/responseutil"
	"github.com/maxiancillotti/access-control/internal/dto"
	"github.com/maxiancillotti/access-control/internal/services"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"
)

type UsersRESTPermissionsController interface {

	// API RESOURCE URL /UsersRESTPermissions

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
	POST(rw http.ResponseWriter, req *http.Request)

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
	DELETE(rw http.ResponseWriter, req *http.Request)

	/**************************************************/

	// Intersection of Resources /Users and /UsersRESTPermissions
	// API RESOURCE URL /Users/{userID}/RESTPermissions/{Attribute to refer a specific subject}

	// Retrieves all REST Permissions associated to the given userID.
	// API RESOURCE /Users/{userID}/RESTPermissions
	GETCollectionByUserID(rw http.ResponseWriter, req *http.Request)

	// Retrieves all REST Permissions associated to the given userID, with all its descriptions.
	// API RESOURCE /Users/{userID}/RESTPermissionsWithDescriptions
	GETCollectionWithDescriptionsByUserID(rw http.ResponseWriter, req *http.Request)

	// API RESOURCE /Users/{userID}/RESTPermissions/PathsMethods
	// This shoulnd't by exposed. The service method retrieveAllPathMethodsByUserID
	// is for internal calls
	//GETPathMethodsCollectionByUserID(rw http.ResponseWriter, req *http.Request)
}

type usersRESTPermissionsController struct {
	service             services.UsersRESTPermissionsServices
	serviceErrorChecker services.ServiceErrorChecker
	presenter           presenter.Presenter
	writerLogger        responseutil.WriterLogger
}

func NewUsersRESTPermissionsController(svc services.UsersRESTPermissionsServices, svcErrorChecker services.ServiceErrorChecker, presenter presenter.Presenter, writerLogger responseutil.WriterLogger) UsersRESTPermissionsController {
	return &usersRESTPermissionsController{
		service:             svc,
		serviceErrorChecker: svcErrorChecker,
		presenter:           presenter,
		writerLogger:        writerLogger,
	}
}

func (c *usersRESTPermissionsController) POST(rw http.ResponseWriter, req *http.Request) {
	var reqBody dto.UserRESTPermission
	var err error

	if requestutil.IsDebugDataEnabled(req) {
		reqBodyBytes := requestutil.GetRequestBodyFromCtx(req.Context())
		err = json.NewDecoder(bytes.NewReader(reqBodyBytes)).Decode(&reqBody)
	} else {
		err = json.NewDecoder(req.Body).Decode(&reqBody)
	}
	if err != nil {
		respErr := errors.Wrap(err, internal.ErrMsgBlockUnmarshalingReqBody)
		c.writerLogger.WriteErrorRespAndLog(req.Context(), rw, http.StatusBadRequest, respErr, respErr) //respErr == logErr
		return
	}

	if err = reqBody.ValidateEmpty(); err != nil {
		respErr := errors.Wrap(err, internal.ErrMsgBlockReqBodyDoesntHaveProperFmt)
		c.writerLogger.WriteErrorRespAndLog(req.Context(), rw, http.StatusBadRequest, respErr, respErr) //respErr == logErr
		return
	}

	err = c.service.Create(reqBody)
	if err != nil {

		var status int
		var respErr error
		localErrMsg := fmt.Sprintf(internal.ErrMsgFmtCreating, "user REST permissions")
		privateErr := errors.Wrap(err, localErrMsg)

		if c.serviceErrorChecker.ErrorIsInvalidInputIdentifier(err) {
			// a resource already exists with this combination of IDs
			status, respErr = http.StatusConflict, privateErr
		} else {
			status, respErr = http.StatusInternalServerError, errors.Wrap(internal.ErrRespInternalUnexpected, localErrMsg)
		}
		c.writerLogger.WriteErrorRespAndLog(req.Context(), rw, status, respErr, privateErr)
		return
	}

	c.presenter.PresentResponse(req.Context(), rw, http.StatusOK, &reqBody)
}

func (c *usersRESTPermissionsController) DELETE(rw http.ResponseWriter, req *http.Request) {

	var reqBody dto.UserRESTPermission
	var err error

	if requestutil.IsDebugDataEnabled(req) {
		reqBodyBytes := requestutil.GetRequestBodyFromCtx(req.Context())
		err = json.NewDecoder(bytes.NewReader(reqBodyBytes)).Decode(&reqBody)
	} else {
		err = json.NewDecoder(req.Body).Decode(&reqBody)
	}
	if err != nil {
		respErr := errors.Wrap(err, internal.ErrMsgBlockUnmarshalingReqBody)
		c.writerLogger.WriteErrorRespAndLog(req.Context(), rw, http.StatusBadRequest, respErr, respErr) //respErr == logErr
		return
	}

	if err = reqBody.ValidateEmpty(); err != nil {
		respErr := errors.Wrap(err, internal.ErrMsgBlockReqBodyDoesntHaveProperFmt)
		c.writerLogger.WriteErrorRespAndLog(req.Context(), rw, http.StatusBadRequest, respErr, respErr) //respErr == logErr
		return
	}

	err = c.service.Delete(reqBody)
	if err != nil {

		var status int
		var respErr error
		localErrMsg := fmt.Sprintf(internal.ErrMsgFmtDeleting, "user REST permissions")
		privateErr := errors.Wrap(err, localErrMsg)

		if c.serviceErrorChecker.ErrorIsInvalidInputIdentifier(err) {
			status, respErr = http.StatusNotFound, privateErr
		} else {
			status, respErr = http.StatusInternalServerError, errors.Wrap(internal.ErrRespInternalUnexpected, localErrMsg)
		}
		c.writerLogger.WriteErrorRespAndLog(req.Context(), rw, status, respErr, privateErr)
		return
	}

	c.presenter.SuccessResp(req.Context(), rw, http.StatusOK, "permissions deleted OK")
}

func (c *usersRESTPermissionsController) GETCollectionByUserID(rw http.ResponseWriter, req *http.Request) {

	urlVars := mux.Vars(req)
	userID, err := strconv.ParseUint(urlVars["userID"], 0, 64)
	if err != nil {
		respErr := errors.Wrapf(err, internal.ErrMsgFmtParsingIDFromURL, "user")
		c.writerLogger.WriteErrorRespAndLog(req.Context(), rw, http.StatusBadRequest, respErr, respErr) //respErr == logErr
		return
	}

	permissions, err := c.service.RetrieveAllByUserID(uint(userID))
	if err != nil {

		var status int
		var respErr error
		localErrMsg := fmt.Sprintf(internal.ErrMsgFmtGetting, "user REST permissions")
		privateErr := errors.Wrap(err, localErrMsg)

		if c.serviceErrorChecker.ErrorIsEmptyResult(err) {
			status, respErr = http.StatusNotFound, privateErr
		} else {
			status, respErr = http.StatusInternalServerError, errors.Wrap(internal.ErrRespInternalUnexpected, localErrMsg)
		}
		c.writerLogger.WriteErrorRespAndLog(req.Context(), rw, status, respErr, privateErr)
		return
	}

	c.presenter.PresentResponse(req.Context(), rw, http.StatusOK, permissions)
}

func (c *usersRESTPermissionsController) GETCollectionWithDescriptionsByUserID(rw http.ResponseWriter, req *http.Request) {

	urlVars := mux.Vars(req)
	userID, err := strconv.ParseUint(urlVars["userID"], 0, 64)
	if err != nil {
		respErr := errors.Wrapf(err, internal.ErrMsgFmtParsingIDFromURL, "user")
		c.writerLogger.WriteErrorRespAndLog(req.Context(), rw, http.StatusBadRequest, respErr, respErr) //respErr == logErr
		return
	}

	permissions, err := c.service.RetrieveAllWithDescriptionsByUserID(uint(userID))
	if err != nil {

		var status int
		var respErr error
		localErrMsg := fmt.Sprintf(internal.ErrMsgFmtGetting, "user REST permissions")
		privateErr := errors.Wrap(err, localErrMsg)

		if c.serviceErrorChecker.ErrorIsEmptyResult(err) {
			status, respErr = http.StatusNotFound, privateErr
		} else {
			status, respErr = http.StatusInternalServerError, errors.Wrap(internal.ErrRespInternalUnexpected, localErrMsg)
		}
		c.writerLogger.WriteErrorRespAndLog(req.Context(), rw, status, respErr, privateErr)
		return
	}

	c.presenter.PresentResponse(req.Context(), rw, http.StatusOK, permissions)
}
