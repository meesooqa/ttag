package repositories

import "github.com/meesooqa/ttag/app/model"

type Repository interface {
	UpsertMany(messagesChan <-chan model.Message)
}
