package adapter

import (
	"log/slog"

	"github.com/meesooqa/ttag/app/analysis"
)

type CooccMatrixD3DataAdapter struct {
	log *slog.Logger
}

func NewCooccMatrixD3DataAdapter(log *slog.Logger) *CooccMatrixD3DataAdapter {
	return &CooccMatrixD3DataAdapter{
		log: log,
	}
}

func (a *CooccMatrixD3DataAdapter) PrepareData(analyzedData analysis.AnalyzedData) any {
	data, ok := analyzedData.(*analysis.CooccMatrixData)
	if !ok {
		a.log.Error("invalid CooccMatrix data")
		return nil
	}
	return data
}
