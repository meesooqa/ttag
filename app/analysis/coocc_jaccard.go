package analysis

import (
	"context"
	"log/slog"

	"github.com/meesooqa/ttag/app/repositories"
)

type CooccJaccardDataProvider struct {
	log  *slog.Logger
	repo repositories.Repository
}

type CooccJaccardData struct{}

func NewCooccJaccardDataProvider(log *slog.Logger, repo repositories.Repository) *CooccJaccardDataProvider {
	return &CooccJaccardDataProvider{
		log:  log,
		repo: repo,
	}
}

func (p *CooccJaccardDataProvider) GetData(ctx context.Context, group string) AnalyzedData {
	//filter := bson.M{}
	//if group != "" {
	//	filter = bson.M{"group": group}
	//}
	//messages, err := p.repo.Find(ctx, filter)
	//if err != nil {
	//	p.log.Error("all messages finding", "err", err)
	//}

	return &CooccJaccardData{}
}
