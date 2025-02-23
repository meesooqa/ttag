package adapter

import (
	"log/slog"

	"github.com/meesooqa/ttag/app/analysis"
)

type CooccConfidenceD3DataAdapter struct {
	log *slog.Logger
}

func NewCooccConfidenceD3DataAdapter(log *slog.Logger) *CooccConfidenceD3DataAdapter {
	return &CooccConfidenceD3DataAdapter{
		log: log,
	}
}

func (a *CooccConfidenceD3DataAdapter) PrepareData(analyzedData analysis.AnalyzedData) any {
	data, ok := analyzedData.(*analysis.CooccConfidenceData)
	if !ok {
		a.log.Error("invalid CooccConfidence data")
		return nil
	}
	return data
}
