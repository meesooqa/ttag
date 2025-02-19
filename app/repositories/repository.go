package repositories

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/meesooqa/ttag/app/model"
)

type Repository interface {
	UpsertMany(messagesChan <-chan model.Message)
	GetUniqueValues(ctx context.Context, fieldName string) ([]string, error)
	Find(ctx context.Context, filter bson.M, opts ...*options.FindOptions) ([]*model.Message, error)
}
