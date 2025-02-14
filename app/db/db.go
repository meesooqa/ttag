package db

import (
	"go.uber.org/zap"
	"time"

	"github.com/meesooqa/ttag/app/tg"
)

type DB interface {
	Upsert(message tg.ArchivedMessage) error
}

type MongoDB struct {
	log *zap.Logger
}

func NewMongoDB(log *zap.Logger) *MongoDB {
	return &MongoDB{
		log: log,
	}
}

func (db *MongoDB) Upsert(message tg.ArchivedMessage) error {
	db.log.Debug("Upsert")
	time.Sleep(time.Second)
	return nil
}
