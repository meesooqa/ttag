package adapter

import (
	"log/slog"
	"sort"
	"strings"

	"github.com/meesooqa/ttag/app/analysis"
)

type CooccPairsD3DataAdapter struct {
	log *slog.Logger
}

func NewCooccPairsD3DataAdapter(log *slog.Logger) *CooccPairsD3DataAdapter {
	return &CooccPairsD3DataAdapter{
		log: log,
	}
}

func (a *CooccPairsD3DataAdapter) PrepareData(analyzedData analysis.AnalyzedData) any {
	data, ok := analyzedData.(*analysis.CooccPairsData)
	if !ok {
		a.log.Error("invalid CooccPairs data")
		return nil
	}

	nodes := make([]Node, 0, len(data.TagCounts))
	for tag, count := range data.TagCounts {
		nodes = append(nodes, Node{ID: tag, Count: count})
	}
	links := make([]Link, 0, len(data.CooccPairsMap))
	for pair, count := range data.CooccPairsMap {
		parts := strings.Split(pair, "|")
		links = append(links, Link{
			Source: parts[0],
			Target: parts[1],
			Count:  count,
		})
	}

	// Сортировка узлов и связей
	sort.Slice(nodes, func(i, j int) bool { return nodes[i].Count > nodes[j].Count })
	sort.Slice(links, func(i, j int) bool { return links[i].Count > links[j].Count })

	return &D3GraphData{Nodes: nodes, Links: links}
}
