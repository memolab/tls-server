package types

import (
	"net/http"

	"go.uber.org/zap"
)

type (

	// Middleware type interface
	Middleware interface {
		Handler() MiddlewareHandler
		Shutdown()
	}

	// MiddlewareHandler middleware http handler type
	MiddlewareHandler func(http.Handler) http.Handler

	// CTXKey for type context withvalue key
	CTXKey string

	// APICTL main api controller type
	APICTL interface {
		RespJSON(http.ResponseWriter, int, interface{})
		Abort(http.ResponseWriter, int)
		Log() *zap.Logger
	}
)
