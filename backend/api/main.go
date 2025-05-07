package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

type Node struct {
	Name  string `json:"name"`
	Idx   int    `json:"idx"`
	Left  *Node  `json:"left,omitempty"`
	Right *Node  `json:"right,omitempty"`
}

type SearchRequest struct {
	Target     string `json:"target"`
	Algorithm  string `json:"algorithm"`
	Mode       string `json:"mode"`
	MaxRecipes int    `json:"max_recipes,omitempty"`
}

type SearchResponse struct {
	Trees         []Node  `json:"trees"`
	NumberOfPaths int     `json:"numberOfPaths"`
	NodesVisited  int     `json:"nodesVisited"`
	ElapsedTime   float64 `json:"elapsedTime"`
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

	dummyTree := Node{
		Name: "Brick",
		Idx:  0,
		Left: &Node{
			Name: "Mud",
			Idx:  1,
			Left: &Node{
				Name: "Water",
				Idx:  2,
			},
			Right: &Node{
				Name: "Earth",
				Idx:  3,
			},
		},
		Right: &Node{
			Name: "Fire",
			Idx:  4,
		},
	}

	response := SearchResponse{
		Trees:         []Node{dummyTree},
		NumberOfPaths: 1,
		NodesVisited:  11,
		ElapsedTime:   0.23,
	}

	treeJSON, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(string(treeJSON))
	}
}
