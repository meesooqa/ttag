package adapter

import (
	"log/slog"

	"github.com/meesooqa/ttag/app/analysis"
)

type CooccSupportD3DataAdapter struct {
	log *slog.Logger
}

func NewCooccSupportD3DataAdapter(log *slog.Logger) *CooccSupportD3DataAdapter {
	return &CooccSupportD3DataAdapter{
		log: log,
	}
}

func (a *CooccSupportD3DataAdapter) PrepareData(analyzedData analysis.AnalyzedData) any {
	data, ok := analyzedData.(*analysis.CooccSupportData)
	if !ok {
		a.log.Error("invalid CooccSupport data")
		return nil
	}
	return data
}
