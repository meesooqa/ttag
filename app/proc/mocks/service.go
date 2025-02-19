package mocks

import "github.com/meesooqa/ttag/app/model"

type ServiceMock struct {
	CallCount int
	Err       error
}

func (fs *ServiceMock) ParseArchivedFile(filename string, messagesChan chan<- model.Message) error {
	fs.CallCount++
	messagesChan <- model.Message{
		MessageID: filename,
	}
	return fs.Err
}
