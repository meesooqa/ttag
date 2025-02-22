package adapter

import (
	"log/slog"

	"github.com/meesooqa/ttag/app/analysis"
)

type CooccPmiD3DataAdapter struct {
	log *slog.Logger
}

func NewCooccPmiD3DataAdapter(log *slog.Logger) *CooccPmiD3DataAdapter {
	return &CooccPmiD3DataAdapter{
		log: log,
	}
}

func (a *CooccPmiD3DataAdapter) PrepareData(analyzedData analysis.AnalyzedData) any {
	data, ok := analyzedData.(*analysis.CooccPmiData)
	if !ok {
		a.log.Error("invalid CooccPmi data")
		return nil
	}
	return data
}
