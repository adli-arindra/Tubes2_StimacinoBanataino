package main

import (
	"fmt"
	"log"
	"scraping/data_scraping"
)

func main() {
	catalog, err := data_scraping.ScrapeAll()
	if err != nil {
		log.Fatalf("Error during scraping: %v", err)
	}

	fmt.Println("Scraping selesai!")
	fmt.Printf("Total Tiers: %d\n", len(catalog.Tiers))

	for _, tier := range catalog.Tiers {
		fmt.Printf("Tier: %s\n", tier.Name)
		fmt.Printf("Jumlah Elemen: %d\n", len(tier.Elements))
		for _, element := range tier.Elements {
			fmt.Printf("  Elemen: %s\n", element.Name)
			fmt.Printf("    SVG Path: %s\n", element.LocalSVGPath)
			fmt.Printf("    Resep: %v\n", element.Recipes)
			fmt.Printf("    Original SVG URL: %s\n", element.OriginalSVGURL)
		}
	}

	fmt.Println("Proses scraping selesai dan data berhasil disimpan!")
}