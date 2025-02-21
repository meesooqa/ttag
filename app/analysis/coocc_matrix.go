package analysis

import (
	"context"
	"log/slog"

	"github.com/meesooqa/ttag/app/repositories"
)

type CooccMatrixDataProvider struct {
	log  *slog.Logger
	repo repositories.Repository
}

type CooccMatrixData struct{}

func NewCooccMatrixDataProvider(log *slog.Logger, repo repositories.Repository) *CooccMatrixDataProvider {
	return &CooccMatrixDataProvider{
		log:  log,
		repo: repo,
	}
}

func (p *CooccMatrixDataProvider) GetData(ctx context.Context, group string) AnalyzedData {
	//filter := bson.M{}
	//if group != "" {
	//	filter = bson.M{"group": group}
	//}
	//messages, err := p.repo.Find(ctx, filter)
	//if err != nil {
	//	p.log.Error("all messages finding", "err", err)
	//}
	//messages := []model.Message{
	//	{Tags: []string{"tag1", "tag2"}},
	//	{Tags: []string{"tag1", "tag3"}},
	//	{Tags: []string{"tag1", "tag4"}},
	//	{Tags: []string{"tag2", "tag4"}},
	//}

	return nil
}
