package db

import (
	"context"
	"log"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Inserter представляет сущность, поддерживающую пакетную вставку документов.
type Inserter interface {
	InsertMany(ctx context.Context, docs []interface{}, opts ...*options.InsertManyOptions) (*mongo.InsertManyResult, error)
}

// Saver отвечает за сбор и пакетную отправку данных в MongoDB.
type Saver struct {
	collection  Inserter
	dataChan    chan bson.M
	batchSize   int
	flushPeriod time.Duration
	wg          sync.WaitGroup
}

// NewSaver создаёт новый Saver с указанными параметрами.
func NewSaver(collection Inserter, batchSize int, flushPeriod time.Duration, bufferSize int) *Saver {
	s := &Saver{
		collection:  collection,
		dataChan:    make(chan bson.M, bufferSize),
		batchSize:   batchSize,
		flushPeriod: flushPeriod,
	}
	s.wg.Add(1)
	go s.run()
	return s
}

// run запускает обработку канала и периодическое сохранение в MongoDB.
func (s *Saver) run() {
	defer s.wg.Done()
	ticker := time.NewTicker(s.flushPeriod)
	defer ticker.Stop()

	var batch []bson.M

	for {
		select {
		case doc, ok := <-s.dataChan:
			if !ok {
				if len(batch) > 0 {
					s.saveBatch(batch)
				}
				return
			}
			batch = append(batch, doc)

			if len(batch) >= s.batchSize {
				s.saveBatch(batch)
				batch = batch[:0]
			}
		case <-ticker.C:
			if len(batch) > 0 {
				s.saveBatch(batch)
				batch = batch[:0]
			}
		}
	}
}

// saveBatch сохраняет батч документов в MongoDB.
func (s *Saver) saveBatch(batch []bson.M) {
	// Приводим []bson.M к []interface{}
	docs := make([]interface{}, len(batch))
	for i, doc := range batch {
		docs[i] = doc
	}
	_, err := s.collection.InsertMany(context.TODO(), docs)
	if err != nil {
		log.Println("Mongo InsertMany error:", err)
	}
}

// Save добавляет документ в очередь сохранения.
func (s *Saver) Save(doc bson.M) {
	s.dataChan <- doc
}

// Close завершает работу и сохраняет остатки.
func (s *Saver) Close() {
	close(s.dataChan)
	s.wg.Wait()
}
