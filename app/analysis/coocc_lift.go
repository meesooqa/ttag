package analysis

import (
	"context"
	"log/slog"
	"sort"
	"strings"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/meesooqa/ttag/app/repositories"
)

type CooccLiftDataProvider struct {
	log  *slog.Logger
	repo repositories.Repository
}

// CooccLiftData содержит срез пар тегов с их Lift-метрикой.
type CooccLiftData struct {
	Pairs []TagPairLift `json:"pairs"`
}

// TagPairLift содержит пару тегов и рассчитанное значение Lift.
type TagPairLift struct {
	TagA string  `json:"tagA"`
	TagB string  `json:"tagB"`
	Lift float64 `json:"lift"`
}

func NewCooccLiftDataProvider(log *slog.Logger, repo repositories.Repository) *CooccLiftDataProvider {
	return &CooccLiftDataProvider{
		log:  log,
		repo: repo,
	}
}

func (p *CooccLiftDataProvider) GetData(ctx context.Context, group string) AnalyzedData {
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
		return &CooccLiftData{Pairs: []TagPairLift{}}
	}

	// Подсчёт индивидуальных появлений тегов и совместных появлений пар тегов
	tagFreq := make(map[string]int)
	pairFreq := make(map[string]int)

	for _, msg := range messages {
		// Используем множество, чтобы посчитать каждый тег только один раз для данного сообщения
		tagSet := make(map[string]struct{})
		for _, tag := range msg.Tags {
			tagSet[tag] = struct{}{}
		}

		// Подсчёт появления каждого тега
		var tags []string
		for tag := range tagSet {
			tagFreq[tag]++
			tags = append(tags, tag)
		}

		// Сортируем теги, чтобы пары имели единый порядок (A|B, а не B|A)
		sort.Strings(tags)
		for i := 0; i < len(tags); i++ {
			for j := i + 1; j < len(tags); j++ {
				key := tags[i] + "|" + tags[j]
				pairFreq[key]++
			}
		}
	}

	// Вычисление Lift для каждой пары
	var result []TagPairLift
	for key, fAB := range pairFreq {
		parts := strings.Split(key, "|")
		if len(parts) != 2 {
			continue // на всякий случай
		}
		tagA := parts[0]
		tagB := parts[1]
		fA := tagFreq[tagA]
		fB := tagFreq[tagB]

		// Предотвращаем деление на ноль
		if fA == 0 || fB == 0 {
			continue
		}
		lift := (float64(fAB) * float64(totalMessages)) / (float64(fA) * float64(fB))
		result = append(result, TagPairLift{
			TagA: tagA,
			TagB: tagB,
			Lift: lift,
		})
	}

	return &CooccLiftData{
		Pairs: result,
	}
}
