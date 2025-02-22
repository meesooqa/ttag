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
	Pairs   []TagPairLift  `json:"pairs"`
	TagFreq map[string]int `json:"tagFreq"`
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

// Меры ассоциации: Lift
// Lift: Отношение наблюдаемой совместной частоты к ожидаемой при независимости появления тегов.
// Lift: отношение совместной вероятности появления тегов к произведению их индивидуальных вероятностей.
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
		return &CooccLiftData{
			Pairs:   []TagPairLift{},
			TagFreq: map[string]int{},
		}
	}

	// Подсчет частот для каждого тега и совместных появлений для пар тегов.
	tagFreq := make(map[string]int)
	pairFreq := make(map[string]int)

	for _, msg := range messages {
		// Для каждого сообщения создаем множество уникальных тегов.
		tagSet := make(map[string]struct{})
		for _, tag := range msg.Tags {
			tagSet[tag] = struct{}{}
		}

		var tags []string
		for tag := range tagSet {
			tagFreq[tag]++ // Увеличиваем частоту для тега (учитываем только уникальное появление в сообщении)
			tags = append(tags, tag)
		}

		// Сортировка тегов для единообразного формирования ключа пары.
		sort.Strings(tags)
		for i := 0; i < len(tags); i++ {
			for j := i + 1; j < len(tags); j++ {
				key := tags[i] + "|" + tags[j]
				pairFreq[key]++
			}
		}
	}

	// Вычисление Lift для каждой пары.
	var pairs []TagPairLift
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
		pairs = append(pairs, TagPairLift{
			TagA: tagA,
			TagB: tagB,
			Lift: lift,
		})
	}

	return &CooccLiftData{
		Pairs:   pairs,
		TagFreq: tagFreq,
	}
}
