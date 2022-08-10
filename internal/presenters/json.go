package presenters

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/maxiancillotti/access-control/internal/controllers/presenter"
	"github.com/maxiancillotti/access-control/internal/dto"
	"github.com/pkg/errors"
)

const (
	contentTypeJSON = "application/json"
	acceptHeaderKey = "Accept"
)

type jsonPresenter struct{}

func NewJSONPresenter() presenter.Presenter {
	return &jsonPresenter{}
}

// PRESENTER

func (p *jsonPresenter) SuccessResp(ctx context.Context, rw http.ResponseWriter, status int, successMsg string) {

	responseBody := dto.Success{Message: successMsg}
	p.PresentResponse(ctx, rw, status, responseBody)
}

func (p *jsonPresenter) ErrorResp(ctx context.Context, rw http.ResponseWriter, status int, errorResp error) {

	responseBody := dto.Error{Message: errorResp.Error()}
	p.PresentResponse(ctx, rw, status, responseBody)
}

func (p *jsonPresenter) PresentResponse(ctx context.Context, rw http.ResponseWriter, status int, responseBody interface{}) {

	// Header first because Write method without calling WriteHeader
	// will write status 200 automatically.
	rw.Header().Set("Content-Type", contentTypeJSON)
	rw.WriteHeader(status)

	err := json.NewEncoder(rw).Encode(responseBody)
	if err != nil {
		panic(errors.Wrap(err, errMsgPanicCannotPresentResponseAsJSON))
	}

	// WHAT HAPPENS TO PREVIOUS OPERATION?
	// Rollback? Client knows with http status 500 that repeating operation can duplicate things?
	// This doesn't apply to this microservice, though. Because operations cannot be rolledback.
	// It cannot do any harm to generate another token or validate it again anyways.
}
