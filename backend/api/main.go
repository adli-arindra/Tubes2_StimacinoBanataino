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

	bidirectionalGraph "bidirectional/graph"
	bidirectionalSearch "bidirectional/search"

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

type MultipleSearchResponse struct {
	Tree         []Node  `json:"trees"`
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
	case "Bidirectional":
		{
			directory = "../bidirectional/result_bidirectional"
			filePath = fmt.Sprintf("%s/%s_bidirectional.json", directory, target)
		}
	default:
		http.Error(w, "Algoritma tidak ada", http.StatusBadRequest)
		return
	}

	graphBFS, err := bfsGraph.LoadRecipes("../scraping/data_scraping/scraped_data.json")
	graphDFS, err := dfsGraph.LoadRecipes("../scraping/data_scraping/scraped_data.json")
	graphBidirectional, err := bidirectionalGraph.LoadRecipes("../scraping/data_scraping/scraped_data.json")
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
			multiResult, err := bfsSearch.MultiBFS(req.Target, graphBFS, *req.MaxRecipes, elementTiers)
			if err != nil {
				http.Error(w, "Pencarian gagal: "+err.Error(), http.StatusInternalServerError)
				return
			}

			var nodes []Node
			for _, t := range multiResult.Trees {
				nodes = append(nodes, convertTreeNodeToNodeBFS(t))
			}

			result = MultipleSearchResponse{
				Tree:         nodes,
				Algorithm:    multiResult.Algorithm,
				Duration:     multiResult.DurationMS,
				VisitedNodes: multiResult.VisitedNodes,
			}
		}
	case "DFS":
		if req.Mode == "single" {
			result, err = dfsSearch.DFS(req.Target, graphDFS, elementTiers)
		} else {
			multiResult, err := dfsSearch.MultiDFS(req.Target, graphDFS, *req.MaxRecipes, elementTiers)
			if err != nil {
				http.Error(w, "Pencarian gagal: "+err.Error(), http.StatusInternalServerError)
				return
			}

			var nodes []Node
			for _, t := range multiResult.Trees {
				nodes = append(nodes, convertTreeNodeToNodeDFS(t))
			}

			result = MultipleSearchResponse{
				Tree:         nodes,
				Algorithm:    multiResult.Algorithm,
				Duration:     multiResult.DurationMS,
				VisitedNodes: multiResult.VisitedNodes,
			}
		}
	case "Bidirectional":
		result, err = bidirectionalSearch.Bidirectional(req.Target, graphBidirectional, elementTiers)
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

	jsonFile, err := os.Open(filePath)
	if err != nil {
		http.Error(w, "Gagal membuka file hasil", http.StatusInternalServerError)
		return
	}
	defer jsonFile.Close()

	if req.Mode == "single" {
		var response SearchResponse
		if err := json.NewDecoder(jsonFile).Decode(&response); err != nil {
			http.Error(w, "Gagal decode hasil pencarian", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	} else {
		var response MultipleSearchResponse
		if err := json.NewDecoder(jsonFile).Decode(&response); err != nil {
			http.Error(w, "Gagal decode hasil pencarian", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}

// Converter for multiple recipe
func convertTreeNodeToNodeBFS(t *bfsGraph.TreeNode) Node {
	if t == nil {
		return Node{}
	}

	children := make([]Node, len(t.Children))
	for i, child := range t.Children {
		children[i] = convertTreeNodeToNodeBFS(child)
	}

	return Node{
		Name:           t.Name,
		NodeDiscovered: t.NodeDiscovered,
		Children:       children,
	}
}

func convertTreeNodeToNodeDFS(t *dfsGraph.TreeNode) Node {
	if t == nil {
		return Node{}
	}

	children := make([]Node, len(t.Children))
	for i, child := range t.Children {
		children[i] = convertTreeNodeToNodeDFS(child)
	}

	return Node{
		Name:           t.Name,
		NodeDiscovered: t.NodeDiscovered,
		Children:       children,
	}
}
