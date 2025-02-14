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

func NewProcessor(log *zap.Logger) *Processor {
	return &Processor{
		log:     log,
		service: tg.NewService(log),
		db:      db.NewMongoDB(log),
	}
}

func (p *Processor) ProcessFile(filesChan <-chan string, wg *sync.WaitGroup) {
	defer wg.Done()

	messagesChan := make(chan tg.ArchivedMessage, 3)

	var wgm sync.WaitGroup
	wgm.Add(1)
	go func() {
		defer wgm.Done()
		p.processMessages(messagesChan)
	}()

	for filename := range filesChan {
		p.service.ParseArchivedFile(filename, messagesChan)
	}
	close(messagesChan)

	wgm.Wait()
}

func (p *Processor) processMessages(messagesChan <-chan tg.ArchivedMessage) {
	var wg sync.WaitGroup

	for message := range messagesChan {
		wg.Add(1)
		go func(msg tg.ArchivedMessage) {
			defer wg.Done()
			if err := p.processMessage(msg); err != nil {
				p.log.Error("Error processing message", zap.Error(err))
			}
		}(message)
	}

	wg.Wait()
}

func (p *Processor) processMessage(message tg.ArchivedMessage) error {
	err := p.db.Upsert(message)
	if err != nil {
		return err
	}
	return nil
}
