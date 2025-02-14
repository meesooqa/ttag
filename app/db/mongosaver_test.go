package db

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// fakeCollection реализует Inserter для тестирования.
type fakeCollection struct {
	mu         sync.Mutex
	inserted   [][]bson.M // каждая вставка сохраняется как отдельный батч
	failInsert bool       // если true, InsertMany возвращает ошибку
}

func (f *fakeCollection) InsertMany(ctx context.Context, docs []interface{}, opts ...*options.InsertManyOptions) (*mongo.InsertManyResult, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	var batch []bson.M
	for _, doc := range docs {
		if m, ok := doc.(bson.M); ok {
			batch = append(batch, m)
		}
	}
	if f.failInsert {
		return nil, errors.New("insert error")
	}
	f.inserted = append(f.inserted, batch)
	return &mongo.InsertManyResult{}, nil
}

// TestSaver_BatchSaving проверяет пакетную вставку по достижению batchSize.
func TestSaver_BatchSaving(t *testing.T) {
	fc := &fakeCollection{}
	saver := NewSaver(fc, 3, 100*time.Millisecond, 10)

	// Отправляем 5 документов.
	saver.Save(bson.M{"a": 1})
	saver.Save(bson.M{"a": 2})
	saver.Save(bson.M{"a": 3}) // первый батч (3 дока)
	saver.Save(bson.M{"a": 4})
	saver.Save(bson.M{"a": 5})

	// Ждём, чтобы сработал тикер и завершились вставки.
	time.Sleep(200 * time.Millisecond)
	saver.Close()

	fc.mu.Lock()
	defer fc.mu.Unlock()

	// Ожидаем два батча: первый из 3 документов, второй из 2.
	if len(fc.inserted) != 2 {
		t.Fatalf("ожидали 2 батча, получили %d", len(fc.inserted))
	}
	if len(fc.inserted[0]) != 3 {
		t.Errorf("ожидали 3 документа в первом батче, получили %d", len(fc.inserted[0]))
	}
	if len(fc.inserted[1]) != 2 {
		t.Errorf("ожидали 2 документа во втором батче, получили %d", len(fc.inserted[1]))
	}
}

// TestSaver_FlushOnTicker проверяет, что оставшиеся документы сохраняются по таймеру.
func TestSaver_FlushOnTicker(t *testing.T) {
	fc := &fakeCollection{}
	saver := NewSaver(fc, 10, 100*time.Millisecond, 10)

	// Отправляем меньше документов, чем batchSize, чтобы сработал таймер.
	saver.Save(bson.M{"a": 1})
	saver.Save(bson.M{"a": 2})

	// Ждём срабатывания flushPeriod.
	time.Sleep(200 * time.Millisecond)
	saver.Close()

	fc.mu.Lock()
	defer fc.mu.Unlock()

	if len(fc.inserted) != 1 {
		t.Fatalf("ожидали 1 батч, получили %d", len(fc.inserted))
	}
	if len(fc.inserted[0]) != 2 {
		t.Errorf("ожидали 2 документа в батче, получили %d", len(fc.inserted[0]))
	}
}

// TestSaver_InsertError проверяет, что ошибка вставки не вызывает панику.
func TestSaver_InsertError(t *testing.T) {
	fc := &fakeCollection{failInsert: true}
	saver := NewSaver(fc, 2, 50*time.Millisecond, 10)

	saver.Save(bson.M{"a": 1})
	saver.Save(bson.M{"a": 2})

	time.Sleep(100 * time.Millisecond)
	// Если ошибка вставки вызовет панику, тест упадёт.
	saver.Close()
}
