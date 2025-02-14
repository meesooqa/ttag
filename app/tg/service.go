package tg

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
)

type ArchivedMessage struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	MessageID string             `bson:"message_id" json:"MessageID"`
	Datetime  time.Time          `bson:"datetime" json:"Datetime"`
	Tags      []string           `bson:"tags" json:"Tags"`
}

type Service interface {
	ParseArchivedFile(filename string, messagesChan chan<- ArchivedMessage) error
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

func (s *TgService) ParseArchivedFile(filename string, messagesChan chan<- ArchivedMessage) error {
	return s.parser.ParseFile(filename, messagesChan)
}
