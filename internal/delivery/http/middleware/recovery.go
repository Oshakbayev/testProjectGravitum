package middleware

import (
	"fmt"
	"go.uber.org/zap"
	"net/http"
)

type RecoveryMiddleware struct {
	responder responder
	logger    *zap.Logger
}

func NewRecoveryMiddleware(responder responder, logger *zap.Logger) *RecoveryMiddleware {
	return &RecoveryMiddleware{
		responder: responder,
		logger:    logger,
	}
}

func (rm *RecoveryMiddleware) Recover(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				uuid := rm.responder.WithInternalError(w, "internal error")
				var errMsg string
				if e, ok := err.(error); ok {
					errMsg = e.Error()
				} else {
					errMsg = fmt.Sprintf("%v", err)
				}
				rm.logger.Error(
					"panic recovered",
					zap.String("error_uuid", uuid),
					zap.String("error", errMsg))

			}
		}()

		next.ServeHTTP(w, r)
	})
}
