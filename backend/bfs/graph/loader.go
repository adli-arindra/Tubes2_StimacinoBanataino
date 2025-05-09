package graph

import (
	"encoding/json"
	"os"
)

func LoadRecipes(path string) (Graph, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var raw ScrapedData
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, err
	}

	graph := make(Graph)
	for _, tier := range raw.Tiers {
		for _, el := range tier.Element {
			for _, recipe := range el.Recipes {
				if len(recipe) == 2 {
					graph[el.Name] = append(graph[el.Name], recipe)
				}
			}
		}
	}
	return graph, nil
}