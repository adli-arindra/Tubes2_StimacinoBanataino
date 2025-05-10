package search

import (
	"dfs/graph"
	"time"
)

func DFS(target string, g graph.Graph) (graph.TreeResult, error) {
	start := time.Now()

	startingElements := []string{"Air", "Fire", "Water", "Earth"}
	discovered := map[string]*graph.TreeNode{}
	parentMap := map[string][]string{}
	visitedNodes := 0
	found := false

	var dfs func(current string)
	dfs = func(current string) {
		if found {
			return
		}
		visitedNodes++

		for product, recipes := range g {
			if _, alreadyFound := discovered[product]; alreadyFound {
				continue
			}
			for _, r := range recipes {
				if len(r) != 2 {
					continue
				}
				a, b := r[0], r[1]
				if (a == current || b == current) && discovered[a] != nil && discovered[b] != nil {
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
