package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"dfs/graph"
	"dfs/search"
)

func main() {
	const recipeFile = "../scraping/data_scraping/scraped_data.json"

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Masukkan nama elemen target: ")
	rawInput, err := reader.ReadString('\n')
	if err != nil {
		log.Fatalf("Gagal membaca input: %v", err)
	}
	target := strings.TrimSpace(rawInput)

	g, err := graph.LoadRecipes(recipeFile)
	if err != nil {
		log.Fatalf("Gagal load graph: %v", err)
	}

	result, err := search.DFS(target, g)
	if err != nil {
		log.Fatalf("Gagal menjalankan DFS: %v", err)
	}

	output, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		log.Fatalf("Gagal encode hasil ke JSON: %v", err)
	}

	err = os.MkdirAll("result", os.ModePerm)
	if err != nil {
		log.Fatalf("Gagal membuat folder result: %v", err)
	}

	fileSafeName := strings.ReplaceAll(strings.ToLower(target), " ", "_")
	filePath := fmt.Sprintf("result/%s_dfs.json", fileSafeName)

	err = os.WriteFile(filePath, output, 0644)
	if err != nil {
		log.Fatalf("Gagal menulis file hasil: %v", err)
	}

	fmt.Printf("Berhasil menyimpan hasil ke %s\n", filePath)
}