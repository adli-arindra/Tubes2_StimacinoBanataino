package search

import (
    "sync"
    "runtime"
	"time"
	"dfs/graph"
)

// Menyimpan dua elemen pembentuk
type Recipe struct {
	Source  string
	Partner string
}

func MultiDFS(target string, g graph.Graph, maxRecipes int, tierMap map[string]int) (graph.MultiTreeResult, error) {
	start := time.Now()

    // menyimpan hasil DFS dari tiap element
	memo := map[string]map[string]Recipe{}

    // Primary route
	mainPath, ok := dfsFromTarget(target, g, tierMap, memo)
	if !ok {
		return graph.MultiTreeResult{Algorithm: "Multi_DFS", Trees: nil}, nil
	}

	primarySig := generateSignature(mainPath)
	trees := []*graph.TreeNode{buildTreeFromMap(target, mainPath)}
	altRecipes := exploreAlternativesParallel(mainPath, g, tierMap, maxRecipes-1)

	for _, recipeMap := range altRecipes {
        // Jika ada duplikat dengan primary route, skip ae
		if generateSignature(recipeMap) == primarySig {
			continue
		}
		trees = append(trees, buildTreeFromMap(target, recipeMap))
	}

	for i, tree := range trees {
		if i == 0 {
			idx := 0
			setDiscoveredIndexMultiple(tree, &idx) // untuk primary route
		} else {
			markNodeDiscoveredMinusOne(tree) // untuk alternative route
		}
	}

	duration := float64(time.Since(start).Microseconds()) / 1000.0
	return graph.MultiTreeResult{
		Trees:        trees,
		Algorithm:    "Multi_DFS",
		DurationMS:   duration,
		VisitedNodes: len(mainPath),
	}, nil
}

func dfsFromTarget(cur string, g graph.Graph, tierMap map[string]int, memo map[string]map[string]Recipe) (map[string]Recipe, bool) {
	base := map[string]bool{"Air": true, "Fire": true, "Water": true, "Earth": true}

    // Jika starting element, tidak ada resep
	if base[cur] {
		return map[string]Recipe{}, true
	}

    // Jika sudah pernah ditemui, ambil dari memo
	if m, ok := memo[cur]; ok {
		return m, true
	}

    // Coba semua resep untuk membentuk current element
	for _, r := range g[cur] {
		if len(r) != 2 || tierMap[r[0]] >= tierMap[cur] || tierMap[r[1]] >= tierMap[cur] {
			continue
		}

        // DFS untuk kedua elemen pembentuk
		left, ok1 := dfsFromTarget(r[0], g, tierMap, memo)
		right, ok2 := dfsFromTarget(r[1], g, tierMap, memo)

        // Menggabungkan hasil jika dapat dibentuk
		if ok1 && ok2 {
			m := map[string]Recipe{
				cur: {Source: r[0], Partner: r[1]},
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

func exploreAlternativesParallel(base map[string]Recipe, g graph.Graph, tierMap map[string]int, limit int) []map[string]Recipe {
	seen := sync.Map{} // Menyimpan kombinasi recipe yang sudah didapatkan
	seen.Store(generateSignature(base), true)

	results := make([]map[string]Recipe, 0, limit)
	pending := make(chan map[string]Recipe, limit*2)
	output := make(chan map[string]Recipe, limit)
	var wg sync.WaitGroup

	pending <- base

    // worker untuk eksplorasi kombinasi baru
	worker := func() {
		for current := range pending {
			for elem, original := range current {
				for _, alt := range g[elem] {
					if len(alt) != 2 {
						continue
					}
                    // Jika kombinasi alternatif sama dengan asli, skip ae
					if (alt[0] == original.Source && alt[1] == original.Partner) ||
						(alt[1] == original.Source && alt[0] == original.Partner) {
						continue
					}
                    // Jika tier element pembentuk >= tier produk, skip ae
					if tierMap[alt[0]] >= tierMap[elem] || tierMap[alt[1]] >= tierMap[elem] {
						continue
					}
                    // Membuat salinan dari recipe
					altMap := copyRecipeMap(current)
					altMap[elem] = Recipe{Source: alt[0], Partner: alt[1]}

					if validateAndRepairRecipe(elem, altMap, g, tierMap) {
						sig := generateSignature(altMap)
						if _, exists := seen.LoadOrStore(sig, true); !exists {
							output <- altMap
							pending <- altMap
						}
					}
				}
			}
		}
		wg.Done()
	}


	numWorkers := runtime.NumCPU() * 2
	wg.Add(numWorkers)
	for i := 0; i < numWorkers; i++ {
		go worker()
	}

    // goroutine untuk menutup chnnael setelah worker selesai
	go func() {
		wg.Wait()
		close(output)
		close(pending)
	}()

	for recipe := range output {
		results = append(results, recipe)
		if len(results) >= limit {
			break
		}
	}

	return results
}

func validateAndRepairRecipe(root string, m map[string]Recipe, g graph.Graph, tierMap map[string]int) bool {
	base := map[string]bool{"Air": true, "Fire": true, "Water": true, "Earth": true}
	var dfs func(string) bool
	dfs = func(cur string) bool {
		if base[cur] {
			return true
		}

		r, ok := m[cur]
		if !ok {
            // Membangun node dengan DFS jika belum ada di map
			repaired, ok := dfsFromTarget(cur, g, tierMap, map[string]map[string]Recipe{})
			if !ok {
				return false
			}
			for k, v := range repaired {
				m[k] = v
			}
			r = m[cur]
		}

        // Pengecekan urutan tier
		if tierMap[r.Source] >= tierMap[cur] || tierMap[r.Partner] >= tierMap[cur] {
			return false
		}
		return dfs(r.Source) && dfs(r.Partner)
	}
	return dfs(root)
}

func buildTreeFromMap(cur string, m map[string]Recipe) *graph.TreeNode {
	n := &graph.TreeNode{Name: cur, Children: []*graph.TreeNode{}}
	if r, ok := m[cur]; ok {
		n.Children = append(n.Children, buildTreeFromMap(r.Source, m))
		n.Children = append(n.Children, buildTreeFromMap(r.Partner, m))
	}
	return n
}

func setDiscoveredIndexMultiple(node *graph.TreeNode, counter *int) {
	if node == nil {
		return
	}
	node.NodeDiscovered = *counter
	*counter++
	for _, child := range node.Children {
		setDiscoveredIndexMultiple(child, counter)
	}
}

func markNodeDiscoveredMinusOne(node *graph.TreeNode) {
	if node == nil {
		return
	}
	node.NodeDiscovered = -1
	for _, child := range node.Children {
		markNodeDiscoveredMinusOne(child)
	}
}

func copyRecipeMap(src map[string]Recipe) map[string]Recipe {
	dst := make(map[string]Recipe)
	for k, v := range src {
		dst[k] = v
	}
	return dst
}

func generateSignature(m map[string]Recipe) string {
	s := ""
	for k, v := range m {
		s += k + ":" + v.Source + "+" + v.Partner + "|"
	}
	return s
}