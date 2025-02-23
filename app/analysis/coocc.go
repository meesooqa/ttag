package analysis

import (
	"context"
	"log/slog"
	"sort"
	"strings"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/meesooqa/ttag/app/model"
	"github.com/meesooqa/ttag/app/repositories"
)

type CooccMetric struct {
	Name string
	Calc func(fA, fB, fAB, totalMessages int) (float64, error)
}

type CooccDataProvider struct {
	log  *slog.Logger
	repo repositories.Repository
}

type CooccData struct {
	Pairs   []TagPair
	TagFreq map[string]int
}

type TagPair struct {
	TagA  string
	TagB  string
	Value float64
}

type freq struct {
	tagFreq  map[string]int
	pairFreq map[string]int
}

func NewCooccDataProvider(log *slog.Logger, repo repositories.Repository) *CooccDataProvider {
	return &CooccDataProvider{
		log:  log,
		repo: repo,
	}
}

func (p *CooccDataProvider) GetData(ctx context.Context, group string) AnalyzedData {
	return &CooccData{Pairs: []TagPair{}, TagFreq: map[string]int{}}
}

func (p *CooccDataProvider) getMessages(ctx context.Context, group string) []*model.Message {
	filter := bson.M{}
	if group != "" {
		filter = bson.M{"group": group}
	}
	messages, err := p.repo.Find(ctx, filter)
	if err != nil {
		p.log.Error("all messages finding", "err", err)
	}
	return messages
}

func (p *CooccDataProvider) analyzeMessages(ctx context.Context, group string, metrics []*CooccMetric) AnalyzedData {
	messages := p.getMessages(ctx, group)
	if len(messages) == 0 {
		return &CooccData{Pairs: []TagPair{}, TagFreq: map[string]int{}}
	}

	freqData := p.getFreqData(messages)
	totalMessages := len(messages)
	var pairs []TagPair
	for key, fAB := range freqData.pairFreq {
		parts := strings.Split(key, "|")
		if len(parts) != 2 {
			continue // на всякий случай
		}
		tagA, tagB := parts[0], parts[1]
		fA, fB := freqData.tagFreq[tagA], freqData.tagFreq[tagB]

		for _, metric := range metrics {
			value, err := metric.Calc(fA, fB, fAB, totalMessages)
			if err != nil {
				p.log.Error("calculateValue", "err", err)
			}
			pairs = append(pairs, TagPair{
				TagA:  tagA,
				TagB:  tagB,
				Value: value,
			})
		}
	}
	return &CooccData{Pairs: pairs, TagFreq: freqData.tagFreq}
}

func (p *CooccDataProvider) getFreqData(messages []*model.Message) *freq {
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
	return &freq{
		tagFreq:  tagFreq,
		pairFreq: pairFreq,
	}
}
