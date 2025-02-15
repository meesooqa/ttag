package proc

import (
	"sync"

	"go.uber.org/zap"

	"github.com/meesooqa/ttag/app/db"
	"github.com/meesooqa/ttag/app/tg"
)

type Processor struct {
	log     *zap.Logger
	service tg.Service
	db      db.DB
}

func NewProcessor(log *zap.Logger, service tg.Service, db db.DB) *Processor {
	return &Processor{
		log:     log,
		service: service,
		db:      db,
	}
}

func (p *Processor) ProcessFile(filesChan <-chan string, wg *sync.WaitGroup) {
	defer wg.Done()

	messagesChan := make(chan tg.ArchivedMessage, 10)

	var wgm sync.WaitGroup
	wgm.Add(1)
	go func() {
		defer wgm.Done()
		p.db.UpsertMany(messagesChan)
	}()

	for filename := range filesChan {
		if err := p.service.ParseArchivedFile(filename, messagesChan); err != nil {
			p.log.Error("Error processing file", zap.String("filename", filename), zap.Error(err))
		}
	}
	close(messagesChan)

	wgm.Wait()
}
