package graph

// Menyimpan elemen hasil scraping yang berisi nama dan resep
type Element struct {
	Name	string		`json:"name"`
	Recipes	[][]string	`json:"recipes"`
}

// Menyimpan tier yang berisi nama tier dan element di setiap tier
type Tier struct {
	Name	string		`json:"name"`
	Elements []Element	`json:"elements"`
}

// Root dari JSON
type ScrapedData struct {
	Tiers	[]Tier		`json:"tiers"`
}

// Graph untuk resep
type Graph map[string][][]string

// Menyimpan satu node dalam pohon
type TreeNode struct {
	Name			string			`json:"name"`
	NodeDiscovered 	int           	`json:"node_discovered"`
	Children		[]*TreeNode		`json:"children"`
}

// Menyimpan struktur output untuk Bidirectional One Recipe
type TreeResult struct {
	Tree        	*TreeNode 	`json:"tree"`       
	Algorithm   	string    	`json:"algorithm"`  
	DurationMS  	float64   	`json:"duration_ms"` 
	VisitedNodes	int      	`json:"visited_nodes"` 
}