package adapter

import (
	"log/slog"

	"github.com/meesooqa/ttag/app/analysis"
)

type CooccJaccardD3DataAdapter struct {
	log *slog.Logger
}

func NewCooccJaccardD3DataAdapter(log *slog.Logger) *CooccJaccardD3DataAdapter {
	return &CooccJaccardD3DataAdapter{
		log: log,
	}
}

func (a *CooccJaccardD3DataAdapter) PrepareData(analyzedData analysis.AnalyzedData) any {
	data, ok := analyzedData.(*analysis.CooccJaccardData)
	if !ok {
		a.log.Error("invalid CooccJaccard data")
		return nil
	}
	return data
}
