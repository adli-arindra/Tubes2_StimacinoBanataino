package graph

type Element struct {
	Name	string		`json:"name"`
	Recipes	[][]string	`json:"recipes"`
}

type Tier struct {
	Name	string		`json:"name"`
	Element []Element	`json:"elements"`
}

type ScrapedData struct {
	Tiers	[]Tier		`json:"tiers"`
}

type Graph map[string][][]string

type Node struct {
	Value		string
	Children	[]*Node
}