package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	bfsGraph "bfs/graph"
	bfsSearch "bfs/search"

	dfsGraph "dfs/graph"
	dfsSearch "dfs/search"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

type Node struct {
	Name           string `json:"name"`
	NodeDiscovered int    `json:"node_discovered"`
	Children       []Node `json:"children"`
}

type SearchRequest struct {
	Target     string `json:"target"`
	Algorithm  string `json:"algorithm"`
	Mode       string `json:"mode"`
	MaxRecipes *int   `json:"max_recipes,omitempty"`
}

type SearchResponse struct {
	Tree         Node    `json:"tree"`
	Algorithm    string  `json:"algorithm"`
	Duration     float64 `json:"duration_ms"`
	VisitedNodes int     `json:"visited_nodes"`
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/search", processSearchRequest).Methods("POST")

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type"},
		AllowCredentials: true,
	})

	handler := c.Handler(r)

	fmt.Println("server listening on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", handler))
}

func processSearchRequest(w http.ResponseWriter, r *http.Request) {
	var req SearchRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	target := strings.ReplaceAll(strings.ToLower(strings.TrimSpace(req.Target)), " ", "_")

	var filePath string
	var directory string

	switch req.Algorithm {
	case "BFS":
		if req.Mode == "single" {
			directory = "../bfs/result_BFS"
			filePath = fmt.Sprintf("%s/%s_bfs.json", directory, target)
		} else {
			directory = "../bfs/result_multi_BFS"
			filePath = fmt.Sprintf("%s/%s_multi_bfs_level.json", directory, target)
		}
	case "DFS":
		if req.Mode == "single" {
			directory = "../dfs/result_DFS"
			filePath = fmt.Sprintf("%s/%s_dfs.json", directory, target)
		} else {
			directory = "../dfs/result_multi_DFS"
			filePath = fmt.Sprintf("%s/%s_multi_dfs_level.json", directory, target)
		}
	default:
		http.Error(w, "Algoritma tidak ada", http.StatusBadRequest)
		return
	}

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		fmt.Println("File belum ada, file akan dibuat untuk mencari target")

		graphBFS, err := bfsGraph.LoadRecipes("../scraping/data_scraping/scraped_data.json")
		graphDFS, err := dfsGraph.LoadRecipes("../scraping/data_scraping/scraped_data.json")
		catalog, err := dfsGraph.LoadCatalog("../scraping/data_scraping/scraped_data.json")
		elementTiers := dfsGraph.MapElementToTier(catalog)

		if err != nil {
			http.Error(w, "Gagal load graph: "+err.Error(), http.StatusInternalServerError)
			return
		}

		var result interface{}
		switch req.Algorithm {
		case "BFS":
			if req.Mode == "single" {
				result, err = bfsSearch.BFS(req.Target, graphBFS, elementTiers)
			} else {
				// result, err = bfsSearch.MultiBFS(req.Target, graphDFS, *req.MaxRecipes,elementTiers)
			}
		case "DFS":
			if req.Mode == "single" {
				result, err = dfsSearch.DFS(req.Target, graphDFS, elementTiers)
			} else {
				// result, err = dfsSearch.MultiDFS(req.Target, graphDFS, *req.MaxRecipes,elementTiers)
			}
		}
		if err != nil {
			http.Error(w, "Pencarian gagal: "+err.Error(), http.StatusInternalServerError)
			return
		}

		output, err := json.MarshalIndent(result, "", "  ")
		if err != nil {
			http.Error(w, "Gagal encode hasil", http.StatusInternalServerError)
			return
		}

		if err := os.MkdirAll(directory, os.ModePerm); err != nil {
			http.Error(w, "Gagal membuat folder", http.StatusInternalServerError)
			return
		}

		if err := os.WriteFile(filePath, output, 0644); err != nil {
			http.Error(w, "Gagal menulis file hasil", http.StatusInternalServerError)
			return
		}
	}

	jsonFile, err := os.Open(filePath)
	if err != nil {
		http.Error(w, "Gagal membuka file hasil", http.StatusInternalServerError)
		return
	}
	defer jsonFile.Close()

	var response SearchResponse
	if err := json.NewDecoder(jsonFile).Decode(&response); err != nil {
		http.Error(w, "Gagal decode hasil pencarian", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
