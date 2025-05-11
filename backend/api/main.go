package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

type Node struct {
	Name           string  `json:"name"`
	NodeDiscovered int     `json:"node discovered"`
	Children       []Node  `json:"children"`
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
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		fmt.Println("Error decoding body:", err)
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	var filePath string

	switch req.Algorithm {
	case "BFS":
		if req.Mode == "single" {
			filePath = fmt.Sprintf("../bfs/result_BFS/%s_bfs.json", req.Target)
		} else {
			filePath = fmt.Sprintf("../bfs/result_multi_BFS/%s_multi_bfs_level.json", req.Target)
		}
	case "DFS":
		if req.Mode == "single" {
			filePath = fmt.Sprintf("../dfs/result_BFS/%s_bfs.json", req.Target)
		} else {
			filePath = fmt.Sprintf("../dfs/result_multi_BFS/%s_multi_bfs_level.json", req.Target)
		}
	default:
		http.Error(w, "Algoritma tidak ada", http.StatusBadRequest)
		return
	}

	jsonFile, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Error opening JSON file:", err)
		http.Error(w, "Error reading file", http.StatusInternalServerError)
		return
	}
	defer jsonFile.Close()

	var response SearchResponse
	err = json.NewDecoder(jsonFile).Decode(&response)
	if err != nil {
		fmt.Println("Error decoding JSON file:", err)
		http.Error(w, "Error parsing tree data", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}