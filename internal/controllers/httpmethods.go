package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/maxiancillotti/access-control/internal/controllers/internal"
	"github.com/maxiancillotti/access-control/internal/controllers/presenter"
	"github.com/maxiancillotti/access-control/internal/controllers/requestutil"
	"github.com/maxiancillotti/access-control/internal/controllers/responseutil"
	"github.com/maxiancillotti/access-control/internal/dto"
	"github.com/maxiancillotti/access-control/internal/services"

	"github.com/pkg/errors"
)

type HttpMethodsController interface {
	// Retrieves an http method given its name.
	// Receives a body with the field dto.HttpMethod.Name.
	/*
		{
			"name": string
		}
	*/
	// API RESOURCE URL /HttpMethods
	GETByName(rw http.ResponseWriter, req *http.Request)

	// Retrieves a collection of all http methods.
	// API RESOURCE URL /HttpMethods
	GETCollection(rw http.ResponseWriter, req *http.Request)
}

type httpMethodsController struct {
	service             services.HttpMethodsServices
	serviceErrorChecker services.ServiceErrorChecker
	presenter           presenter.Presenter
	writerLogger        responseutil.WriterLogger
}

func NewHttpMethodsController(svc services.HttpMethodsServices, svcErrorChecker services.ServiceErrorChecker, presenter presenter.Presenter, writerLogger responseutil.WriterLogger) HttpMethodsController {
	return &httpMethodsController{
		service:             svc,
		serviceErrorChecker: svcErrorChecker,
		presenter:           presenter,
		writerLogger:        writerLogger,
	}
}

func (c *httpMethodsController) GETByName(rw http.ResponseWriter, req *http.Request) {

	var reqBody dto.HttpMethod
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

	method, err := c.service.RetrieveByName(*reqBody.Name)
	if err != nil {

		var status int
		var respErr error
		localErrMsg := fmt.Sprintf(internal.ErrMsgFmtGetting, "method")
		privateErr := errors.Wrap(err, localErrMsg)

		if c.serviceErrorChecker.ErrorIsInvalidInputIdentifier(err) {
			status, respErr = http.StatusNotFound, privateErr
		} else {
			status, respErr = http.StatusInternalServerError, errors.Wrap(internal.ErrRespInternalUnexpected, localErrMsg)
		}
		c.writerLogger.WriteErrorRespAndLog(req.Context(), rw, status, respErr, privateErr)
		return
	}

	c.presenter.PresentResponse(req.Context(), rw, http.StatusOK, method)
}

func (c *httpMethodsController) GETCollection(rw http.ResponseWriter, req *http.Request) {

	methods, err := c.service.RetrieveAll()
	if err != nil {

		var status int
		var respErr error
		localErrMsg := fmt.Sprintf(internal.ErrMsgFmtGetting, "methods")
		privateErr := errors.Wrap(err, localErrMsg)

		if c.serviceErrorChecker.ErrorIsEmptyResult(err) {
			status, respErr = http.StatusNotFound, privateErr
		} else {
			status, respErr = http.StatusInternalServerError, errors.Wrap(internal.ErrRespInternalUnexpected, localErrMsg)
		}
		c.writerLogger.WriteErrorRespAndLog(req.Context(), rw, status, respErr, privateErr)
		return
	}

	c.presenter.PresentResponse(req.Context(), rw, http.StatusOK, methods)
}
