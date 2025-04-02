package mongo

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"template/internal/config"
	"time"
)

func NewMongoClient(config *config.MongoRepo) (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	opts := options.Client().ApplyURI(config.Address)
	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		return nil, err
	}

	if err = client.Ping(ctx, nil); err != nil {
		return nil, err
	}
	collConfig := getCollectionsConfig()
	if err = ensureCollections(context.Background(), client.Database(config.Database), collConfig); err != nil {
		return nil, err
	}
	return client, nil
}

type IndexConfig struct {
	Keys   bson.D
	Unique bool
}

type CollectionConfig struct {
	Name    string
	Indexes []IndexConfig
}

func getCollectionsConfig() []CollectionConfig {
	return []CollectionConfig{

		{Name: "users",
			Indexes: []IndexConfig{
				{
					Keys: bson.D{
						{Key: "email", Value: 1},
					}},
			}},
	}
}

func ensureCollections(ctx context.Context, db *mongo.Database, collections []CollectionConfig) error {
	existing, err := db.ListCollectionNames(ctx, bson.D{})
	if err != nil {
		return fmt.Errorf("listing collections: %w", err)
	}

	existingMap := make(map[string]bool)
	for _, name := range existing {
		existingMap[name] = true
	}

	for _, col := range collections {
		if !existingMap[col.Name] {
			if err := db.CreateCollection(ctx, col.Name); err != nil {
				return fmt.Errorf("creating collection %s: %w", col.Name, err)
			}
		}

		if len(col.Indexes) > 0 {
			for _, idx := range col.Indexes {
				indexOpts := options.Index()
				indexOpts.SetUnique(idx.Unique)

				_, err := db.Collection(col.Name).Indexes().CreateOne(ctx, mongo.IndexModel{
					Keys:    idx.Keys,
					Options: indexOpts,
				})
				if err != nil {
					return fmt.Errorf("creating index for %s: %w", col.Name, err)
				}
			}
		}
	}

	return nil
}
