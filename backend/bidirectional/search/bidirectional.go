package search

import (
	"bidirectional/graph"
	"time"
)

func Bidirectional(target string, g graph.Graph, tierMap map[string]int) (graph.TreeResult, error) {
	start := time.Now()

	// Starting element
	startingElements := []string{"Air", "Fire", "Water", "Earth"}

	// Kalau target starting element
	if isStartingElement(target, startingElements) {
		return graph.TreeResult{
			Tree:         &graph.TreeNode{
				Name: target, 
				NodeDiscovered: 0,
				Children: []*graph.TreeNode{},
			},
			Algorithm:    "Bidirectional",
			DurationMS:   float64(time.Since(start).Microseconds()) / 1000.0,
			VisitedNodes: 1,
		}, nil
	}

	// Menyimpan map node yang telah dikunjungi
	forwardVisited := map[string]struct{}{}
	backwardVisited := map[string]struct{}{target: {}}

	// Menyimpan parent node untuk membentuk pohon
	forwardParent := map[string][]string{}
	backwardParent := map[string][]string{}

	// Menyimpan node tier untuk tiap arah
	forwardTier := map[int][]string{}
	backwardTier := map[int][]string{}

	// Inisiasi forward search
	for _, el := range startingElements {
		forwardVisited[el] = struct{}{}
		if t, ok := tierMap[el]; ok {
			forwardTier[t] = append(forwardTier[t], el)
		}
	}

	// Inisiasi backward search
	if t, ok := tierMap[target]; ok {
		backwardTier[t] = append(backwardTier[t], target)
	}

	visitedNodes := 0
	maxTier := -1
	for _, t := range tierMap {
		if t > maxTier {
			maxTier = t
		}
	}

	var meetingPoint string // Titik temu
	found := false

	for depth := 0; depth <= maxTier && !found; depth++ {
		// Forward search
		for _, curr := range forwardTier[depth] {
			_ = curr
			for product, recipes := range g {
				if _, seen := forwardVisited[product]; seen {
					continue
				}
				for _, r := range recipes {
					if len(r) != 2 {
						continue
					}
					a, b := r[0], r[1]
					_, aOk := forwardVisited[a]
					_, bOk := forwardVisited[b]
					// Skip ae kalo ketemu element yang belum dikunjungi
					if !aOk || !bOk {
						continue
					}
					// Cek tier
					productTier, okP := tierMap[product]
					aTier, okA := tierMap[a]
					bTier, okB := tierMap[b]
					if !okP || !okA || !okB || aTier >= productTier || bTier >= productTier {
						continue
					}
					forwardVisited[product] = struct{}{}
					forwardParent[product] = []string{a, b}
					forwardTier[productTier] = append(forwardTier[productTier], product)
					visitedNodes++

					// Mengecek udah ditemui di backward atau belum
					if _, ok := backwardVisited[product]; ok {
						meetingPoint = product
						found = true
						break
					}
				}
				if found {
					break
				}
			}
			if found {
				break
			}
		}

		// Backward search
		for _, curr := range backwardTier[maxTier-depth] {
			_ = curr
			for product, recipes := range g {
				if _, ok := tierMap[product]; !ok || product != curr {
					continue
				}
				for _, r := range recipes {
					if len(r) != 2 {
						continue
					}
					a, b := r[0], r[1]
					for _, ing := range []string{a, b} {
						if _, seen := backwardVisited[ing]; seen {
							continue
						}
						backwardVisited[ing] = struct{}{}
						backwardParent[product] = []string{a, b}
						if t, ok := tierMap[ing]; ok {
							backwardTier[t] = append(backwardTier[t], ing)
						}
						visitedNodes++

						// Mengecek udah ditemui di forward atau belum
						if _, ok := forwardVisited[ing]; ok {
							meetingPoint = ing
							found = true
							break
						}
					}
					if found {
						break
					}
				}
				if found {
					break
				}
			}
			if found {
				break
			}
		}
	}

	if found {
		// Gabung parent dan build tree dari titik temu
		merged := mergeMaps(forwardParent, backwardParent)
		tree := buildTree(meetingPoint, merged)
		zeroNodeDiscovered(tree) // Karena ga pake live update jadinya di set 0 aja
		return graph.TreeResult{
			Tree:         tree,
			Algorithm:    "Bidirectional",
			DurationMS:   float64(time.Since(start).Microseconds()) / 1000.0,
			VisitedNodes: visitedNodes,
		}, nil
	}

	// Kalau target ga ditemukan
	return graph.TreeResult{
		Tree:         nil,
		Algorithm:    "Bidirectional",
		DurationMS:   float64(time.Since(start).Microseconds()) / 1000.0,
		VisitedNodes: visitedNodes,
	}, nil
}

func isStartingElement(name string, starters []string) bool {
	for _, el := range starters {
		if name == el {
			return true
		}
	}
	return false
}

// Gabungin dua parent
func mergeMaps(m1, m2 map[string][]string) map[string][]string {
	merged := map[string][]string{}
	for k, v := range m1 {
		merged[k] = v
	}
	for k, v := range m2 {
		if _, exists := merged[k]; !exists {
			merged[k] = v
		}
	}
	return merged
}

func buildTree(current string, parentMap map[string][]string) *graph.TreeNode {
	node := &graph.TreeNode{Name: current}
	parents, ok := parentMap[current]
	if !ok {
		node.Children = []*graph.TreeNode{}
		return node
	}
	for _, p := range parents {
		child := buildTree(p, parentMap)
		node.Children = append(node.Children, child)
	}
	return node
}

func zeroNodeDiscovered(root *graph.TreeNode) {
	if root == nil {
		return
	}
	queue := []*graph.TreeNode{root}
	for len(queue) > 0 {
		curr := queue[0]
		queue = queue[1:]
		curr.NodeDiscovered = 0
		queue = append(queue, curr.Children...)
	}
}
