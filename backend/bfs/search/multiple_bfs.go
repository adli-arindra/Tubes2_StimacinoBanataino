package search

import (
	"fmt"
	"runtime"
	"sync"
	"time"

	"bfs/graph"
)

// Menyimpan task kombinasi 
type levelTask struct {
	Product string
	Recipe  []string
}

// Menyimpan result dari task yang berhasil
type levelResult struct {
	Product string
	Node    *graph.TreeNode
}

func MultiBFS(target string, g graph.Graph, maxRecipes int, tierMap map[string]int) (graph.MultiTreeResult, error) {
	start := time.Now()

	startingElements := []string{"Air", "Fire", "Water", "Earth"}
	existing := sync.Map{}
	visitedCombo := sync.Map{}
	visitedElem := sync.Map{}

	// Inisialisasi dari starting element
	for _, e := range startingElements {
		existing.Store(e, &graph.TreeNode{Name: e, Children: []*graph.TreeNode{}})
		visitedElem.Store(e, true)
	}

	foundRecipes := []*graph.TreeNode{}
	tier := 0
	numWorkers := runtime.NumCPU() * 2

	// Proses Multiple BFS
	for len(foundRecipes) < maxRecipes {
		tasks := make(chan levelTask, 1000)
		results := make(chan levelResult, 1000)
		var wg sync.WaitGroup

		// Worker pool multithreading
		for i := 0; i < numWorkers; i++ {
			wg.Add(1)
			go func(id int) {
				defer wg.Done()
				for task := range tasks {
					product := task.Product
					recipe := task.Recipe

					// Pengecekan untuk resep harus di bawah elemen yang akan dibentuk
					productTier := tierMap[product]
					valid := true
					for _, ing := range recipe {
						if tierMap[ing] >= productTier {
							valid = false
							break
						}
					}
					if !valid {
						continue
					}

					// Pengecekan elemen dari recipe
					n1Raw, ok1 := existing.Load(recipe[0])
					n2Raw, ok2 := existing.Load(recipe[1])
					if !ok1 || !ok2 {
						continue
					}

					// Menghindari duplikasi kalo ada kombinasi yang sama
					comboKey := recipe[0] + "+" + recipe[1] + ">" + product
					if _, dup := visitedCombo.LoadOrStore(comboKey, true); dup {
						continue
					}

					newNode := &graph.TreeNode{Name: product, Children: []*graph.TreeNode{n1Raw.(*graph.TreeNode), n2Raw.(*graph.TreeNode)}}
					results <- levelResult{Product: product, Node: newNode}
					fmt.Printf("Tier-%d %s dibuat dari %v\n", tier, product, recipe) // debugging
				}
			}(i)
		}

		// Mengirim semua kombinasi ke worker
		for product, recipes := range g {
			if _, seen := visitedElem.Load(product); seen {
				continue
			}
			for _, r := range recipes {
				if len(r) != 2 {
					continue
				}
				tasks <- levelTask{Product: product, Recipe: r}
			}
		}
		close(tasks)

		go func() {
			wg.Wait()
			close(results)
		}()

		nextCount := 0
		for res := range results {
			existing.Store(res.Product, res.Node)
			visitedElem.Store(res.Product, true)
			nextCount++
			
			if res.Product == target {
				foundRecipes = append(foundRecipes, res.Node)
				if len(foundRecipes) >= maxRecipes {
					break
				}
			}
		}

		if nextCount == 0 {
			break
		}
		tier++
	}

	nodeCount := 0
	visitedElem.Range(func(_, _ any) bool {
		nodeCount++
		return true
	})

	duration := float64(time.Since(start).Microseconds()) / 1000.0
	return graph.MultiTreeResult{
		Trees:        foundRecipes,
		Algorithm:    "Multi_BFS",
		DurationMS:   duration,
		VisitedNodes: nodeCount,
	}, nil
}