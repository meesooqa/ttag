package analysis

import (
	"context"
	"log/slog"
	"sort"
	"strings"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/meesooqa/ttag/app/repositories"
)

type CooccJaccardDataProvider struct {
	log  *slog.Logger
	repo repositories.Repository
}

type CooccJaccardData struct {
	Pairs []TagPairJaccard `json:"pairs"`
}

type TagPairJaccard struct {
	TagA    string  `json:"tagA"`
	TagB    string  `json:"tagB"`
	Jaccard float64 `json:"jaccard"`
}

func NewCooccJaccardDataProvider(log *slog.Logger, repo repositories.Repository) *CooccJaccardDataProvider {
	return &CooccJaccardDataProvider{
		log:  log,
		repo: repo,
	}
}

// GetData вычисляет Jaccard Index для пар тегов.
// Jaccard Index определяется как отношение количества сообщений, в которых встречаются оба тега,
// к количеству сообщений, в которых встречается хотя бы один из этих тегов.
func (p *CooccJaccardDataProvider) GetData(ctx context.Context, group string) AnalyzedData {
	filter := bson.M{}
	if group != "" {
		filter = bson.M{"group": group}
	}
	messages, err := p.repo.Find(ctx, filter)
	if err != nil {
		p.log.Error("all messages finding", "err", err)
	}
	totalMessages := len(messages)
	if totalMessages == 0 {
		return &CooccJaccardData{
			Pairs: []TagPairJaccard{},
		}
	}

	tagFreq := make(map[string]int)
	pairFreq := make(map[string]int)

	// Подсчитываем частоты для каждого тега и для каждой пары тегов.
	for _, msg := range messages {
		tagSet := make(map[string]struct{})
		for _, tag := range msg.Tags {
			tagSet[tag] = struct{}{}
		}
		var tags []string
		for tag := range tagSet {
			tagFreq[tag]++
			tags = append(tags, tag)
		}
		sort.Strings(tags)
		for i := 0; i < len(tags); i++ {
			for j := i + 1; j < len(tags); j++ {
				key := tags[i] + "|" + tags[j]
				pairFreq[key]++
			}
		}
	}

	var jaccardPairs []TagPairJaccard
	for key, fAB := range pairFreq {
		parts := strings.Split(key, "|")
		if len(parts) != 2 {
			continue
		}
		tagA := parts[0]
		tagB := parts[1]
		fA := tagFreq[tagA]
		fB := tagFreq[tagB]
		// Вычисляем объединение – количество сообщений, где встречается хотя бы один из тегов
		union := fA + fB - fAB
		if union == 0 {
			continue
		}
		jaccard := float64(fAB) / float64(union)
		jaccardPairs = append(jaccardPairs, TagPairJaccard{
			TagA:    tagA,
			TagB:    tagB,
			Jaccard: jaccard,
		})
	}

	return &CooccJaccardData{
		Pairs: jaccardPairs,
	}
}
