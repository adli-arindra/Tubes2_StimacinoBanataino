package search

import (
	"fmt"
	"runtime"
	"sync"
	"time"

	"bfs/graph"
)

// Menyimpan tiap kombinasi resep
type levelTask struct {
	Product string
	Recipe  []string
}

func MultiBFS(target string, g graph.Graph, maxRecipes int, tierMap map[string]int) (graph.MultiTreeResult, error) {
	start := time.Now()

	startingElements := []string{"Air", "Fire", "Water", "Earth"}
	existing := map[string][]*graph.TreeNode{}
	visitedCombo := map[string]bool{}
	visitedElem := map[string]bool{}

	// Inisialisasi dari starting element
	for _, e := range startingElements {
		existing[e] = []*graph.TreeNode{{Name: e, Children: []*graph.TreeNode{}}}
		visitedElem[e] = true
	}

	// Jika target elemen dasar
	if contains(startingElements, target) {
		n := &graph.TreeNode{
			Name: target, 
			NodeDiscovered: 0,
			Children: []*graph.TreeNode{},
		}
		return graph.MultiTreeResult{
		Trees:        []*graph.TreeNode{n},
		Algorithm:    "Multi_BFS",
		DurationMS:   float64(time.Since(start).Microseconds()) / 1000.0,
		VisitedNodes: 1,
		}, nil
	}

	foundRecipes := []*graph.TreeNode{}
	tier := 0
	numWorkers := runtime.NumCPU() * 2
	targetTier := tierMap[target]
	done := false

	// Proses utama Multiple BFS
	for !done {
		newNodesThisTier := map[string][]*graph.TreeNode{}
		var mutex sync.Mutex
		tasks := make(chan levelTask, 1000)
		var wg sync.WaitGroup

		// Worker pool untuk memproses recipe secara paralel
		for i := 0; i < numWorkers; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for task := range tasks {
					if maxRecipes > 0 && len(foundRecipes) >= maxRecipes {
						return
					}

					product := task.Product
					recipe := task.Recipe
					productTier := tierMap[product]

					// Cek elemen pembentuk tidak lebih tinggi dari elemen yang dibentuk
					if tierMap[recipe[0]] >= productTier || tierMap[recipe[1]] >= productTier {
						continue
					}

					lefts, ok1 := existing[recipe[0]]
					rights, ok2 := existing[recipe[1]]
					if !ok1 || !ok2 {
						continue
					}

					// Build semua kombinasi dari elemen yang sudah ada
					for _, l := range lefts {
						for _, r := range rights {
							key := fmt.Sprintf("%s+%s>%s#%p|%p", recipe[0], recipe[1], product, l, r)
							mutex.Lock()
							if visitedCombo[key] {
								mutex.Unlock()
								continue
							}
							
							visitedCombo[key] = true
							mutex.Unlock()

       						newNode := &graph.TreeNode{Name: product, Children: []*graph.TreeNode{l, r}}

							mutex.Lock()
							newNodesThisTier[product] = append(newNodesThisTier[product], newNode)
							if product == target {
								foundRecipes = append(foundRecipes, deepCopyTree(newNode))
								if maxRecipes > 0 && len(foundRecipes) >= maxRecipes {
									done = true
        						}
      						}
       						mutex.Unlock()
      					}
     				}
    			}
   			}()
  		}

		// Enqueue semua recipe yang valid
  		for product, recipes := range g {
			productTier := tierMap[product]
			if productTier > targetTier {
				continue
			}
			if productTier == targetTier && product != target {
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
		wg.Wait()

		if len(newNodesThisTier) == 0 {
			break
		}

		for k, v := range newNodesThisTier {
			existing[k] = append(existing[k], v...)
			visitedElem[k] = true
		}

		tier++
	}

	// Ambil rute pertama sebagai primary rute sebagai live update
	for i, t := range foundRecipes {
		if i == 0 {
			idx := 0
			setDiscoveredIndexMultiple(t, &idx)
		} else {
			markTreeAsAlternative(t)
		}
	}

	// Cek apakah resep sudah cukup atau belum
	if maxRecipes > 0 && len(foundRecipes) > maxRecipes {
		foundRecipes = foundRecipes[:maxRecipes]
	}

	return graph.MultiTreeResult{
		Trees:        foundRecipes,
		Algorithm:    "Multi_BFS_All_Paths",
		DurationMS:   float64(time.Since(start).Microseconds()) / 1000.0,
		VisitedNodes: len(visitedElem),
	}, nil
}

// Membuat salinan dari pohon node
func deepCopyTree(node *graph.TreeNode) *graph.TreeNode {
	if node == nil {
		return nil
	}

	copy := &graph.TreeNode{
		Name: node.Name,
		NodeDiscovered: node.NodeDiscovered,
		Children: []*graph.TreeNode{},
	}
	for _, child := range node.Children {
		copy.Children = append(copy.Children, deepCopyTree(child))
	}
	return copy
}

// Sama kayak di BFS
func setDiscoveredIndexMultiple(node *graph.TreeNode, counter *int) {
if node == nil {
		return
	}

	queue := []*graph.TreeNode{node}

	for len(queue) > 0 {
		curr := queue[0]
		queue = queue[1:]

		curr.NodeDiscovered = *counter
		*counter++

		queue = append(queue, curr.Children...)
	}
}

// yang alternative nilai node discoverednya adalah -1
func markTreeAsAlternative(node *graph.TreeNode) {
	if node == nil {
		return
	}
	for _, child := range node.Children {
		markTreeAsAlternative(child)
	}
	node.NodeDiscovered = -1
}

// Contains untuk memeriksa apakah elemen ada di dalam list rating atau tidak
func contains(list []string, val string) bool {
	for _, v := range list {
		if v == val {
			return true
		}
	}
	return false
}