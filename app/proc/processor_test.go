package proc

import (
	"errors"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest"
	"go.uber.org/zap/zaptest/observer"

	"github.com/meesooqa/ttag/app/proc/mocks"
)

func TestProcessor_ProcessFile_Success(t *testing.T) {
	logger := zaptest.NewLogger(t)
	fService := &mocks.ServiceMock{}
	fRepo := &mocks.RepositoryMock{}
	processor := NewProcessor(logger, fService, fRepo)

	filesChan := make(chan string, 1)
	filesChan <- "file1.txt"
	close(filesChan)

	var wg sync.WaitGroup
	wg.Add(1)
	processor.ProcessFile(filesChan, &wg)
	wg.Wait()

	assert.Equal(t, 1, fService.CallCount, "Ожидается, что ParseArchivedFile будет вызван один раз")
	assert.Equal(t, 1, len(fRepo.UpsertCalls), "Ожидается, что одно сообщение будет обработано")
	assert.Equal(t, "file1.txt", fRepo.UpsertCalls[0].MessageID)
}

func TestProcessor_ProcessFile_Error(t *testing.T) {
	core, observedLogs := observer.New(zap.ErrorLevel)
	logger := zap.New(core)
	parseErr := errors.New("db upsert error")
	fService := &mocks.ServiceMock{
		Err: parseErr,
	}
	fRepo := &mocks.RepositoryMock{}
	processor := NewProcessor(logger, fService, fRepo)

	filesChan := make(chan string, 1)
	filesChan <- "file2.txt"
	close(filesChan)

	var wg sync.WaitGroup
	wg.Add(1)
	processor.ProcessFile(filesChan, &wg)
	wg.Wait()

	assert.Equal(t, 1, fService.CallCount, "Ожидается, что ParseArchivedFile будет вызван один раз")
	assert.Equal(t, 1, len(fRepo.UpsertCalls), "Ожидается, что одно сообщение будет обработано")

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
