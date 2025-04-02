package middleware

import (
	"context"
	"errors"
	"go.uber.org/zap"
	"net/http"
	"strings"
	"template/internal/entity/user"
	"template/internal/service/auth"
	"template/internal/utils"
)

const AuthorizationHeaderKey = "Authorization"

type AuthMiddleware struct {
	userService userService
	authService authService
	responder   responder
	logger      *zap.Logger
}

func NewAuthMiddleware(userService userService, authService authService, responder responder, logger *zap.Logger) *AuthMiddleware {
	return &AuthMiddleware{
		userService: userService,
		authService: authService,
		responder:   responder,
		logger:      logger,
	}
}

func (c *AuthMiddleware) AuthorizedJWT(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, ok := r.Header[AuthorizationHeaderKey]; !ok {
			uuid := c.responder.WithUnauthorizedError(w)
			c.logger.Warn("'Authorization' key missing from headers", zap.String("error_uuid", uuid))
			return
		}

		jwtToken, ok := getTokenFromHeader(r)
		if !ok {
			uuid := c.responder.WithUnauthorizedError(w)
			c.logger.Warn("Failed to get token from header", zap.String("error_uuid", uuid))
			return
		}

		userId, err := c.authService.VerifyToken(jwtToken, auth.TokenTypeAccess)

		var expError utils.ExplorerError
		if errors.As(err, &expError) {
			switch expError.ErrorCategory {
			case utils.ErrorCategoryInternal:
				uuid := c.responder.WithInternalError(w, "internal error")
				c.logger.Error("failed authService.VerifyToken", zap.String("error_uuid", uuid), zap.String("reason", err.Error()))
			case utils.ErrorCategoryUserError:
				c.responder.With(expError.StatusCode, w, expError.Message)
			}

			return
		} else if err != nil {
			uuid := c.responder.WithInternalError(w, "internal error")
			c.logger.Error("failed authService.VerifyToken", zap.String("error_uuid", uuid), zap.String("reason", err.Error()))
			return
		}

		u, err := c.userService.GetUserById(r.Context(), userId)
		if errors.As(err, &expError) {
			switch expError.ErrorCategory {
			case utils.ErrorCategoryInternal:
				uuid := c.responder.WithInternalError(w, "internal error")
				c.logger.Error("failed userService.GetUserById", zap.String("error_uuid", uuid), zap.String("reason", err.Error()))
			case utils.ErrorCategoryUserError:
				c.responder.With(expError.StatusCode, w, expError.Message)
			}

			return
		} else if err != nil {
			uuid := c.responder.WithInternalError(w, "internal error")
			c.logger.Error("failed userService.GetUserById", zap.String("error_uuid", uuid), zap.String("reason", err.Error()))
			return
		}

		newCtx := context.WithValue(r.Context(), user.ContextUserKey{}, &u)
		next.ServeHTTP(w, r.WithContext(newCtx))
	})
}

func (c *AuthMiddleware) AuthorizedAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctxUser, err := c.userService.GetContextUser(r.Context())
		if err != nil {
			c.responder.WithInternalError(w, "internal error")
			return
		}

		if ctxUser.Role != string(user.Admin) {
			c.responder.WithForbiddenError(w)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func getTokenFromHeader(r *http.Request) (string, bool) {
	bearer := r.Header.Get(AuthorizationHeaderKey)
	if len(bearer) > 7 && strings.ToUpper(bearer[0:6]) == "BEARER" {
		return bearer[7:], true
	} else if len(bearer) > 6 && strings.ToUpper(bearer[0:5]) == "BASIC" {
		return bearer[6:], true
	}

	return "", false
}
