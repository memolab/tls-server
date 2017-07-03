package middlewares

import (
	"net/http"

	"tls-server/api/types"
)

// IsAdmin depend on authMiddleware
func IsAdmin(ctl types.APICTL, adminID string) types.MiddlewareHandler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			uid := r.Context().Value(types.CTXKey("uid")).(string)

			if adminID != uid {
				ctl.Abort(rw, http.StatusForbidden)
				return
			}

			next.ServeHTTP(rw, r)
		})
	}
}
