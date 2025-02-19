package mocks

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/meesooqa/ttag/app/model"
)

type RepositoryMock struct {
	UpsertCalls []model.Message
	Err         error
}

func (f *RepositoryMock) UpsertMany(messagesChan <-chan model.Message) {
	for m := range messagesChan {
		f.UpsertCalls = append(f.UpsertCalls, m)
	}
}

func (f *RepositoryMock) GetUniqueValues(ctx context.Context, fieldName string) ([]string, error) {
	return nil, nil
}

func (f *RepositoryMock) Find(ctx context.Context, filter bson.M, opts ...*options.FindOptions) ([]*model.Message, error) {
	return nil, nil
}
