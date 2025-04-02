package user

import (
	"context"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"template/internal/config"
	dto "template/internal/dto/user"
	"template/internal/entity/user"
	"template/internal/utils"
)

type Service struct {
	userRepo userRepo
	logger   *zap.Logger
}

func NewService(userRepo userRepo, logger *zap.Logger) *Service {
	return &Service{
		userRepo: userRepo,
		logger:   logger,
	}
}

func (s *Service) GetUserById(ctx context.Context, id string) (dto.User, error) {
	return s.userRepo.GetUserById(ctx, id)

}

func (s *Service) GetContextUser(ctx context.Context) (*dto.User, error) {
	if usr, ok := ctx.Value(user.ContextUserKey{}).(*dto.User); ok {
		return usr, nil
	}

	return nil, errors.New("could not find user by contextUserKey{}")
}

func (s *Service) UpdateUser(ctx context.Context, userID string, upd dto.UpdateUserRequest) error {
	if upd.Email != nil {
		usr, err := s.userRepo.GetUserByEmail(ctx, *upd.Email)
		if err == nil && usr.ID != userID {
			return utils.ErrEmailAlreadyExists
		} else if !errors.Is(err, utils.ErrUserNotFound) {
			return err
		}
	}
	return s.userRepo.UpdateUser(ctx, userID, upd)
}

func (s *Service) CreatUser(ctx context.Context, request dto.CreateUserRequest) (string, error) {

	_, err := s.userRepo.GetUserByEmail(ctx, request.Email)
	if err == nil {
		return "", utils.ErrEmailAlreadyExists
	} else if !errors.Is(err, utils.ErrUserNotFound) {
		return "", err
	}

	hashedPassword, err := utils.GeneratePasswordHash(request.Password)
	if err != nil {
		return "", utils.NewInternalError("Can not generate password hash: " + err.Error())
	}
	mongoUser := user.User{
		Name:     request.Name,
		Password: hashedPassword,
		Email:    request.Email,
		Role:     user.Usr,
	}

	userId, err := s.userRepo.CreateUser(ctx, mongoUser)
	if err != nil {
		return "", err
	}

	return userId, nil
}

func (s *Service) InitAdmin(adminCreds *config.AdminConf) error {
	ctx := context.TODO()
	err := s.userRepo.AdminExists(ctx)
	if err != nil && !errors.Is(err, utils.ErrUserNotFound) {
		return err
	} else if err == nil {
		return nil
	}
	hashedPassword, err := utils.GeneratePasswordHash(adminCreds.Password)
	if err != nil {
		return utils.NewInternalError("Can not generate password hash: " + err.Error())
	}
	admin := user.User{
		Name:     adminCreds.Name,
		Email:    adminCreds.Email,
		Password: hashedPassword,
		Role:     user.Admin,
	}
	_, err = s.userRepo.CreateUser(ctx, admin)
	if err != nil {
		return fmt.Errorf("can not create admin: %w", err)
	}
	return nil
}
