package repositories

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/meesooqa/ttag/app/model"
)

type Repository interface {
	UpsertMany(messagesChan <-chan model.Message)
	Create(ctx context.Context, item *model.Message) error
	GetByID(ctx context.Context, id primitive.ObjectID) (*model.Message, error)
	Find(ctx context.Context, filter bson.M, opts ...*options.FindOptions) ([]*model.Message, error)
	Update(ctx context.Context, id primitive.ObjectID, update bson.M) error
	Delete(ctx context.Context, id primitive.ObjectID) error
}
