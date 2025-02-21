package adapter

import "github.com/meesooqa/ttag/app/analysis"

type D3DataAdapter interface {
	PrepareData(analyzedData analysis.AnalyzedData) any
}

type D3GraphData struct {
	Nodes []Node `json:"nodes"`
	Links []Link `json:"links"`
}

type Node struct {
	ID    string `json:"id"`
	Count int    `json:"count"`
}

type Link struct {
	Source string `json:"source"`
	Target string `json:"target"`
	Count  int    `json:"count"`
}
