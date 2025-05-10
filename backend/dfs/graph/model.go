package graph

type Element struct {
	Name	string		`json:"name"`
	Recipes	[][]string	`json:"recipes"`
}

type Tier struct {
	Name	string		`json:"name"`
	Elements []Element	`json:"elements"`
}

type ScrapedData struct {
	Tiers	[]Tier		`json:"tiers"`
}

type Graph map[string][][]string

type TreeNode struct {
	Name		string			`json:"name"`
	Children	[]*TreeNode		`json:"children,omitempty"`
}

type TreeResult struct {
	Tree        	*TreeNode 	`json:"tree"`       
	Algorithm   	string    	`json:"algorithm"`  
	DurationMS  	float64   	`json:"duration_ms"` 
	VisitedNodes	int      	`json:"visited_nodes"` 
}