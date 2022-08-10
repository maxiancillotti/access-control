package presenters

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/maxiancillotti/access-control/internal/controllers/presenter"
	"github.com/maxiancillotti/access-control/internal/dto"

	"github.com/stretchr/testify/assert"
)

var (
	testPresenter presenter.Presenter = NewJSONPresenter()
)

//////////////////////////////////////////////

func TestSuccessResp(t *testing.T) {

	// Initialization
	req := httptest.NewRequest(http.MethodPost, "http://localhost:8001/authtoken", nil)
	rwr := httptest.NewRecorder()

	// Execution
	testPresenter.SuccessResp(req.Context(), rwr, http.StatusOK, "Successful")

	// Check
	statusCode := rwr.Result().StatusCode
	assert.Equal(t, http.StatusOK, statusCode)

	var successResp dto.Success
	err := json.NewDecoder(rwr.Result().Body).Decode(&successResp)
	assert.Nil(t, err)

	assert.Equal(t, "Successful", successResp.Message)
}

func TestErrorResp(t *testing.T) {

	// Initialization
	req := httptest.NewRequest(http.MethodPost, "http://localhost:8001/authtoken", nil)
	rwr := httptest.NewRecorder()

	errMsg := errors.New("ErrorMsg")

	// Execution
	testPresenter.ErrorResp(req.Context(), rwr, http.StatusBadRequest, errMsg)

	// Check
	statusCode := rwr.Result().StatusCode
	assert.Equal(t, http.StatusBadRequest, statusCode)

	var errResp dto.Error
	err := json.NewDecoder(rwr.Result().Body).Decode(&errResp)
	assert.Nil(t, err)

	assert.Equal(t, errMsg.Error(), errResp.Message)
}

func TestPresentResponse(t *testing.T) {

	t.Run("Success", func(t *testing.T) {
		// Initialization
		req := httptest.NewRequest(http.MethodPost, "http://localhost:8001/authtoken", nil)
		rwr := httptest.NewRecorder()

		token := "header.payload.sign"
		tokenResp := dto.TokenResp{
			Token: token,
		}

		// Execution
		testPresenter.PresentResponse(req.Context(), rwr, http.StatusOK, tokenResp)

		// Check
		statusCode := rwr.Result().StatusCode
		assert.Equal(t, http.StatusOK, statusCode)

		var tokenRespBodyCheck dto.TokenResp
		err := json.NewDecoder(rwr.Result().Body).Decode(&tokenRespBodyCheck)
		assert.Nil(t, err)

		assert.Equal(t, token, tokenRespBodyCheck.Token)
	})

	t.Run("Panic", func(t *testing.T) {
		// Initialization
		req := httptest.NewRequest(http.MethodPost, "http://localhost:8001/authtoken", nil)
		rwr := httptest.NewRecorder()

		nonJsonVar := make(chan int) // var that doesn't support JSON encoding

		var panicMsgOut interface{}

		defer func() {
			if panicMsgOut = recover(); panicMsgOut != nil {
				t.Log("Panic recovered ok")

				// Check
				panicErr, ok := panicMsgOut.(error)
				assert.True(t, ok)

				assert.Contains(t, panicErr.Error(), errMsgPanicCannotPresentResponseAsJSON)
			}
		}()

		// Execution
		testPresenter.PresentResponse(req.Context(), rwr, http.StatusOK, nonJsonVar)
	})
}
