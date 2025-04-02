package auth

import (
	"errors"
	"template/internal/dto/user"
	"template/internal/utils"
)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (r *LoginRequest) Validate() error {
	if !utils.IsValidEmail(r.Email) {
		return errors.New("invalid email")
	}
	return nil
}

type LoginResponse struct {
	AccessToken  string     `json:"jwt"`
	RefreshToken string     `json:"refresh_token"`
	User         *user.User `json:"user,omitempty"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}
