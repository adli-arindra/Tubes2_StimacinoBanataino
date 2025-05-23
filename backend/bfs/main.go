package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"bfs/graph"
	"bfs/search"
)

func main1() {
	const recipeFile = "../scraping/data_scraping/scraped_data.json"

    // Input
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Masukkan nama elemen target: ")
	rawInput, err := reader.ReadString('\n')
	if err != nil {
		log.Fatalf("Gagal membaca input: %v", err) // debugging
	}
	target := strings.TrimSpace(rawInput)

	catalog, err := graph.LoadCatalog(recipeFile)
	if err != nil {
		log.Fatalf("Gagal load catalog: %v", err)
	}

	elementTier := graph.MapElementToTier(catalog)

	// Load scraped_data
	g, err := graph.LoadRecipes(recipeFile)
	if err != nil {
		log.Fatalf("Gagal load graph: %v", err) // debugging
	}

	// Process
	result, err := search.BFS(target, g, elementTier)
	if err != nil {
		log.Fatalf("Gagal menjalankan BFS: %v", err) // debugging
	}

	// Save
	output, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		log.Fatalf("Gagal encode hasil ke JSON: %v", err) // debugging
	}

	err = os.MkdirAll("result_BFS", os.ModePerm)
	if err != nil {
		log.Fatalf("Gagal membuat folder result: %v", err) // debugging
	}

	fileSafeName := strings.ReplaceAll(strings.ToLower(target), " ", "_")
	filePath := fmt.Sprintf("result_BFS/%s_bfs.json", fileSafeName)

	err = os.WriteFile(filePath, output, 0644)
	if err != nil {
		log.Fatalf("Gagal menulis file hasil: %v", err) // debugging
	}

	fmt.Printf("Berhasil menyimpan hasil\n")
}