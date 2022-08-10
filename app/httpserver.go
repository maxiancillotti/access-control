package app

import (
	"fmt"
	"net/http"
	"time"

	"github.com/maxiancillotti/access-control/app/config"
	"github.com/maxiancillotti/gohttpserver"
)

func StartHttpServer(config *config.HttpServerConfig, handlers http.Handler) {

	gohttpserver.NewBuilder().
		SetReadHeaderTimeout(time.Duration(config.ReadHeaderTimeout.GetDuration())).
		SetWriteTimeout(time.Duration(config.WriteTimeout.GetDuration())).
		SetAddr(fmt.Sprintf("%s:%d", config.Hostname, config.HostPort)).
		Build(handlers).
		ListenAndServe()
}
