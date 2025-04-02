package user

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	dto "template/internal/dto/user"
)

type Role string

const (
	Admin Role = "admin"
	Usr   Role = "user"
)

type ContextUserKey struct{}
type User struct {
	ID       *primitive.ObjectID `bson:"_id,omitempty"`
	Name     string              `bson:"name"`
	Email    string              `bson:"email"`
	Password string              `bson:"password"`
	Role     Role                `bson:"role"`
}

func (u *User) Json() dto.User {
	return dto.User{
		ID:    u.ID.Hex(),
		Name:  u.Name,
		Email: u.Email,
		Role:  string(u.Role),
	}
}
