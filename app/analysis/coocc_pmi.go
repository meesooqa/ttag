package analysis

import (
	"context"
	"log/slog"
	"math"
	"sort"
	"strings"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/meesooqa/ttag/app/repositories"
)

type CooccPmiDataProvider struct {
	log  *slog.Logger
	repo repositories.Repository
}

type CooccPmiData struct {
	Pairs []TagPairPmi `json:"pairs"`
}

type TagPairPmi struct {
	TagA string  `json:"tagA"`
	TagB string  `json:"tagB"`
	PMI  float64 `json:"pmi"`
}

func NewCooccPmiDataProvider(log *slog.Logger, repo repositories.Repository) *CooccPmiDataProvider {
	return &CooccPmiDataProvider{
		log:  log,
		repo: repo,
	}
}

// GetData вычисляет PMI (Pointwise Mutual Information) для пар тегов.
// PMI определяется как логарифмическое соотношение между совместной вероятностью появления двух тегов
// и произведением их индивидуальных вероятностей:
// PMI(A, B) = log2((f(A,B) * N) / (f(A) * f(B)))
func (p *CooccPmiDataProvider) GetData(ctx context.Context, group string) AnalyzedData {
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
		return &CooccPmiData{
			Pairs: []TagPairPmi{},
		}
	}

	tagFreq := make(map[string]int)
	pairFreq := make(map[string]int)

	// Подсчет частот для каждого тега и для каждой пары тегов.
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

	var pmiPairs []TagPairPmi
	for key, fAB := range pairFreq {
		parts := strings.Split(key, "|")
		if len(parts) != 2 {
			continue
		}
		tagA := parts[0]
		tagB := parts[1]
		fA := tagFreq[tagA]
		fB := tagFreq[tagB]
		// Предотвращаем деление на ноль
		if fA == 0 || fB == 0 {
			continue
		}
		// Вычисляем отношение вероятностей и берем логарифм по основанию 2.
		value := (float64(fAB) * float64(totalMessages)) / (float64(fA) * float64(fB))
		pmi := math.Log2(value)
		pmiPairs = append(pmiPairs, TagPairPmi{
			TagA: tagA,
			TagB: tagB,
			PMI:  pmi,
		})
	}

	return &CooccPmiData{
		Pairs: pmiPairs,
	}
}
