package repositories

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/meesooqa/ttag/app/db"
	"github.com/meesooqa/ttag/app/model"
)

type MessageRepository struct {
	log        *slog.Logger
	collection *mongo.Collection
}

func NewMessageRepository(log *slog.Logger, db *db.MongoDB) *MessageRepository {
	return &MessageRepository{
		log:        log,
		collection: db.GetCollectionMessages(),
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
				r.log.Error("Saver error", "err", err)
			}
		}
	}()
	time.Sleep(time.Duration(flushPeriod+1) * time.Second) // wait flushPeriod
	saver.Close()
	r.log.Debug("all data has been successfully saved to MongoDB")
}

func (r *MessageRepository) Find(ctx context.Context, filter bson.M, opts ...*options.FindOptions) ([]*model.Message, error) {
	cursor, err := r.collection.Find(ctx, filter, opts...)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	var items []*model.Message
	for cursor.Next(ctx) {
		var item model.Message
		if err := cursor.Decode(&item); err != nil {
			return nil, err
		}
		items = append(items, &item)
	}
	if err := cursor.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

func (r *MessageRepository) GetUniqueValues(ctx context.Context, fieldName string) ([]string, error) {
	values, err := r.collection.Distinct(ctx, fieldName, bson.D{})
	if err != nil {
		return nil, fmt.Errorf("distinct failed: %w", err)
	}
	return convertToStrings(values), nil
}

func convertToStrings(values []interface{}) []string {
	result := make([]string, 0, len(values))
	for _, v := range values {
		if str, ok := v.(string); ok {
			result = append(result, str)
		}
	}
	return result
}
