package main

import (
	"fmt"
	"bfs/graph"
	"bfs/search"
)

func main() {
	g, err := graph.LoadRecipes("../scraping/data_scraping/scraped_data.json")
	if err != nil {
		panic(err)
	}

	var target string
	fmt.Println("Masukkan elemen target yang ingin dicari:")
	fmt.Scan(&target)

	result, err := search.BFS(target, g)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		// Menampilkan hasil
		fmt.Printf("%s: %v\n", target, result[target])
	}
}