package graph

import (
	"encoding/json"
	"os"
)

func LoadCatalog(path string) (ScrapedData, error) {
	// Membaca file sesuai dengan path
	data, err := os.ReadFile(path)
	if err != nil {
		return ScrapedData{}, err
	}
	var catalog ScrapedData
	// Mengubah data JSON menjadi struct dari ScrapedData
	err = json.Unmarshal(data, &catalog)
	return catalog, err
}

// Mengkonversikan ScrapedData menjadi Graph
func LoadRecipes(path string) (Graph, error) {
	catalog, err := LoadCatalog(path)
	if err != nil {
		return nil, err
	}
 
	graph := make(Graph)
	for _, tier := range catalog.Tiers {
		for _, el := range tier.Elements {
			for _, recipe := range el.Recipes {
				if len(recipe) == 2 {
					graph[el.Name] = append(graph[el.Name], recipe)
				}
			}
		}
	}
	return graph, nil
}

// Menyimpan map dari nama elemen ke indeks tier
func MapElementToTier(catalog ScrapedData) map[string]int {
	elementTier := make(map[string]int)
	for t, tier := range catalog.Tiers {
		for _, e := range tier.Elements {
			elementTier[e.Name] = t
		}
	}
	return elementTier
}