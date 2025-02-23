package analysis

import (
	"context"
	"fmt"
	"log/slog"
	"math"

	"github.com/meesooqa/ttag/app/repositories"
)

type CooccSupportDataProvider struct {
	*CooccDataProvider
}

func NewCooccSupportDataProvider(log *slog.Logger, repo repositories.Repository) *CooccSupportDataProvider {
	return &CooccSupportDataProvider{&CooccDataProvider{
		log:  log,
		repo: repo,
	}}
}

// GetData calculates меру ассоциации Support (поддержку)
// Частота совместного появления пары в совокупности сообщений
// Support определяется как отношение количества сообщений с данной парой тегов к общему числу сообщений
func (p *CooccSupportDataProvider) GetData(ctx context.Context, group string) AnalyzedData {
	return p.analyzeMessages(ctx, group, []*CooccMetric{{
		Name: "Support",
		Calc: func(fA, fB, fAB, totalMessages int) (float64, error) {
			return float64(fAB) / float64(totalMessages), nil
		},
	}})
}

type CooccPmiDataProvider struct {
	*CooccDataProvider
}

func NewCooccPmiDataProvider(log *slog.Logger, repo repositories.Repository) *CooccPmiDataProvider {
	return &CooccPmiDataProvider{&CooccDataProvider{
		log:  log,
		repo: repo,
	}}
}

// GetData вычисляет PMI (Pointwise Mutual Information) для пар тегов.
// PMI определяется как логарифмическое соотношение между совместной вероятностью появления двух тегов
// и произведением их индивидуальных вероятностей:
// PMI(A, B) = log2((f(A,B) * N) / (f(A) * f(B)))
func (p *CooccPmiDataProvider) GetData(ctx context.Context, group string) AnalyzedData {
	return p.analyzeMessages(ctx, group, []*CooccMetric{{
		Name: "PMI",
		Calc: func(fA, fB, fAB, totalMessages int) (float64, error) {
			if fA == 0 || fB == 0 {
				return 0.00, fmt.Errorf("pmi: fA or fB must not be 0")
			}
			value := (float64(fAB) * float64(totalMessages)) / (float64(fA) * float64(fB))
			return math.Log2(value), nil
		},
	}})
}

type CooccLiftDataProvider struct {
	*CooccDataProvider
}

func NewCooccLiftDataProvider(log *slog.Logger, repo repositories.Repository) *CooccLiftDataProvider {
	return &CooccLiftDataProvider{&CooccDataProvider{
		log:  log,
		repo: repo,
	}}
}

// Меры ассоциации: Lift
// Lift: Отношение наблюдаемой совместной частоты к ожидаемой при независимости появления тегов.
// Lift: отношение совместной вероятности появления тегов к произведению их индивидуальных вероятностей.
func (p *CooccLiftDataProvider) GetData(ctx context.Context, group string) AnalyzedData {
	return p.analyzeMessages(ctx, group, []*CooccMetric{{
		Name: "Lift",
		Calc: func(fA, fB, fAB, totalMessages int) (float64, error) {
			if fA == 0 || fB == 0 {
				return 0.00, fmt.Errorf("lift: fA or fB must not be 0")
			}
			return (float64(fAB) * float64(totalMessages)) / (float64(fA) * float64(fB)), nil
		},
	}})
}

type CooccJaccardDataProvider struct {
	*CooccDataProvider
}

func NewCooccJaccardDataProvider(log *slog.Logger, repo repositories.Repository) *CooccJaccardDataProvider {
	return &CooccJaccardDataProvider{&CooccDataProvider{
		log:  log,
		repo: repo,
	}}
}

// GetData вычисляет Jaccard Index для пар тегов.
// Jaccard Index определяется как отношение количества сообщений, в которых встречаются оба тега,
// к количеству сообщений, в которых встречается хотя бы один из этих тегов.
func (p *CooccJaccardDataProvider) GetData(ctx context.Context, group string) AnalyzedData {
	return p.analyzeMessages(ctx, group, []*CooccMetric{{
		Name: "Jaccard",
		Calc: func(fA, fB, fAB, totalMessages int) (float64, error) {
			union := fA + fB - fAB
			if union == 0 {
				return 0.00, fmt.Errorf("jaccard: union must not be 0")
			}
			return float64(fAB) / float64(union), nil
		},
	}})
}

type CooccConfidenceDataProvider struct {
	*CooccDataProvider
}

func NewCooccConfidenceDataProvider(log *slog.Logger, repo repositories.Repository) *CooccConfidenceDataProvider {
	return &CooccConfidenceDataProvider{&CooccDataProvider{
		log:  log,
		repo: repo,
	}}
}

// GetData вычисляет метрику Confidence для пар тегов.
// Для каждой пары тегов, найденной в сообщениях, вычисляются два правила:
//   - Confidence(A → B) = f(A,B)/f(A)
//   - Confidence(B → A) = f(A,B)/f(B)
func (p *CooccConfidenceDataProvider) GetData(ctx context.Context, group string) AnalyzedData {
	return p.analyzeMessages(ctx, group, []*CooccMetric{{
		Name: "Confidence A",
		Calc: func(fA, fB, fAB, totalMessages int) (float64, error) {
			// Вычисляем Confidence для правила tagA -> tagB, если fA > 0
			if fA > 0 {
				return float64(fAB) / float64(fA), nil
			}
			return 0.00, fmt.Errorf("confidenceA: fA must be greater than 0")
		},
	}, {
		Name: "Confidence B",
		Calc: func(fA, fB, fAB, totalMessages int) (float64, error) {
			// Вычисляем Confidence для правила tagB -> tagA, если fB > 0
			if fB > 0 {
				return float64(fAB) / float64(fB), nil
			}
			return 0.00, fmt.Errorf("confidenceB: fB must be greater than 0")
		},
	}})
}
