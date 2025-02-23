package analysis

import (
	"context"
	"log/slog"

	"github.com/meesooqa/ttag/app/repositories"
)

type CooccConfidenceDataProvider struct {
	log  *slog.Logger
	repo repositories.Repository
}

type CooccConfidenceData struct {
}

func NewCooccConfidenceDataProvider(log *slog.Logger, repo repositories.Repository) *CooccConfidenceDataProvider {
	return &CooccConfidenceDataProvider{
		log:  log,
		repo: repo,
	}
}

func (p *CooccConfidenceDataProvider) GetData(ctx context.Context, group string) AnalyzedData {
	//filter := bson.M{}
	//if group != "" {
	//	filter = bson.M{"group": group}
	//}
	//messages, err := p.repo.Find(ctx, filter)
	//if err != nil {
	//	p.log.Error("all messages finding", "err", err)
	//}

	return &CooccConfidenceData{}
}
