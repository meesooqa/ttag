package repositories

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

func (r *MessageRepository) Create(ctx context.Context, item *model.Message) error {
	result, err := r.collection.InsertOne(ctx, item)
	if err != nil {
		return err
	}
	if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
		item.ID = oid
	}
	return nil
}

func (r *MessageRepository) GetByID(ctx context.Context, id primitive.ObjectID) (*model.Message, error) {
	var item model.Message
	filter := bson.M{"_id": id}
	err := r.collection.FindOne(ctx, filter).Decode(&item)
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *MessageRepository) Update(ctx context.Context, id primitive.ObjectID, update bson.M) error {
	filter := bson.M{"_id": id}
	_, err := r.collection.UpdateOne(ctx, filter, update)
	return err
}

func (r *MessageRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	filter := bson.M{"_id": id}
	_, err := r.collection.DeleteOne(ctx, filter)
	return err
}
