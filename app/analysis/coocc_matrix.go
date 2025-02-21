package analysis

import (
	"context"
	"log/slog"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/meesooqa/ttag/app/repositories"
)

type CooccMatrixDataProvider struct {
	log  *slog.Logger
	repo repositories.Repository
}

type CooccMatrixData struct {
	Tags   []string `json:"tags"`
	Matrix [][]int  `json:"matrix"`
}

func NewCooccMatrixDataProvider(log *slog.Logger, repo repositories.Repository) *CooccMatrixDataProvider {
	return &CooccMatrixDataProvider{
		log:  log,
		repo: repo,
	}
}

func (p *CooccMatrixDataProvider) GetData(ctx context.Context, group string) AnalyzedData {
	filter := bson.M{}
	if group != "" {
		filter = bson.M{"group": group}
	}
	messages, err := p.repo.Find(ctx, filter)
	if err != nil {
		p.log.Error("all messages finding", "err", err)
	}
	//messages := []model.Message{
	//	{Tags: []string{"tag1", "tag2"}},
	//	{Tags: []string{"tag1", "tag3"}},
	//	{Tags: []string{"tag1", "tag4"}},
	//	{Tags: []string{"tag2", "tag4"}},
	//}

	// Строим матрицу совместной встречаемости
	tagIndex := make(map[string]int)
	var tags []string
	matrix := make(map[string]map[string]int)

	// Собираем уникальные теги
	for _, msg := range messages {
		for _, tag := range msg.Tags {
			if _, exists := tagIndex[tag]; !exists {
				tagIndex[tag] = len(tags)
				tags = append(tags, tag)
				matrix[tag] = make(map[string]int)
			}
		}
	}

	// Заполняем матрицу
	for _, msg := range messages {
		for i := 0; i < len(msg.Tags); i++ {
			for j := i; j < len(msg.Tags); j++ {
				tag1 := msg.Tags[i]
				tag2 := msg.Tags[j]
				matrix[tag1][tag2]++
				if i != j {
					matrix[tag2][tag1]++
				}
			}
		}
	}

	// Преобразуем в массив
	size := len(tags)
	result := make([][]int, size)
	for i := range result {
		result[i] = make([]int, size)
		tagI := tags[i]
		for j := range result[i] {
			result[i][j] = matrix[tagI][tags[j]]
		}
	}

	// Отправляем ответ
	return &CooccMatrixData{
		Tags:   tags,
		Matrix: result,
	}
}
