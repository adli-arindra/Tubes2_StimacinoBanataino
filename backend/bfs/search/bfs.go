package search

import (
    "fmt"
    "bfs/graph"
)

type QueueItem struct {
    Node    string
    Paths   []string
    Recipes [][]string
}

func BFS(target string, g graph.Graph) (map[string][][]string, error) {
    queue := []QueueItem{}
    startingElements := []string{"Air", "Fire", "Earth", "Water"}
    visited := make(map[string]bool)
    result := make(map[string][][]string)

    // Inisialisasi elemen dasar
    for _, elem := range startingElements {
        visited[elem] = true
        queue = append(queue, QueueItem{
            Node:    elem,
            Paths:   []string{elem},
            Recipes: [][]string{},
        })
    }

    for len(queue) > 0 {
        item := queue[0]
        queue = queue[1:]

        if item.Node == target {
            result[target] = item.Recipes
            continue
        }

        if currentRecipes, exists := g[item.Node]; exists {
            for _, recipe := range currentRecipes {
                if len(recipe) != 2 {
                    continue
                }

                resultElement := recipe[1] // Selemen yag dibuat
                if visited[resultElement] {
                    continue
                }

                if !visited[recipe[0]] {
                    continue
                }

                visited[resultElement] = true
                fullRecipe := []string{item.Node, recipe[0], resultElement}
                newRecipes := append(item.Recipes, fullRecipe)
                
                queue = append(queue, QueueItem{
                    Node:    resultElement,
                    Paths:   append(item.Paths, resultElement),
                    Recipes: newRecipes,
                })
            }
        }

        for resultElement, recipes := range g {
            if visited[resultElement] {
                continue
            }

            for _, recipe := range recipes {
                if len(recipe) != 2 || !visited[recipe[0]] {
                    continue
                }

                if recipe[1] == item.Node {
                    visited[resultElement] = true
                    fullRecipe := []string{recipe[0], item.Node, resultElement}
                    newRecipes := append(item.Recipes, fullRecipe)
                    
                    queue = append(queue, QueueItem{
                        Node:    resultElement,
                        Paths:   append(item.Paths, resultElement),
                        Recipes: newRecipes,
                    })
                }
            }
        }
    }

    if len(result) == 0 {
        return nil, fmt.Errorf("target %s not found", target)
    }

    return result, nil
}