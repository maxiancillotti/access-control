package presenters

import "github.com/pkg/errors"

var (
	errRespCheckHeadersContentTypeNotSupported = errors.New("the accepted content type for the response is not supported, JSON only")

	errMsgPanicCannotPresentResponseAsJSON = "cannot write present response body as a JSON"
)
