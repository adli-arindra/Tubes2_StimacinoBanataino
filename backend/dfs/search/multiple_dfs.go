package search

import (
	"sync"
	"runtime"
	"time"
	"dfs/graph"
)

// Menyimpan dua elemen pembentuk
type Recipe struct {
	Recipe1	string
	Recipe2	string
}

func MultiDFS(target string, g graph.Graph, maxRecipes int, tierMap map[string]int) (graph.MultiTreeResult, error) {
	start := time.Now()

	memo := map[string]map[string]Recipe{} // Nyimpen memo untuk DFS

	primaryPath, ok := dfsFromTarget(target, g, tierMap, memo)

	if !ok {
		return graph.MultiTreeResult{
			Algorithm: "Multi_DFS", 
			Trees: nil,
		}, nil
	}

	primaryKey := generateKey(primaryPath) // Membuat key unik dari primary path biar tidak ada duplikasi
	trees := []*graph.TreeNode{buildTreeFromMap(target, primaryPath)}

	// Menemukan node-node dalam primary path
	elementsToExplore := findElementsWithAlternatives(target, primaryPath, g, tierMap)

	// Eksplorasi alternative
	altRecipes := exploreAlternatives(primaryPath, g, tierMap, maxRecipes-1, elementsToExplore, primaryKey) // Eksplorasi recipe alternative

	// Membuat key untuk masing-masing recipe alternative
	for _, recipeMap := range altRecipes {
		if generateKey(recipeMap) == primaryKey {
			continue
		}
		trees = append(trees, buildTreeFromMap(target, recipeMap))
	}

	// Set primary route untuk live update
	for i, tree := range trees {
		if i == 0 {
			idx := 0
			setDiscoveredIndexMultiple(tree, &idx)
		} else {
			markTreeAsAlternative(tree)
		}
	}

	visited := map[string]bool{}

	for _, recipe := range append([]map[string]Recipe{primaryPath}, altRecipes...) {
		for product, r := range recipe {
			visited[product] = true
			visited[r.Recipe1] = true
			visited[r.Recipe2] = true
		}
	}
	for _, base := range []string{"Air", "Fire", "Water", "Earth"} {
		visited[base] = true
	}

	duration := float64(time.Since(start).Microseconds()) / 1000.0
	return graph.MultiTreeResult{
		Trees:        trees,
		Algorithm:    "Multi_DFS",
		DurationMS:   duration,
		VisitedNodes: len(visited), // Bentar ini harusnya semua node
	}, nil
}

func dfsFromTarget(cur string, g graph.Graph, tierMap map[string]int, memo map[string]map[string]Recipe) (map[string]Recipe, bool) {
	base := map[string]bool{"Air": true, "Fire": true, "Water": true, "Earth": true}

	// Basis dari DFS
	if base[cur] {
		return map[string]Recipe{}, true
	}
	if m, ok := memo[cur]; ok {
		return m, true
	}

	// Iteratating semua kombinasi
	for _, r := range g[cur] {
		// Validasi elemen pembentuk
		if len(r) != 2 || tierMap[r[0]] >= tierMap[cur] || tierMap[r[1]] >= tierMap[cur] {
			continue
		}

		// Lakukan DFS untuk kedua elemen pembentuk
		left, ok1 := dfsFromTarget(r[0], g, tierMap, memo)
		right, ok2 := dfsFromTarget(r[1], g, tierMap, memo)

		if ok1 && ok2 {
			m := map[string]Recipe{
				cur: {Recipe1: r[0], Recipe2: r[1]},
			}
			for k, v := range left {
				m[k] = v
			}
			for k, v := range right {
				m[k] = v
			}
			memo[cur] = m
			return m, true
		}
	}
	return nil, false
}

func findElementsWithAlternatives(root string, recipe map[string]Recipe, g graph.Graph, tierMap map[string]int) []string {
	visited := map[string]bool{}
	var result []string

	var dfs func(string)
	dfs = func(cur string) {
		if visited[cur] {
			return
		}
		visited[cur] = true

		r, ok := recipe[cur]
		if !ok {
			return
		}

		count := 0
		// Mengecek alternatif untuk resep lain
		for _, alt := range g[cur] {
			if len(alt) != 2 {
				continue
			}
			if (alt[0] == r.Recipe1 && alt[1] == r.Recipe2) || (alt[1] == r.Recipe1 && alt[0] == r.Recipe2) {
				continue
			}
			if tierMap[alt[0]] >= tierMap[cur] || tierMap[alt[1]] >= tierMap[cur] {
				continue
			}
			count++
			if count > 0 {
				result = append(result, cur)
				break
			}
		}

		dfs(r.Recipe1)
		dfs(r.Recipe2)
	}

	dfs(root)
	return result
}

func exploreAlternatives(base map[string]Recipe, g graph.Graph, tierMap map[string]int, limit int, elementsToExplore []string, primarySig string) []map[string]Recipe {
	if len(elementsToExplore) == 0 || limit <= 0 {
		return nil
	}

	seen := sync.Map{}
	seen.Store(primarySig, true)

	results := make([]map[string]Recipe, 0, limit)
	pending := make(chan map[string]Recipe, limit*2)
	output := make(chan map[string]Recipe, limit)

	var wg sync.WaitGroup

	numWorkers := runtime.NumCPU()
	if numWorkers < 1 {
		numWorkers = 1
	}

	pending <- base

	worker := func() {
		defer wg.Done()
		for current := range pending {
			for _, elem := range elementsToExplore {
				original, exists := current[elem]
				if !exists {
					continue
				}
				for _, alt := range g[elem] {
					if len(alt) != 2 {
						continue
					}
					if (alt[0] == original.Recipe1 && alt[1] == original.Recipe2) ||
						(alt[1] == original.Recipe1 && alt[0] == original.Recipe2) {
						continue
					}
					if tierMap[alt[0]] >= tierMap[elem] || tierMap[alt[1]] >= tierMap[elem] {
						continue
					}

					altMap := copyRecipeMap(current)
					altMap[elem] = Recipe{Recipe1: alt[0], Recipe2: alt[1]}

					if validateAndRepairRecipe(elem, altMap, g, tierMap) {
						sig := generateKey(altMap)
						if _, exists := seen.LoadOrStore(sig, true); !exists {
							output <- altMap
							if len(results) < limit {
								pending <- altMap
							}
						}
					}
				}
			}
		}
	}

	// Menjalankan worker
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go worker()
	}

	// Mengambil hasil dari output channel
	go func() {
		wg.Wait()
		close(output)
	}()

	for recipe := range output {
		results = append(results, recipe)
		if len(results) >= limit {
			break
		}
	}

	return results
}

// Pengecekan apakah jalur/rute valid atau tidak, jika tidak valid akan dicoba ulang.
func validateAndRepairRecipe(root string, m map[string]Recipe, g graph.Graph, tierMap map[string]int) bool {
	base := map[string]bool{"Air": true, "Fire": true, "Water": true, "Earth": true}

	var dfs func(string) bool
	dfs = func(cur string) bool {
		if base[cur] {
			return true
		}

		r, ok := m[cur]
		if !ok {
			repaired, ok := dfsFromTarget(cur, g, tierMap, map[string]map[string]Recipe{})
			if !ok {
				return false
			}
			for k, v := range repaired {
				m[k] = v
			}
			r = m[cur]
		}

		if tierMap[r.Recipe1] >= tierMap[cur] || tierMap[r.Recipe2] >= tierMap[cur] {
			return false
		}
		return dfs(r.Recipe1) && dfs(r.Recipe2)
	}
	return dfs(root)
}

func buildTreeFromMap(cur string, m map[string]Recipe) *graph.TreeNode {
	n := &graph.TreeNode{Name: cur, Children: []*graph.TreeNode{}}
	if r, ok := m[cur]; ok {
		n.Children = append(n.Children, buildTreeFromMap(r.Recipe1, m))
		n.Children = append(n.Children, buildTreeFromMap(r.Recipe2, m))
	}
	return n
}

// Sama kayak BFS Multiple
func setDiscoveredIndexMultiple(node *graph.TreeNode, counter *int) {
	if node == nil {
		return
	}
	node.NodeDiscovered = *counter
	*counter++
	for _, child := range node.Children {
		setDiscoveredIndex(child, counter)
	}
}

// Sama kayak BFS Multiple
func markTreeAsAlternative(node *graph.TreeNode) {
	if node == nil {
		return
	}
	node.NodeDiscovered = -1
	for _, child := range node.Children {
		markTreeAsAlternative(child)
	}
}

func copyRecipeMap(src map[string]Recipe) map[string]Recipe {
	dst := make(map[string]Recipe)
	for k, v := range src {
		dst[k] = v
	}
	return dst
}

func generateKey(m map[string]Recipe) string {
	s := ""
	for k, v := range m {
		s += k + ":" + v.Recipe1 + "+" + v.Recipe2 + "|"
	}
	return s
}