package adapter

import (
	"log/slog"

	"github.com/meesooqa/ttag/app/analysis"
)

type CooccLiftD3DataAdapter struct {
	log *slog.Logger
}

func NewCooccLiftD3DataAdapter(log *slog.Logger) *CooccLiftD3DataAdapter {
	return &CooccLiftD3DataAdapter{
		log: log,
	}
}

func (a *CooccLiftD3DataAdapter) PrepareData(analyzedData analysis.AnalyzedData) any {
	data, ok := analyzedData.(*analysis.CooccLiftData)
	if !ok {
		a.log.Error("invalid CooccLift data")
		return nil
	}
	return data
}
