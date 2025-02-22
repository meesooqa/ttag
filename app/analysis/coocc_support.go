package analysis

import (
	"context"
	"log/slog"
	"sort"
	"strings"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/meesooqa/ttag/app/repositories"
)

type CooccSupportDataProvider struct {
	log  *slog.Logger
	repo repositories.Repository
}

type CooccSupportData struct {
	Pairs   []TagPairSupport `json:"pairs"`
	TagFreq map[string]int   `json:"tagFreq"`
}

type TagPairSupport struct {
	TagA    string  `json:"tagA"`
	TagB    string  `json:"tagB"`
	Support float64 `json:"support"`
}

func NewCooccSupportDataProvider(log *slog.Logger, repo repositories.Repository) *CooccSupportDataProvider {
	return &CooccSupportDataProvider{
		log:  log,
		repo: repo,
	}
}

// Меры ассоциации: Support (поддержка)
// Частота совместного появления пары в совокупности сообщений
// Support определяется как отношение количества сообщений с данной парой тегов к общему числу сообщений
func (p *CooccSupportDataProvider) GetData(ctx context.Context, group string) AnalyzedData {
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
		return &CooccSupportData{
			Pairs: []TagPairSupport{},
		}
	}

	// Подсчет совместного появления тегов.
	tagFreq := make(map[string]int)
	pairFreq := make(map[string]int)
	for _, msg := range messages {
		// Создаем множество уникальных тегов для каждого сообщения.
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

	// Формирование результатов с расчетом Support.
	var supportPairs []TagPairSupport
	for key, fAB := range pairFreq {
		parts := strings.Split(key, "|")
		if len(parts) != 2 {
			continue
		}
		support := float64(fAB) / float64(totalMessages)
		supportPairs = append(supportPairs, TagPairSupport{
			TagA:    parts[0],
			TagB:    parts[1],
			Support: support,
		})
	}

	return &CooccSupportData{
		Pairs:   supportPairs,
		TagFreq: tagFreq,
	}
}
