package mongo

import (
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"template/internal/config"
	dto "template/internal/dto/user"
	"template/internal/entity/user"
	"template/internal/utils"
)

type UserRepo struct {
	cli *mongo.Client
	cfg *config.MongoRepo
}

func NewUserRepo(cfg *config.MongoRepo, cli *mongo.Client) *UserRepo {
	return &UserRepo{cli: cli, cfg: cfg}
}

func (repo *UserRepo) CreateUser(ctx context.Context, user user.User) (string, error) {
	collection := repo.cli.Database(repo.cfg.Database).Collection("users")
	res, err := collection.InsertOne(ctx, user)

	if err != nil {
		return "", err
	}

	objectId, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		return fmt.Sprint(res.InsertedID), nil
	}
	return objectId.Hex(), nil
}

func (repo *UserRepo) GetUserById(ctx context.Context, userId string) (dto.User, error) {
	u := user.User{}
	collection := repo.cli.Database(repo.cfg.Database).Collection("users")

	id, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return dto.User{}, utils.ErrWrongUserId
	}
	filter := bson.M{"_id": id}

	cursor := collection.FindOne(ctx, filter)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return dto.User{}, utils.ErrUserNotFound
	} else if err != nil {
		return dto.User{}, err
	}
	if err := cursor.Decode(&u); err != nil {
		return dto.User{}, err
	}

	return u.Json(), nil
}

func (repo *UserRepo) UpdateUser(ctx context.Context, userID string, upd dto.UpdateUserRequest) error {
	collection := repo.cli.Database(repo.cfg.Database).Collection("users")
	objId, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return err
	}
	filter := bson.M{"_id": objId}
	update := bson.M{"$set": bson.M{}}
	if upd.Name != nil {
		update["$set"] = bson.M{"name": *upd.Name}
	}
	if upd.Email != nil {
		update["$set"] = bson.M{"email": *upd.Email}
	}
	res, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return utils.NewInternalError("unable to update user " + err.Error())
	}
	if res.MatchedCount == 0 {
		return utils.ErrUserNotFound
	}
	return nil
}
func (repo *UserRepo) GetUserByEmail(ctx context.Context, email string) (dto.User, error) {
	u := user.User{}
	collection := repo.cli.Database(repo.cfg.Database).Collection("users")

	filter := bson.M{
		"email": email}
	err := collection.FindOne(ctx, filter).Decode(&u)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return dto.User{}, utils.ErrUserNotFound
	} else if err != nil {
		return dto.User{}, err
	}

	return u.Json(), nil
}

func (repo *UserRepo) AdminExists(ctx context.Context) error {
	u := user.User{}
	collection := repo.cli.Database(repo.cfg.Database).Collection("users")

	filter := bson.M{"role": "admin"}

	err := collection.FindOne(ctx, filter).Decode(&u)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return utils.ErrUserNotFound
	} else if err != nil {
		return err
	}

	return nil
}
