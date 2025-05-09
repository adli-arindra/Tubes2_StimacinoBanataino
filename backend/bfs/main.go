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
        return
    }

    // Format output sesuai scraped_data.json
    fmt.Printf("\n{\n  \"name\": \"%s\",\n  \"recipes\": [\n", target)
    
    recipes := result[target]
    for i, recipe := range recipes {
        fmt.Printf("    [\n      \"%s\",\n      \"%s\"\n    ]", recipe[0], recipe[1])
        if i < len(recipes)-1 {
            fmt.Printf(",")
        }
        fmt.Println()
    }
    fmt.Println("  ]\n}")
}