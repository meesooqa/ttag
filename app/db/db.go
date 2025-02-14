package db

import (
	"context"
	"github.com/meesooqa/ttag/app/tg"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
	"time"
)

type DB interface {
	UpsertMany(messagesChan <-chan tg.ArchivedMessage)
}

type MongoDB struct {
	log *zap.Logger
}

func NewMongoDB(log *zap.Logger) *MongoDB {
	return &MongoDB{
		log: log,
	}
}

func (db *MongoDB) UpsertMany(messagesChan <-chan tg.ArchivedMessage) {
	db.log.Debug("UpsertMany")

	// Подключаемся к MongoDB
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		db.log.Fatal("Ошибка подключения к MongoDB:", zap.Error(err))
	}
	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			db.log.Fatal("Ошибка отключения от MongoDB:", zap.Error(err))
		}
	}()

	collection := client.Database("test").Collection("logs")
	saver := NewSaver(collection, 10, 5*time.Second, 50)
	go func() {
		for msg := range messagesChan {
			doc := bson.M{
				"message_id": msg.MessageID,
				"datetime":   msg.Datetime,
				"tags":       msg.Tags,
			}
			saver.Save(doc)
		}
	}()
	saver.Close()
	db.log.Debug("Все данные успешно сохранены в MongoDB")
}
