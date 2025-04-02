package mongo

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
	"template/internal/config"
	"template/internal/dto/auth"
	dto "template/internal/dto/user"
	entity "template/internal/entity/user"
	"template/internal/utils"
)

type AuthRepo struct {
	cli *mongo.Client
	cfg *config.MongoRepo
}

func NewAuthRepo(cfg *config.MongoRepo, cli *mongo.Client) *AuthRepo {
	return &AuthRepo{cli: cli, cfg: cfg}
}

func (repo *AuthRepo) Login(ctx context.Context, req auth.LoginRequest) (dto.User, error) {
	u := entity.User{}
	collection := repo.cli.Database(repo.cfg.Database).Collection("users")

	filter := bson.M{
		"email": req.Email,
	}

	err := collection.FindOne(ctx, filter).Decode(&u)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return dto.User{}, utils.ErrUserNotFound
		}
		return dto.User{}, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(req.Password))
	if err != nil {
		return dto.User{}, utils.ErrUserNotFound
	}

	return u.Json(), nil
}
