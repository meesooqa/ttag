package db

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"

	"github.com/meesooqa/ttag/app/config"
	"github.com/meesooqa/ttag/app/model"
)

type DB interface {
	UpsertMany(messagesChan <-chan model.Message)
}

type MongoDB struct {
	log    *zap.Logger
	conf   *config.MongoConfig
	client *mongo.Client
}

func NewMongoDB(log *zap.Logger, conf *config.MongoConfig) *MongoDB {
	return &MongoDB{
		log:  log,
		conf: conf,
	}
}

func (db *MongoDB) UpsertMany(messagesChan <-chan model.Message) {
	batchSize := 10
	flushPeriod := 2 // Seconds

	collection := db.GetCollection(db.conf.CollectionMessages)
	saver := NewSaver(db.log, collection, batchSize, time.Duration(flushPeriod)*time.Second, 50)
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
	time.Sleep(time.Duration(flushPeriod+1) * time.Second) // wait flushPeriod
	saver.Close()
	db.log.Debug("all data has been successfully saved to MongoDB")
}

func (db *MongoDB) Init() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	connectedClient, err := mongo.Connect(ctx, options.Client().ApplyURI(db.conf.URI))
	if err != nil {
		return err
	}

	err = connectedClient.Ping(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to ping MongoDB: %v", err)
	}

	db.client = connectedClient

	if err := db.createUniqueUuidIndex(context.TODO()); err != nil {
		db.log.Fatal("creating index", zap.Error(err))
	}

	return nil
}

func (db *MongoDB) Close() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if db.client != nil {
		if err := db.client.Disconnect(ctx); err != nil {
			db.log.Error("failed to disconnect MongoDB", zap.Error(err))
		}
	}
}

func (db *MongoDB) GetDatabase() *mongo.Database {
	return db.client.Database(db.conf.Database)
}

func (db *MongoDB) GetCollection(collectionName string) *mongo.Collection {
	return db.GetDatabase().Collection(collectionName)
}

func (db *MongoDB) createUniqueUuidIndex(ctx context.Context) error {
	collection := db.GetCollection(db.conf.CollectionMessages)
	indexModel := mongo.IndexModel{
		Keys:    bson.D{{Key: "uuid", Value: 1}},
		Options: options.Index().SetUnique(true),
	}
	_, err := collection.Indexes().CreateOne(ctx, indexModel)
	return err
}
