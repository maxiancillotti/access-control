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

type ResourcesController interface {
	// Creates a resource given its Path.
	// Receives body with a dto.Resource.Path field:
	/*
		{
			"path": string
		}

	*/
	// API RESOURCE URL /Resources
	POST(rw http.ResponseWriter, req *http.Request)

	// Deletes a resource given its ID.
	// API RESOURCE URL /Resources/{id}
	DELETE(rw http.ResponseWriter, req *http.Request)

	// Retrieves a resource given its ID.
	// API RESOURCE URL /Resources/{id}
	GET(rw http.ResponseWriter, req *http.Request)

	// Retrieves a resource given its Path.
	// Receives body with a dto.Resource.Path field:
	/*
		{
			"path": string
		}

	*/
	// API RESOURCE URL /Resources
	GETByPath(rw http.ResponseWriter, req *http.Request)

	// Retrieves a collection of all resources.
	// API RESOURCE URL /Resources
	GETCollection(rw http.ResponseWriter, req *http.Request)
	// TO-DO: Filtered search.
}

type resourcesController struct {
	service             services.ResourcesServices
	serviceErrorChecker services.ServiceErrorChecker
	presenter           presenter.Presenter
	writerLogger        responseutil.WriterLogger
}

func NewResourcesController(svc services.ResourcesServices, svcErrorChecker services.ServiceErrorChecker, presenter presenter.Presenter, writerLogger responseutil.WriterLogger) ResourcesController {
	return &resourcesController{
		service:             svc,
		serviceErrorChecker: svcErrorChecker,
		presenter:           presenter,
		writerLogger:        writerLogger,
	}
}

func (c *resourcesController) POST(rw http.ResponseWriter, req *http.Request) {

	var reqBody dto.Resource
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

	resource, err := c.service.Create(reqBody)
	if err != nil {

		var status int
		var respErr error
		localErrMsg := fmt.Sprintf(internal.ErrMsgFmtCreating, "resource")
		privateErr := errors.Wrap(err, localErrMsg)

		if c.serviceErrorChecker.ErrorIsInvalidInputIdentifier(err) {
			// a resource already exists with this path
			status, respErr = http.StatusConflict, privateErr
		} else {
			status, respErr = http.StatusInternalServerError, errors.Wrap(internal.ErrRespInternalUnexpected, localErrMsg)
		}
		c.writerLogger.WriteErrorRespAndLog(req.Context(), rw, status, respErr, privateErr)
		return
	}

	c.presenter.PresentResponse(req.Context(), rw, http.StatusOK, resource)
}

func (c *resourcesController) DELETE(rw http.ResponseWriter, req *http.Request) {

	urlVars := mux.Vars(req)
	id, err := strconv.ParseUint(urlVars["id"], 0, 64)
	if err != nil {
		respErr := errors.Wrapf(err, internal.ErrMsgFmtParsingIDFromURL, "resource")
		c.writerLogger.WriteErrorRespAndLog(req.Context(), rw, http.StatusBadRequest, respErr, respErr) //respErr == logErr
		return
	}

	err = c.service.Delete(uint(id))
	if err != nil {

		var status int
		var respErr error
		localErrMsg := fmt.Sprintf(internal.ErrMsgFmtDeleting, "resource")
		privateErr := errors.Wrap(err, localErrMsg)

		if c.serviceErrorChecker.ErrorIsInvalidInputIdentifier(err) {
			status, respErr = http.StatusNotFound, privateErr
		} else {
			status, respErr = http.StatusInternalServerError, errors.Wrap(internal.ErrRespInternalUnexpected, localErrMsg)
		}
		c.writerLogger.WriteErrorRespAndLog(req.Context(), rw, status, respErr, privateErr)
		return
	}

	c.presenter.SuccessResp(req.Context(), rw, http.StatusOK, "resource deleted OK")
}

func (c *resourcesController) GET(rw http.ResponseWriter, req *http.Request) {

	urlVars := mux.Vars(req)
	id, err := strconv.ParseUint(urlVars["id"], 0, 64)
	if err != nil {
		respErr := errors.Wrapf(err, internal.ErrMsgFmtParsingIDFromURL, "resource")
		c.writerLogger.WriteErrorRespAndLog(req.Context(), rw, http.StatusBadRequest, respErr, respErr) //respErr == logErr
		return
	}

	resource, err := c.service.Retrieve(uint(id))
	if err != nil {

		var status int
		var respErr error
		localErrMsg := fmt.Sprintf(internal.ErrMsgFmtGetting, "resource")
		privateErr := errors.Wrap(err, localErrMsg)

		if c.serviceErrorChecker.ErrorIsInvalidInputIdentifier(err) {
			status, respErr = http.StatusNotFound, privateErr
		} else {
			status, respErr = http.StatusInternalServerError, errors.Wrap(internal.ErrRespInternalUnexpected, localErrMsg)
		}
		c.writerLogger.WriteErrorRespAndLog(req.Context(), rw, status, respErr, privateErr)
		return
	}

	c.presenter.PresentResponse(req.Context(), rw, http.StatusOK, &resource)
}

func (c *resourcesController) GETByPath(rw http.ResponseWriter, req *http.Request) {

	var reqBody dto.Resource
	var err error

	if requestutil.IsDebugDataEnabled(req) {
		reqBodyBytes := requestutil.GetRequestBodyFromCtx(req.Context())
		err = json.NewDecoder(bytes.NewReader(reqBodyBytes)).Decode(&reqBody)
	} else {
		err = json.NewDecoder(req.Body).Decode(&reqBody)
	}
	if err != nil {
		respErr := errors.Wrapf(err, internal.ErrMsgBlockUnmarshalingReqBody, "resource")
		c.writerLogger.WriteErrorRespAndLog(req.Context(), rw, http.StatusBadRequest, respErr, respErr) //respErr == logErr
		return
	}

	if err := reqBody.ValidateFormat(); err != nil {
		respErr := errors.Wrap(err, internal.ErrMsgBlockReqBodyDoesntHaveProperFmt)
		c.writerLogger.WriteErrorRespAndLog(req.Context(), rw, http.StatusUnprocessableEntity, respErr, respErr) //respErr == logErr
		return
	}

	resource, err := c.service.RetrieveByPath(*reqBody.Path)
	if err != nil {

		var status int
		var respErr error
		localErrMsg := fmt.Sprintf(internal.ErrMsgFmtGetting, "resource")
		privateErr := errors.Wrap(err, localErrMsg)

		if c.serviceErrorChecker.ErrorIsInvalidInputIdentifier(err) {
			status, respErr = http.StatusNotFound, privateErr
		} else {
			status, respErr = http.StatusInternalServerError, errors.Wrap(internal.ErrRespInternalUnexpected, localErrMsg)
		}
		c.writerLogger.WriteErrorRespAndLog(req.Context(), rw, status, respErr, privateErr)
		return
	}

	c.presenter.PresentResponse(req.Context(), rw, http.StatusOK, &resource)
}

func (c *resourcesController) GETCollection(rw http.ResponseWriter, req *http.Request) {

	resources, err := c.service.RetrieveAll()
	if err != nil {

		var status int
		var respErr error
		localErrMsg := fmt.Sprintf(internal.ErrMsgFmtGetting, "resources")
		privateErr := errors.Wrap(err, localErrMsg)

		if c.serviceErrorChecker.ErrorIsEmptyResult(err) {
			status, respErr = http.StatusNotFound, privateErr
		} else {
			status, respErr = http.StatusInternalServerError, errors.Wrap(internal.ErrRespInternalUnexpected, localErrMsg)
		}
		c.writerLogger.WriteErrorRespAndLog(req.Context(), rw, status, respErr, privateErr)
		return
	}

	c.presenter.PresentResponse(req.Context(), rw, http.StatusOK, resources)
}
