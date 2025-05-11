package search

import (
	"dfs/graph"
	"time"

	"fmt"
)

func DFS(target string, g graph.Graph, elementTier map[string]int) (graph.TreeResult, error) {
	start := time.Now()

	// Starting element
	startingElements := []string{"Air", "Fire", "Water", "Earth"}
	discovered := map[string]*graph.TreeNode{}
	parentMap := map[string][]string{}
	visitedNodes := 0
	found := false

	// Jika target adalah starting element
	for _, base := range startingElements {
		if base == target {
			duration := float64(time.Since(start).Microseconds()) / 1000.0
			return graph.TreeResult{
				Tree: &graph.TreeNode{
					Name: base,
					Children: []*graph.TreeNode{},
				},
				Algorithm: "DFS",
				DurationMS: duration,
				VisitedNodes: 1,
			}, nil
		}
	}

	// Mengambil tier dari target
	maxTier, ok := elementTier[target]
	if !ok {
		return graph.TreeResult{
			Tree: nil, 
			Algorithm: "DFS",
			DurationMS: 0,
			VisitedNodes: 0,
		}, nil
	}

	var dfs func(current string)
	dfs = func(current string) {
		if found {
			return
		}
		fmt.Printf("DFS visiting: %s\n", current)
		visitedNodes++

		for product, recipes := range g {
			productTier, productOk := elementTier[product]
			if !productOk {
				continue
			}

			// Kalau elemen lebih tinggi dari target, skip ae
			if productTier > maxTier {
				continue
			}

			// Kalau elemen sudah pernah didapatkan, skip ae 
			if _, alreadyFound := discovered[product]; alreadyFound {
				continue
			}

			for _, r := range recipes {
				if len(r) != 2 {
					continue
				}
				a, b := r[0], r[1]
				if (a == current || b == current) && discovered[a] != nil && discovered[b] != nil {
					// Cek dulu elemen di recipe jangan sampe lebih tinggi dari elemen yang akan dihasilkan
					aTier, aOk := elementTier[a]
					bTier, bOk := elementTier[b]
					if !aOk || !bOk || aTier > productTier || bTier > productTier {
						continue
					}

					fmt.Printf("Membuat: %s dari %s + %s (tier: %d)\n", product, a, b, productTier) // debugging

					discovered[product] = &graph.TreeNode{Name: product}
					parentMap[product] = []string{a, b}

					if product == target {
						found = true
						return
					}

					dfs(product)
					if found {
						return
					}
				}
			}
		}
	}

	// Inisialisasi DFS dari starting element
	for _, element := range startingElements {
		discovered[element] = &graph.TreeNode{Name: element}
		dfs(element)
		if found {
			break
		}
	}

	duration := float64(time.Since(start).Microseconds()) / 1000.0

	if !found {
		return graph.TreeResult{
			Tree:         nil,
			Algorithm:    "DFS",
			DurationMS:   duration,
			VisitedNodes: visitedNodes,
		}, nil
	}

	tree := buildTree(target, discovered, parentMap)
	return graph.TreeResult{
		Tree:         tree,
		Algorithm:    "DFS",
		DurationMS:   duration,
		VisitedNodes: visitedNodes,
	}, nil
}

func buildTree(current string, nodes map[string]*graph.TreeNode, parentMap map[string][]string) *graph.TreeNode {
	node := &graph.TreeNode{
		Name: current,
		Children: []*graph.TreeNode{},
	}
	parents, ok := parentMap[current]
	
	if !ok {
		return node
	}

	for _, p := range parents {
		child := buildTree(p, nodes, parentMap)
		node.Children = append(node.Children, child)
	}

	return node
}
