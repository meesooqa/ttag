package analysis

import (
	"context"
	"log/slog"

	"github.com/meesooqa/ttag/app/repositories"
)

type CooccPmiDataProvider struct {
	log  *slog.Logger
	repo repositories.Repository
}

type CooccPmiData struct{}

func NewCooccPmiDataProvider(log *slog.Logger, repo repositories.Repository) *CooccPmiDataProvider {
	return &CooccPmiDataProvider{
		log:  log,
		repo: repo,
	}
}

func (p *CooccPmiDataProvider) GetData(ctx context.Context, group string) AnalyzedData {
	//filter := bson.M{}
	//if group != "" {
	//	filter = bson.M{"group": group}
	//}
	//messages, err := p.repo.Find(ctx, filter)
	//if err != nil {
	//	p.log.Error("all messages finding", "err", err)
	//}

	return &CooccPmiData{}
}
