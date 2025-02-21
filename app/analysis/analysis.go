package analysis

import "context"

type AnalyzedData interface{}

type AnalyzedDataProvider interface {
	GetData(ctx context.Context, group string) AnalyzedData
}
