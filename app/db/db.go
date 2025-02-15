package db

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"

	"github.com/meesooqa/ttag/app/tg"
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
	// TODO Refactor UpsertMany
	db.log.Debug("UpsertMany")

	ctx := context.TODO()
	// Подключаемся к MongoDB
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		db.log.Fatal("Ошибка подключения к MongoDB:", zap.Error(err))
	}
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			db.log.Fatal("Ошибка отключения от MongoDB:", zap.Error(err))
		}
	}()

	collection := client.Database("db_tags").Collection("tags")
	if err := db.createUniqueUuidIndex(ctx, collection); err != nil {
		db.log.Fatal("Ошибка создания индекса:", zap.Error(err))
	}

	saver := NewSaver(db.log, collection, 10, 2*time.Second, 50)
	go func() {
		for msg := range messagesChan {
			doc := bson.M{
				"message_id": msg.MessageID,
				"datetime":   msg.Datetime,
				"group":      msg.Group,
				"uuid":       msg.UUID,
				"tags":       msg.Tags,
			}
			if err := saver.Save(doc); err != nil {
				db.log.Error("Saver error", zap.Error(err))
			}
		}
	}()
	time.Sleep(3 * time.Second) // wait flushPeriod
	saver.Close()
	db.log.Debug("Все данные успешно сохранены в MongoDB")
}

func (db *MongoDB) createUniqueUuidIndex(ctx context.Context, collection *mongo.Collection) error {
	indexModel := mongo.IndexModel{
		Keys:    bson.D{{Key: "uuid", Value: 1}},
		Options: options.Index().SetUnique(true),
	}
	_, err := collection.Indexes().CreateOne(ctx, indexModel)
	return err
}
