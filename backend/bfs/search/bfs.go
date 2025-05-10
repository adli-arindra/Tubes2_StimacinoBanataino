package search

import (
	"time"
	"bfs/graph"
)

// BFS untuk satu recipe
func BFS(target string, g graph.Graph) (graph.TreeResult, error) {
	start := time.Now()

	startingElements := []string{"Air", "Fire", "Water", "Earth"}
	discovered := map[string]*graph.TreeNode{}
	parentMap := map[string][]string{}
	queue := []string{}
	visitedNodes := 0

	// Inisialisasi dengan elemen dasar
	for _, el := range startingElements {
		discovered[el] = &graph.TreeNode{Name: el}
		queue = append(queue, el)
	}

	for len(queue) > 0 {
		curr := queue[0]
		queue = queue[1:]
		visitedNodes++

		// Mengecek semua elemen dalam Graph
		for product, recipes := range g {
			if _, alreadyFound := discovered[product]; alreadyFound {
				continue
			}
			for _, r := range recipes {
				if len(r) != 2 {
					continue
				}
				a, b := r[0], r[1]
				if (a == curr || b == curr) && discovered[a] != nil && discovered[b] != nil {
					discovered[product] = &graph.TreeNode{Name: product}
					parentMap[product] = []string{a, b}
					queue = append(queue, product)

					if product == target {
						duration := float64(time.Since(start).Microseconds()) / 1000.0
						tree := buildTree(target, discovered, parentMap)
						return graph.TreeResult{
							Tree:         tree,
							Algorithm:    "BFS",
							DurationMS:   duration,
							VisitedNodes: visitedNodes,
						}, nil
					}
					break
				}
			}
		}
	}

	// Kembalikan nil kalau misalnya elemen yang dicari ga ketemu
	duration := float64(time.Since(start).Microseconds()) / 1000.0
	return graph.TreeResult{
		Tree:         nil,
		Algorithm:    "BFS",
		DurationMS:   duration,
		VisitedNodes: visitedNodes,
	}, nil
}

func buildTree(current string, nodes map[string]*graph.TreeNode, parentMap map[string][]string) *graph.TreeNode {
	node := &graph.TreeNode{Name: current}
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