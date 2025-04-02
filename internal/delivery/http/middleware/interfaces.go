package middleware

import (
	"context"
	"net/http"
	dto "template/internal/dto/user"
	authSvc "template/internal/service/auth"
)

type authService interface {
	VerifyToken(userToken string, tokenType authSvc.TokenType) (userId string, err error)
}

type userService interface {
	GetUserById(ctx context.Context, id string) (dto.User, error)
	GetContextUser(ctx context.Context) (*dto.User, error)
}

type responder interface {
	With(statusCode int, w http.ResponseWriter, v interface{})
	WithOK(w http.ResponseWriter, v interface{})
	WithNotFound(w http.ResponseWriter, v interface{})
	WithCreated(w http.ResponseWriter, v interface{})
	WithBadRequest(w http.ResponseWriter, message string) string
	WithInternalError(w http.ResponseWriter, message string) string
	WithUnauthorizedError(w http.ResponseWriter) string
	WithForbiddenError(w http.ResponseWriter) string
	WithTooManyRequests(w http.ResponseWriter) string
}
