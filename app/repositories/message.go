package repositories

import (
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"

	"github.com/meesooqa/ttag/app/model"
)

type MessageRepository struct {
	log        *zap.Logger
	collection *mongo.Collection
}

func NewMessageRepository(log *zap.Logger, collection *mongo.Collection) *MessageRepository {
	return &MessageRepository{
		log:        log,
		collection: collection,
	}
}

func (r *MessageRepository) UpsertMany(messagesChan <-chan model.Message) {
	batchSize := 10
	flushPeriod := 2 // Seconds

	saver := NewSaver(r.log, r.collection, batchSize, time.Duration(flushPeriod)*time.Second, 50)
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
				r.log.Error("Saver error", zap.Error(err))
			}
		}
	}()
	time.Sleep(time.Duration(flushPeriod+1) * time.Second) // wait flushPeriod
	saver.Close()
	r.log.Debug("all data has been successfully saved to MongoDB")
}
