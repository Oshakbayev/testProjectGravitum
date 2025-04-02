package user

import (
	"template/internal/utils"
)

type User struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Role  string `json:"role"`
}

type CreateUserRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (u *CreateUserRequest) Validate() error {
	if !utils.IsValidEmail(u.Email) {
		return utils.ErrInvalidEmail
	}

	if err := utils.IsValidPassword(u.Password); err != nil {
		return err
	}

	return nil
}

type UpdateUserRequest struct {
	Name  *string `json:"name"`
	Email *string `json:"email"`
}

func (u *UpdateUserRequest) Validate() error {
	if u.Email != nil && !utils.IsValidEmail(*u.Email) {
		return utils.ErrInvalidEmail
	}

	return nil
}
