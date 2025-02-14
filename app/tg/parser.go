package tg

import "go.uber.org/zap"

type Parser interface {
	ParseFile(filename string, messagesChan chan<- ArchivedMessage)
}

type TgArchivedHTMLParser struct {
	log *zap.Logger
}

func NewTgArchivedHTMLParser(log *zap.Logger) *TgArchivedHTMLParser {
	return &TgArchivedHTMLParser{
		log: log,
	}
}

func (p *TgArchivedHTMLParser) ParseFile(filename string, messagesChan chan<- ArchivedMessage) {
	p.log.Debug("TgArchivedHTMLParser got file:", zap.String("filename", filename))
}
