package analysis

import (
	"container/heap"
	"context"
	"fmt"
	"log/slog"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/meesooqa/ttag/app/repositories"
)

type CooccClustersDataProvider struct {
	log  *slog.Logger
	repo repositories.Repository
}

// структура для дендрограммы
type CooccClustersData struct {
	Name     string               `json:"name"`
	Children []*CooccClustersData `json:"children,omitempty"`
	Size     int                  `json:"size,omitempty"`
	Parent   *CooccClustersData   `json:"-"`
	Tags     []string             `json:"-"`
}

func NewCooccClustersDataProvider(log *slog.Logger, repo repositories.Repository) *CooccClustersDataProvider {
	return &CooccClustersDataProvider{
		log:  log,
		repo: repo,
	}
}

func (p *CooccClustersDataProvider) GetData(ctx context.Context, group string) AnalyzedData {
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

	// Создаем матрицу совместной встречаемости
	coOccurMatrix := make(map[string]map[string]int)
	for _, msg := range messages {
		tags := msg.Tags
		for i := 0; i < len(tags); i++ {
			if _, ok := coOccurMatrix[tags[i]]; !ok {
				coOccurMatrix[tags[i]] = make(map[string]int)
			}
			for j := i + 1; j < len(tags); j++ {
				coOccurMatrix[tags[i]][tags[j]]++
				if _, ok := coOccurMatrix[tags[j]]; !ok {
					coOccurMatrix[tags[j]] = make(map[string]int)
				}
				coOccurMatrix[tags[j]][tags[i]]++
			}
		}
	}

	// Иерархическая кластеризация
	result := p.hierarchicalClustering(coOccurMatrix)

	// Удаляем узлы, которые стали листьями с Size < 2
	return p.filterClusters(result, func(n *CooccClustersData) bool {
		return true
		//return n.Size > 1 || len(n.Children) > 0
		//return n.Size == 2 || len(n.Children) > 0 // пары
		//return n.Size == 3 || len(n.Children) > 0 // трио
	})
}

func (p *CooccClustersDataProvider) hierarchicalClustering(matrix map[string]map[string]int) *CooccClustersData {
	tags := make([]string, 0, len(matrix))
	for tag := range matrix {
		tags = append(tags, tag)
	}
	n := len(tags)

	// 1. Инициализация кластеров
	clusters := make([]*CooccClustersData, n)
	for i := range clusters {
		clusters[i] = &CooccClustersData{
			Name: tags[i],
			Size: 1,
			Tags: []string{tags[i]},
		}
	}

	// 2. Создаем приоритетную очередь
	pq := make(PriorityQueue, 0)
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			sim := p.findMaxSimilarity(clusters[i].Tags, clusters[j].Tags, matrix)
			pq.Push(&Item{
				Value1:   clusters[i],
				Value2:   clusters[j],
				Priority: sim,
			})
		}
	}
	heap.Init(&pq)

	// 3. Объединение кластеров
	for len(clusters) > 1 {
		if pq.Len() == 0 {
			break
		}

		item := heap.Pop(&pq).(*Item)
		a, b := item.Value1, item.Value2

		// Пропускаем уже объединенные кластеры
		if a.Parent != nil || b.Parent != nil {
			continue
		}

		// Создаем новый кластер
		clusterName := fmt.Sprintf("%s+%s", a.Name, b.Name)
		merged := &CooccClustersData{
			Name:     clusterName,
			Children: []*CooccClustersData{a, b},
			Size:     a.Size + b.Size,
			Tags:     append(a.Tags, b.Tags...),
		}
		a.Parent = merged
		b.Parent = merged

		// Удаляем старые кластеры из списка
		newClusters := make([]*CooccClustersData, 0, len(clusters)-1)
		for _, c := range clusters {
			if c != a && c != b {
				newClusters = append(newClusters, c)
			}
		}
		clusters = append(newClusters, merged)

		// Добавляем новые пары в очередь
		for _, c := range newClusters {
			if c != merged {
				sim := p.findMaxSimilarity(merged.Tags, c.Tags, matrix)
				heap.Push(&pq, &Item{
					Value1:   merged,
					Value2:   c,
					Priority: sim,
				})
			}
		}
	}

	return clusters[0]
}

func (p *CooccClustersDataProvider) findMaxSimilarity(tags1, tags2 []string, matrix map[string]map[string]int) float64 {
	maxSim := 0.0
	for _, t1 := range tags1 {
		for _, t2 := range tags2 {
			if sim := float64(matrix[t1][t2]); sim > maxSim {
				maxSim = sim
			}
		}
	}
	return maxSim
}

func (p *CooccClustersDataProvider) filterClusters(node *CooccClustersData, shouldKeep FilterFunc) *CooccClustersData {
	// Рекурсивная фильтрация детей
	if node == nil {
		return nil
	}

	// Рекурсивно фильтруем детей
	var filteredChildren []*CooccClustersData
	for _, child := range node.Children {
		if filteredChild := p.filterClusters(child, shouldKeep); filteredChild != nil {
			filteredChildren = append(filteredChildren, filteredChild)
		}
	}
	node.Children = filteredChildren

	// Проверяем, нужно ли сохранить текущий узел
	if !shouldKeep(node) {
		return nil
	}

	return node
}

// Вспомогательные структуры для оптимизации
type Item struct {
	Value    *CooccClustersData
	Value1   *CooccClustersData
	Value2   *CooccClustersData
	Priority float64
	Index    int
}
type PriorityQueue []*Item

func (pq PriorityQueue) Len() int           { return len(pq) }
func (pq PriorityQueue) Less(i, j int) bool { return pq[i].Priority > pq[j].Priority } // Max-heap
func (pq PriorityQueue) Swap(i, j int)      { pq[i], pq[j] = pq[j], pq[i] }
func (pq *PriorityQueue) Push(x interface{}) {
	item := x.(*Item)
	item.Index = len(*pq)
	*pq = append(*pq, item)
}
func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	item.Index = -1
	*pq = old[0 : n-1]
	return item
}

type FilterFunc func(*CooccClustersData) bool

// http://localhost:8080/api/coocc_clusters_d3.json?group=test&min_size=5&exclude_tags=test,temp
/*
// Фильтр: удаляем узлы с тегами из черного списка
var blacklist = map[string]bool{"test": true, "temp": true}
var filterByName = func(n *CooccClustersData) bool {
	if len(n.Children) == 0 { // Листовой узел
		return !blacklist[n.Name]
	}
	return true
}

// Фильтр: оставляем узлы глубже 2 уровня
var depthFilter = func(n *CooccClustersData) bool {
	return calculateDepth(n) > 2
}
var combinedFilter = func(n *CooccClustersData) bool {
	return filterBySize(n) && filterByName(n)
}

func calculateDepth(n *CooccClustersData) int {
	if len(n.Children) == 0 {
		return 0
	}
	maxDepth := 0
	for _, child := range n.Children {
		if d := calculateDepth(child); d > maxDepth {
			maxDepth = d
		}
	}
	return maxDepth + 1
}
*/
