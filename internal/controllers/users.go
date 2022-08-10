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

type UsersController interface {
	// Creates a new user and returns it with the password raw so the user knows
	// what value to use when authenticating.
	// API RESOURCE URL /Users
	// Receives body with field dto.User.Username.
	/*
		{
			"username": string
		}

	*/
	POST(rw http.ResponseWriter, req *http.Request)

	// Updates password given user's ID and returns it raw so the
	// wielder knows what value to use when authenticating.
	// API RESOURCE URL /Users/{id}/Password
	PATCHPassword(rw http.ResponseWriter, req *http.Request)

	// Updates Enabled State given the desired one and the user's ID.
	// API RESOURCE URL /Users/{id}/EnabledState
	// Receives body with field dto.User.EnabledState.
	/*
		{
			"enabled_state": bool
		}

	*/
	PATCHEnabledState(rw http.ResponseWriter, req *http.Request)

	// API RESOURCE URL /Users/{id}
	DELETE(rw http.ResponseWriter, req *http.Request)

	// API RESOURCE URL /Users
	// Receives body with field dto.User.Username.
	/*
		{
			"username": string
		}

	*/
	GETByUsername(rw http.ResponseWriter, req *http.Request)
}

type usersController struct {
	service             services.UsersServices
	serviceErrorChecker services.ServiceErrorChecker
	presenter           presenter.Presenter
	writerLogger        responseutil.WriterLogger
}

func NewUsersController(svc services.UsersServices, svcErrorChecker services.ServiceErrorChecker, presenter presenter.Presenter, writerLogger responseutil.WriterLogger) UsersController {
	return &usersController{
		service:             svc,
		serviceErrorChecker: svcErrorChecker,
		presenter:           presenter,
		writerLogger:        writerLogger,
	}
}

func (c *usersController) POST(rw http.ResponseWriter, req *http.Request) {
	var reqBody dto.User
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

	if err := reqBody.ValidateFormat(); err != nil {
		respErr := errors.Wrap(err, internal.ErrMsgBlockReqBodyDoesntHaveProperFmt)
		c.writerLogger.WriteErrorRespAndLog(req.Context(), rw, http.StatusUnprocessableEntity, respErr, respErr) //respErr == logErr
		return
	}

	user, err := c.service.Create(reqBody)
	if err != nil {

		var status int
		var respErr error
		localErrMsg := fmt.Sprintf(internal.ErrMsgFmtCreating, "user")
		privateErr := errors.Wrap(err, localErrMsg)

		if c.serviceErrorChecker.ErrorIsInvalidInputIdentifier(err) {
			// a resource already exists with this username
			status, respErr = http.StatusConflict, privateErr
		} else {
			status, respErr = http.StatusInternalServerError, errors.Wrap(internal.ErrRespInternalUnexpected, localErrMsg)
		}
		c.writerLogger.WriteErrorRespAndLog(req.Context(), rw, status, respErr, privateErr)
		return
	}

	c.presenter.PresentResponse(req.Context(), rw, http.StatusOK, &user)
}

func (c *usersController) PATCHPassword(rw http.ResponseWriter, req *http.Request) {

	urlVars := mux.Vars(req)
	id, err := strconv.ParseUint(urlVars["id"], 0, 64)
	if err != nil {
		respErr := errors.Wrapf(err, internal.ErrMsgFmtParsingIDFromURL, "user")
		c.writerLogger.WriteErrorRespAndLog(req.Context(), rw, http.StatusBadRequest, respErr, respErr) //respErr == logErr
		return
	}
	idUint := uint(id)

	password, err := c.service.UpdatePassword(idUint)
	if err != nil {

		var status int
		var respErr error
		localErrMsg := fmt.Sprintf(internal.ErrMsgFmtUpdating, "user")
		privateErr := errors.Wrap(err, localErrMsg)

		if c.serviceErrorChecker.ErrorIsInvalidInputIdentifier(err) {
			status, respErr = http.StatusNotFound, privateErr
		} else {
			status, respErr = http.StatusInternalServerError, errors.Wrap(internal.ErrRespInternalUnexpected, localErrMsg)
		}
		c.writerLogger.WriteErrorRespAndLog(req.Context(), rw, status, respErr, privateErr)
		return
	}

	userDTO := dto.User{ID: &idUint, Password: &password}
	c.presenter.PresentResponse(req.Context(), rw, http.StatusOK, &userDTO)
}

func (c *usersController) PATCHEnabledState(rw http.ResponseWriter, req *http.Request) {

	var reqBody dto.User
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

	if err := reqBody.ValidateEmptyEnabledState(); err != nil {
		respErr := errors.Wrap(err, internal.ErrMsgBlockReqBodyDoesntHaveProperFmt)
		c.writerLogger.WriteErrorRespAndLog(req.Context(), rw, http.StatusUnprocessableEntity, respErr, respErr) //respErr == logErr
		return
	}

	urlVars := mux.Vars(req)
	id, err := strconv.ParseUint(urlVars["id"], 0, 64)
	if err != nil {
		respErr := errors.Wrapf(err, internal.ErrMsgFmtParsingIDFromURL, "user")
		c.writerLogger.WriteErrorRespAndLog(req.Context(), rw, http.StatusBadRequest, respErr, respErr) //respErr == logErr
		return
	}
	idUint := uint(id)
	reqBody.ID = &idUint

	err = c.service.UpdateEnabledState(reqBody)
	if err != nil {

		var status int
		var respErr error
		localErrMsg := fmt.Sprintf(internal.ErrMsgFmtUpdating, "user")
		privateErr := errors.Wrap(err, localErrMsg)

		if c.serviceErrorChecker.ErrorIsInvalidInputIdentifier(err) {
			status, respErr = http.StatusNotFound, privateErr
		} else {
			status, respErr = http.StatusInternalServerError, errors.Wrap(internal.ErrRespInternalUnexpected, localErrMsg)
		}
		c.writerLogger.WriteErrorRespAndLog(req.Context(), rw, status, respErr, privateErr)
		return
	}

	c.presenter.SuccessResp(req.Context(), rw, http.StatusOK, "user enabled state updated OK")
}

func (c *usersController) DELETE(rw http.ResponseWriter, req *http.Request) {

	urlVars := mux.Vars(req)
	id, err := strconv.ParseUint(urlVars["id"], 0, 64)
	if err != nil {
		respErr := errors.Wrapf(err, internal.ErrMsgFmtParsingIDFromURL, "user")
		c.writerLogger.WriteErrorRespAndLog(req.Context(), rw, http.StatusBadRequest, respErr, respErr) //respErr == logErr
		return
	}

	err = c.service.Delete(uint(id))
	if err != nil {

		var status int
		var respErr error
		localErrMsg := fmt.Sprintf(internal.ErrMsgFmtDeleting, "user")
		privateErr := errors.Wrap(err, localErrMsg)

		if c.serviceErrorChecker.ErrorIsInvalidInputIdentifier(err) {
			status, respErr = http.StatusNotFound, privateErr
		} else {
			status, respErr = http.StatusInternalServerError, errors.Wrap(internal.ErrRespInternalUnexpected, localErrMsg)
		}
		c.writerLogger.WriteErrorRespAndLog(req.Context(), rw, status, respErr, privateErr)
		return
	}

	c.presenter.SuccessResp(req.Context(), rw, http.StatusOK, "user deleted OK")
}

func (c *usersController) GETByUsername(rw http.ResponseWriter, req *http.Request) {

	var reqBody dto.User
	var err error

	if requestutil.IsDebugDataEnabled(req) {
		reqBodyBytes := requestutil.GetRequestBodyFromCtx(req.Context())
		err = json.NewDecoder(bytes.NewReader(reqBodyBytes)).Decode(&reqBody)
	} else {
		err = json.NewDecoder(req.Body).Decode(&reqBody)
	}
	if err != nil {
		respErr := errors.Wrapf(err, internal.ErrMsgBlockUnmarshalingReqBody, "user")
		c.writerLogger.WriteErrorRespAndLog(req.Context(), rw, http.StatusBadRequest, respErr, respErr) //respErr == logErr
		return
	}

	if err := reqBody.ValidateFormat(); err != nil {
		respErr := errors.Wrap(err, internal.ErrMsgBlockReqBodyDoesntHaveProperFmt)
		c.writerLogger.WriteErrorRespAndLog(req.Context(), rw, http.StatusUnprocessableEntity, respErr, respErr) //respErr == logErr
		return
	}

	user, err := c.service.RetrieveByUsername(*reqBody.Username)
	if err != nil {

		var status int
		var respErr error
		localErrMsg := fmt.Sprintf(internal.ErrMsgFmtGetting, "user")
		privateErr := errors.Wrap(err, localErrMsg)

		if c.serviceErrorChecker.ErrorIsInvalidInputIdentifier(err) {
			status, respErr = http.StatusNotFound, privateErr
		} else {
			status, respErr = http.StatusInternalServerError, errors.Wrap(internal.ErrRespInternalUnexpected, localErrMsg)
		}
		c.writerLogger.WriteErrorRespAndLog(req.Context(), rw, status, respErr, privateErr)
		return
	}

	c.presenter.PresentResponse(req.Context(), rw, http.StatusOK, &user)
}
