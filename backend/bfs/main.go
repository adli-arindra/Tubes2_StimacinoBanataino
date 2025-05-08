package main

import (
	"fmt"
	"bfs/graph"
)

func main() {
	g, err := graph.LoadRecipes("../scraping/data_scraping/scraped_data.json")
	if err != nil {
		panic(err)
	}

	fmt.Println("Total elements with recipes:", len(g))
	for k, v := range g {
		fmt.Printf("%s ← %v\n", k, v)
		break
	}
}