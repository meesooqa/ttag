package analysis

import (
	"context"
	"log/slog"

	"github.com/meesooqa/ttag/app/repositories"
)

type CooccSupportDataProvider struct {
	log  *slog.Logger
	repo repositories.Repository
}

type CooccSupportData struct{}

func NewCooccSupportDataProvider(log *slog.Logger, repo repositories.Repository) *CooccSupportDataProvider {
	return &CooccSupportDataProvider{
		log:  log,
		repo: repo,
	}
}

// Меры ассоциации: Support (поддержка)
// Частота совместного появления пары в совокупности сообщений
func (p *CooccSupportDataProvider) GetData(ctx context.Context, group string) AnalyzedData {
	//filter := bson.M{}
	//if group != "" {
	//	filter = bson.M{"group": group}
	//}
	//messages, err := p.repo.Find(ctx, filter)
	//if err != nil {
	//	p.log.Error("all messages finding", "err", err)
	//}
	// TODO CooccSupportDataProvider

	return &CooccSupportData{}
}
