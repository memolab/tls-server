package types

import (
	"net/http"

	"go.uber.org/zap"
)

type (

	// Middleware type interface
	Middleware interface {
		Handler() MiddlewareHandler
		Close()
	}

	// MiddlewareHandler middleware http handler type
	MiddlewareHandler func(http.Handler) http.Handler

	// CTXUIDKey for type context withvalue key
	CTXUIDKey struct{}
	// CTXCACHEKey for type context withvalue key
	CTXCACHEKey struct{}

	// APICTL main api controller type
	APICTL interface {
		RespJSON(http.ResponseWriter, int, interface{})
		Abort(http.ResponseWriter, int)
		Log() *zap.Logger
	}
)
