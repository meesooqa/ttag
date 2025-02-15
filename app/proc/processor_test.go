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
	err       error
}

func (fs *fakeService) ParseArchivedFile(filename string, messagesChan chan<- tg.ArchivedMessage) error {
	fs.callCount++
	messagesChan <- tg.ArchivedMessage{
		MessageID: filename,
	}
	return fs.err
}

type fakeDB struct {
	upsertCalls []tg.ArchivedMessage
	err         error
}

func (f *fakeDB) UpsertMany(messagesChan <-chan tg.ArchivedMessage) {
	for m := range messagesChan {
		f.upsertCalls = append(f.upsertCalls, m)
	}
}

func TestProcessor_ProcessFile_Success(t *testing.T) {
	logger := zaptest.NewLogger(t)
	fService := &fakeService{}
	fDB := &fakeDB{}
	processor := NewProcessor(logger, fService, fDB)

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

func TestProcessor_ProcessFile_Error(t *testing.T) {
	core, observedLogs := observer.New(zap.ErrorLevel)
	logger := zap.New(core)
	parseErr := errors.New("db upsert error")
	fService := &fakeService{
		err: parseErr,
	}
	fDB := &fakeDB{}
	processor := NewProcessor(logger, fService, fDB)

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
		if entry.Message == "Error processing file" && entry.Context != nil {
			if errField := entry.Context[1]; errField.Key == "error" &&
				errField.Interface.(error).Error() == parseErr.Error() {
				if filenameField := entry.Context[0]; filenameField.String == "file2.txt" {
					found = true
					break
				}
			}
		}
	}
	assert.True(t, found, "Ожидается, что ошибка обработки сообщения будет залогирована")
}
