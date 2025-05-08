package search

import (
	"fmt"
	"bfs/graph"
)

// Menyimpan informasi dari item ke queue BFS
type QueueItem struct {
	Node	string
	Paths	[]string
	Recipes [][]string
}

func BFS(target string, g graph.Graph) (map[string][][]string, error) {
	// Mulai dari 4 elemen dasar
	queue := []QueueItem{}
	startingElements := []string{"Air", "Fire", "Earth", "Water"}

	for _, elem := range startingElements {
		queue = append(queue, QueueItem{Node: elem, Paths: []string{elem}, Recipes: [][]string{}})
	}

	visited := make(map[string]bool) // Melacak elemen yang sudah dikunjungi

	result := make(map[string][][]string)

	for len(queue) > 0 {
		item := queue[0]
		queue = queue[1:]

		// Skip elemen yang uda dikunjungi
		if visited[item.Node] {
			continue
		}

		visited[item.Node] = true

		if item.Node == target {
			result[target] = item.Recipes
		}

		for _, recipe := range g[item.Node] {
			for _, next := range recipe {
				if !visited[next] {
					// Menambah elemen berikutnya ke antrian dan resep 
					newRecipes := append(item.Recipes, []string{fmt.Sprintf("%s , %s", recipe[0], recipe[1])})
					queue = append(queue, QueueItem{
						Node:    next,
						Paths:   append(item.Paths, next),
						Recipes: newRecipes,
					})
				}
			}
		}
	}

	// Jika tidak ditemukan jalur, kembalikan error
	if len(result) == 0 {
		return nil, fmt.Errorf("target %s not found", target)
	}

	return result, nil
}