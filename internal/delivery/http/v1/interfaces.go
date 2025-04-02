package v1

import (
	"context"
	"net/http"
	"template/internal/dto/auth"
	dto "template/internal/dto/user"
	authService "template/internal/service/auth"
)

type authSvc interface {
	Login(ctx context.Context, req auth.LoginRequest) (res auth.LoginResponse, err error)
	VerifyToken(userToken string, tokenType authService.TokenType) (userId string, err error)
}

type userSvc interface {
	GetUserById(ctx context.Context, id string) (dto.User, error)
	GetContextUser(ctx context.Context) (*dto.User, error)
	CreatUser(ctx context.Context, request dto.CreateUserRequest) (string, error)
	UpdateUser(ctx context.Context, userID string, upd dto.UpdateUserRequest) error
}
type Responder interface {
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
