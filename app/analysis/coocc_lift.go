package analysis

import (
	"context"
	"log/slog"

	"github.com/meesooqa/ttag/app/repositories"
)

type CooccLiftDataProvider struct {
	log  *slog.Logger
	repo repositories.Repository
}

type CooccLiftData struct{}

func NewCooccLiftDataProvider(log *slog.Logger, repo repositories.Repository) *CooccLiftDataProvider {
	return &CooccLiftDataProvider{
		log:  log,
		repo: repo,
	}
}

func (p *CooccLiftDataProvider) GetData(ctx context.Context, group string) AnalyzedData {
	//filter := bson.M{}
	//if group != "" {
	//	filter = bson.M{"group": group}
	//}
	//messages, err := p.repo.Find(ctx, filter)
	//if err != nil {
	//	p.log.Error("all messages finding", "err", err)
	//}

	return &CooccLiftData{}
}
