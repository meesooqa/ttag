package repositories

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/meesooqa/ttag/app/model"
)

type Repository interface {
	Find(ctx context.Context, filter bson.M, opts ...*options.FindOptions) ([]*model.Message, error)
	UpsertMany(messagesChan <-chan model.Message)
	GetUniqueValues(ctx context.Context, fieldName string) ([]string, error)
	GetTags(ctx context.Context, query string) ([]string, error)
}
