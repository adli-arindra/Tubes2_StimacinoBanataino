package data_scraping

// 
type Element struct {
	Name           string     `json:"name"`
	Recipes        [][]string `json:"recipes"`
	LocalSVGPath   string     `json:"local_svg_path"`
	OriginalSVGURL string     `json:"original_svg_url"`
}

// 
type Tier struct {
	Name     string    `json:"name"`
	Elements []Element `json:"elements"`
}

// 
type Catalog struct {
	Tiers []Tier `json:"tiers"`
}
