package analysis

import (
	"context"
	"log/slog"
	"sort"
	"strings"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/meesooqa/ttag/app/repositories"
)

type CooccConfidenceDataProvider struct {
	log  *slog.Logger
	repo repositories.Repository
}

type CooccConfidenceData struct {
	Pairs []TagPairConfidence `json:"pairs"`
}

type TagPairConfidence struct {
	TagA       string  `json:"tagA"`       // Антецедент
	TagB       string  `json:"tagB"`       // Консеквент
	Confidence float64 `json:"confidence"` // Значение confidence
}

func NewCooccConfidenceDataProvider(log *slog.Logger, repo repositories.Repository) *CooccConfidenceDataProvider {
	return &CooccConfidenceDataProvider{
		log:  log,
		repo: repo,
	}
}

// GetData вычисляет метрику Confidence для пар тегов.
// Для каждой пары тегов, найденной в сообщениях, вычисляются два правила:
//   - Confidence(A → B) = f(A,B)/f(A)
//   - Confidence(B → A) = f(A,B)/f(B)
func (p *CooccConfidenceDataProvider) GetData(ctx context.Context, group string) AnalyzedData {
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
		return &CooccConfidenceData{
			Pairs: []TagPairConfidence{},
		}
	}

	// Подсчет индивидуальных частот тегов и частот появления пар (без учета порядка)
	tagFreq := make(map[string]int)
	pairFreq := make(map[string]int)

	for _, msg := range messages {
		// Создаем множество уникальных тегов в сообщении
		tagSet := make(map[string]struct{})
		for _, tag := range msg.Tags {
			tagSet[tag] = struct{}{}
		}
		var tags []string
		for tag := range tagSet {
			tagFreq[tag]++ // учитываем уникальное появление тега в сообщении
			tags = append(tags, tag)
		}
		// Сортировка тегов для формирования ключа пары в единообразном порядке
		sort.Strings(tags)
		for i := 0; i < len(tags); i++ {
			for j := i + 1; j < len(tags); j++ {
				key := tags[i] + "|" + tags[j]
				pairFreq[key]++
			}
		}
	}

	// Формируем результат для каждого направления правила
	var confidencePairs []TagPairConfidence
	for key, fAB := range pairFreq {
		parts := strings.Split(key, "|")
		if len(parts) != 2 {
			continue
		}
		tagA := parts[0]
		tagB := parts[1]
		fA := tagFreq[tagA]
		fB := tagFreq[tagB]
		// Вычисляем Confidence для правила tagA -> tagB, если fA > 0
		if fA > 0 {
			confidence := float64(fAB) / float64(fA)
			confidencePairs = append(confidencePairs, TagPairConfidence{
				TagA:       tagA,
				TagB:       tagB,
				Confidence: confidence,
			})
		}
		// Вычисляем Confidence для правила tagB -> tagA, если fB > 0
		if fB > 0 {
			confidence := float64(fAB) / float64(fB)
			confidencePairs = append(confidencePairs, TagPairConfidence{
				TagA:       tagB,
				TagB:       tagA,
				Confidence: confidence,
			})
		}
	}

	return &CooccConfidenceData{
		Pairs: confidencePairs,
	}
}
