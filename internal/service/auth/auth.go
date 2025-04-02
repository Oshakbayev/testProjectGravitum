package auth

import (
	"context"
	"crypto/ed25519"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"template/internal/config"
	"template/internal/dto/auth"
	"template/internal/utils"
	"time"
)

type TokenType string

const (
	TokenTypeAccess          TokenType = "access"
	TokenTypeRefresh         TokenType = "refresh"
	defaultAccessExpiration            = 15 * time.Minute
	defaultRefreshExpiration           = 24 * time.Hour
)

type Service struct {
	authRepo authRepo
	cfg      *config.AuthConf
	logger   *zap.Logger
}

func NewService(authRepo authRepo, cfg *config.AuthConf, logger *zap.Logger) *Service {
	if cfg.AccessKeyExpiration == 0 {
		cfg.AccessKeyExpiration = defaultAccessExpiration
	}

	if cfg.RefreshKeyExpiration == 0 {
		cfg.RefreshKeyExpiration = defaultRefreshExpiration
	}
	return &Service{
		authRepo: authRepo,
		cfg:      cfg,
		logger:   logger,
	}
}

func (s *Service) Login(ctx context.Context, req auth.LoginRequest) (res auth.LoginResponse, err error) {
	user, err := s.authRepo.Login(ctx, req)
	if errors.Is(err, utils.ErrUserNotFound) {
		return res, utils.ErrWrongCredentials
	} else if err != nil {
		return res, utils.NewInternalError(err.Error())
	}

	res, err = s.GenerateTokensByUserId(user.ID)
	if err != nil {
		return res, utils.NewInternalError(err.Error())
	}

	res.User = &user
	return res, nil
}

func (s *Service) VerifyToken(userToken string, tokenType TokenType) (userId string, err error) {
	token, err := jwt.Parse(userToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodEd25519); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return ed25519.PublicKey(s.cfg.PublicKeyByte), nil
	})

	var validationErr jwt.ValidationError
	if errors.As(err, &validationErr) {
		if validationErr.Errors != jwt.ValidationErrorExpired {
			s.logger.Error("failed to validate token", zap.String("reason", err.Error()))
		}

		return userId, utils.ErrInvalidToken
	} else if err != nil {
		s.logger.Error("failed to parse token", zap.String("reason", err.Error()))
		return userId, utils.ErrInvalidToken
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		s.logger.Error(fmt.Sprintf("unexpected type %T", claims))
		return userId, utils.ErrInvalidToken
	}

	reqTokenType, ok := claims["token_type"].(string)
	if !ok || reqTokenType != string(tokenType) {
		return userId, utils.ErrUnauthorized
	}

	if !token.Valid {
		s.logger.Error("invalid token")
		return userId, utils.ErrInvalidToken
	}

	return getUserIdFromJwt(claims)
}

func (s *Service) GenerateToken(userId string, tokenType TokenType, expiration time.Duration) (string, error) {
	token := jwt.NewWithClaims(
		jwt.SigningMethodEdDSA,
		jwt.MapClaims{
			"user_id":    userId,
			"exp":        time.Now().Add(expiration).Unix(),
			"iat":        time.Now().Unix(),
			"jti":        uuid.New().String(),
			"token_type": tokenType,
		},
	)

	return token.SignedString(ed25519.PrivateKey(s.cfg.PrivateKeyByte))
}

func (s *Service) GenerateTokensByUserId(userId string) (res auth.LoginResponse, err error) {
	res.AccessToken, err = s.GenerateToken(userId, TokenTypeAccess, s.cfg.AccessKeyExpiration)
	if err != nil {
		return res, err
	}

	res.RefreshToken, err = s.GenerateToken(userId, TokenTypeRefresh, s.cfg.RefreshKeyExpiration)
	if err != nil {
		return res, err
	}

	return res, nil
}

func getUserIdFromJwt(claims jwt.MapClaims) (userId string, err error) {
	userId, ok := claims["user_id"].(string)
	if !ok {
		return userId, utils.ErrUnauthorized
	}

	return userId, nil
}
