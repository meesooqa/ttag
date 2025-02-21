package adapter

import (
	"log/slog"

	"github.com/meesooqa/ttag/app/analysis"
)

type CooccClustersD3DataAdapter struct {
	log *slog.Logger
}

func NewCooccClustersD3DataAdapter(log *slog.Logger) *CooccClustersD3DataAdapter {
	return &CooccClustersD3DataAdapter{
		log: log,
	}
}

func (a *CooccClustersD3DataAdapter) PrepareData(analyzedData analysis.AnalyzedData) any {
	data, ok := analyzedData.(*analysis.CooccClustersData)
	if !ok {
		a.log.Error("invalid CooccClusters data")
		return nil
	}
	return data
}
