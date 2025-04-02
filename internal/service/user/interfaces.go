package user

import (
	"context"
	dto "template/internal/dto/user"
	"template/internal/entity/user"
)

type userRepo interface {
	GetUserById(ctx context.Context, userId string) (dto.User, error)
	CreateUser(ctx context.Context, user user.User) (string, error)
	UpdateUser(ctx context.Context, userID string, upd dto.UpdateUserRequest) error
	GetUserByEmail(ctx context.Context, email string) (dto.User, error)
	AdminExists(ctx context.Context) error
}
