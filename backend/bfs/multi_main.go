package main

import (
	"bufio"
	"bfs/graph"
	"bfs/search"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func main() {
	const recipeFile = "../scraping/data_scraping/scraped_data.json"

	// Input Target
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Masukkan nama elemen target: ")
	rawTarget, _ := reader.ReadString('\n')
	target := strings.TrimSpace(rawTarget)

	// Input Jumlah Recipe yang diinginkan
	fmt.Print("Masukkan jumlah maksimal recipe: ")
	rawCount, _ := reader.ReadString('\n')
	countStr := strings.TrimSpace(rawCount)
	maxRecipes, err := strconv.Atoi(countStr)
	if err != nil || maxRecipes <= 0 {
		log.Fatalf("Jumlah recipe tidak valid: %v", err) // debugging
	}

	// Load scraped_data.JSON
	catalog, err := graph.LoadCatalog(recipeFile)
	if err != nil {
		log.Fatalf("Gagal load catalog: %v", err) // debugging
	}

	// Bangun tier map
	tierMap := make(map[string]int)
	for tierIndex, tier := range catalog.Tiers {
		for _, el := range tier.Elements {
			tierMap[el.Name] = tierIndex
		}
	}

	// Load graph
	g, err := graph.LoadRecipes(recipeFile)
	if err != nil {
		log.Fatalf("Gagal load graph: %v", err) // debugging
	}

	// Process Multiple BFS
	result, err := search.MultiBFS(target, g, maxRecipes, tierMap)
	if err != nil {
		log.Fatalf("Gagal menjalankan multiple BFS: %v", err) // debugging
	}

	// Save Output
	output, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		log.Fatalf("Gagal menyimpan hasil: %v", err) // debugging
	}

	_ = os.MkdirAll("result_multi_BFS", os.ModePerm)
	filename := strings.ReplaceAll(strings.ToLower(target), " ", "_")
	err = os.WriteFile(fmt.Sprintf("result_multi_BFS/%s_multi_bfs_level.json", filename), output, 0644)
	if err != nil {
		log.Fatalf("Gagal menyimpan hasil ke file: %v", err) // debugging
	}

	fmt.Printf("\nBerhasil menyimpan hasil\n")
}