package search

import (
	"time"
	"bfs/graph"
)

// BFS untuk satu recipe
func BFS(target string, g graph.Graph, elementTier map[string]int) (graph.TreeResult, error) {
	start := time.Now()

	startingElements := []string{"Air", "Fire", "Water", "Earth"}
	discovered := map[string]*graph.TreeNode{}
	parentMap := map[string][]string{}
	queue := []string{}
	visitedNodes := 0

	// Inisialisasi dengan elemen dasar
	for _, el := range startingElements {
		discovered[el] = &graph.TreeNode{
			Name: el, 
			Children: []*graph.TreeNode{},
		}
		queue = append(queue, el)
	}

	// Cek target starting element atau bukan
	if node, ok := discovered[target]; ok {
		duration := float64(time.Since(start).Microseconds()) / 1000.0
		return graph.TreeResult{
			Tree:         node,
			Algorithm:    "BFS",
			DurationMS:   duration,
			VisitedNodes: 1,
		}, nil
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

					// Buat pengecekan recipe harus dari elemen yang lebih rendah
					productTier, productOk := elementTier[product]
					aTier, aOk := elementTier[a]
					bTier, bOk := elementTier[b]
					if !productOk || !aOk || !bOk || aTier > productTier || bTier > productTier {
						continue
					}

					discovered[product] = &graph.TreeNode{Name: product}
					parentMap[product] = []string{a, b}
					queue = append(queue, product)

					if product == target {
						duration := float64(time.Since(start).Microseconds()) / 1000.0
						tree := buildTree(target, discovered, parentMap) // Urutan buat fitur live update
						index := 0
						setDiscoveredIndex(tree, &index)
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
		node.Children = []*graph.TreeNode{}
		return node
	}
	for _, p := range parents {
		child := buildTree(p, nodes, parentMap)
		node.Children = append(node.Children, child)
	}
	return node
}

// Node discovered dimulai dari root ke leaf
func setDiscoveredIndex(root *graph.TreeNode, counter *int) {
	if root == nil {
		return
	}

	queue := []*graph.TreeNode{root}

	for len(queue) > 0 {
		curr := queue[0]
		queue = queue[1:]

		curr.NodeDiscovered = *counter
		*counter++

		queue = append(queue, curr.Children...)
	}
}