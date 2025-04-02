package auth

import (
	"context"
	"template/internal/dto/auth"
	"template/internal/dto/user"
)

type authRepo interface {
	Login(ctx context.Context, req auth.LoginRequest) (user.User, error)
}
