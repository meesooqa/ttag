package proc

import (
	"errors"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest"
	"go.uber.org/zap/zaptest/observer"

	"github.com/meesooqa/ttag/app/tg"
)

type fakeService struct {
	callCount int
}

func (fs *fakeService) ParseArchivedFile(filename string, messagesChan chan<- tg.ArchivedMessage) {
	fs.callCount++
	messagesChan <- tg.ArchivedMessage{
		MessageID: filename,
	}
}

type fakeDB struct {
	upsertCalls []tg.ArchivedMessage
	err         error
}

func (f *fakeDB) Upsert(message tg.ArchivedMessage) error {
	f.upsertCalls = append(f.upsertCalls, message)
	return f.err
}

func TestProcessor_ProcessFile_Success(t *testing.T) {
	logger := zaptest.NewLogger(t)
	processor := NewProcessor(logger)

	fService := &fakeService{}
	processor.service = fService

	fDB := &fakeDB{}
	processor.db = fDB

	filesChan := make(chan string, 1)
	filesChan <- "file1.txt"
	close(filesChan)

	var wg sync.WaitGroup
	wg.Add(1)
	processor.ProcessFile(filesChan, &wg)
	wg.Wait()

	assert.Equal(t, 1, fService.callCount, "Ожидается, что ParseArchivedFile будет вызван один раз")
	assert.Equal(t, 1, len(fDB.upsertCalls), "Ожидается, что одно сообщение будет обработано")
	assert.Equal(t, "file1.txt", fDB.upsertCalls[0].MessageID)
}

func TestProcessor_ProcessFile_DBError(t *testing.T) {
	core, observedLogs := observer.New(zap.ErrorLevel)
	logger := zap.New(core)
	processor := NewProcessor(logger)

	fService := &fakeService{}
	processor.service = fService

	dbErr := errors.New("db upsert error")
	fDB := &fakeDB{err: dbErr}
	processor.db = fDB

	filesChan := make(chan string, 1)
	filesChan <- "file2.txt"
	close(filesChan)

	var wg sync.WaitGroup
	wg.Add(1)
	processor.ProcessFile(filesChan, &wg)
	wg.Wait()

	assert.Equal(t, 1, fService.callCount, "Ожидается, что ParseArchivedFile будет вызван один раз")
	assert.Equal(t, 1, len(fDB.upsertCalls), "Ожидается, что одно сообщение будет обработано")

	var found bool
	for _, entry := range observedLogs.All() {
		if entry.Message == "Error processing message" && entry.Context != nil {
			if errField := entry.Context[0]; errField.Key == "error" /*&& errField.String == dbErr.Error()*/ {
				found = true
				break
			}
		}
	}
	assert.True(t, found, "Ожидается, что ошибка обработки сообщения будет залогирована")
}
