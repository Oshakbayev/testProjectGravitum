package server

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type db interface {
	Disconnect(ctx context.Context) error
	Ping(ctx context.Context, rp *readpref.ReadPref) error
}
