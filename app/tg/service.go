package tg

import (
	"go.uber.org/zap"

	"github.com/meesooqa/ttag/app/model"
)

type Service interface {
	ParseArchivedFile(filename string, messagesChan chan<- model.ArchivedMessage) error
}

type TgService struct {
	log    *zap.Logger
	parser Parser
}

func NewService(log *zap.Logger) *TgService {
	return &TgService{
		log:    log,
		parser: NewTgArchivedHTMLParser(log),
	}
}

func (s *TgService) ParseArchivedFile(filename string, messagesChan chan<- model.ArchivedMessage) error {
	return s.parser.ParseFile(filename, messagesChan)
}
