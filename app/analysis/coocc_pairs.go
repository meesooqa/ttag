package analysis

import (
	"context"
	"fmt"
	"log/slog"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/meesooqa/ttag/app/model"
	"github.com/meesooqa/ttag/app/repositories"
)

type CooccPairsDataProvider struct {
	log  *slog.Logger
	repo repositories.Repository
}

type CooccPairsData struct {
	TagCounts     map[string]int
	CooccPairsMap map[string]int
}

func NewCooccPairsDataProvider(log *slog.Logger, repo repositories.Repository) *CooccPairsDataProvider {
	return &CooccPairsDataProvider{
		log:  log,
		repo: repo,
	}
}

func (p *CooccPairsDataProvider) GetData(ctx context.Context, group string) AnalyzedData {
	filter := bson.M{}
	if group != "" {
		filter = bson.M{"group": group}
	}
	messages, err := p.repo.Find(ctx, filter)
	if err != nil {
		p.log.Error("all messages finding", "err", err)
	}

	tagCounts, cooccPairsMap := p.analyzeMessages(messages)

	return &CooccPairsData{
		TagCounts:     tagCounts,
		CooccPairsMap: cooccPairsMap,
	}
}

func (p *CooccPairsDataProvider) analyzeMessages(messages []*model.Message) (map[string]int, map[string]int) {
	tagCounts := make(map[string]int)
	cooccPairsMap := make(map[string]int)

	for _, msg := range messages {
		tags := msg.Tags
		for _, tag := range tags {
			tagCounts[tag]++
		}
		pairs := p.generatePairs(tags)
		for _, pair := range pairs {
			cooccPairsMap[pair]++
		}
	}

	return tagCounts, cooccPairsMap
}

func (p *CooccPairsDataProvider) generatePairs(tags []string) []string {
	var pairs []string
	seen := make(map[string]bool)

	for i := 0; i < len(tags); i++ {
		for j := i + 1; j < len(tags); j++ {
			tag1, tag2 := tags[i], tags[j]
			if tag1 > tag2 {
				tag1, tag2 = tag2, tag1
			}
			pair := fmt.Sprintf("%s|%s", tag1, tag2)
			if !seen[pair] {
				pairs = append(pairs, pair)
				seen[pair] = true
			}
		}
	}
	return pairs
}
