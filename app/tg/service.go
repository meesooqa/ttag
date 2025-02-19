package tg

import (
	"go.uber.org/zap"

	"github.com/meesooqa/ttag/app/config"
	"github.com/meesooqa/ttag/app/model"
)

type Service interface {
	ParseArchivedFile(filename string, messagesChan chan<- model.Message) error
}

type TgService struct {
	log    *zap.Logger
	parser Parser
}

func NewService(log *zap.Logger, conf *config.SystemConfig) *TgService {
	return &TgService{
		log:    log,
		parser: NewTgArchivedHTMLParser(log, conf.DataPath),
	}
}

func (s *TgService) ParseArchivedFile(filename string, messagesChan chan<- model.Message) error {
	return s.parser.ParseFile(filename, messagesChan)
}
