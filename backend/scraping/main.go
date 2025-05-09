package main

import (
	"fmt"
	"log"

	"scraping/data_scraping"
)

func main() {
	data, err := data_scraping.ScrapeAll()
	if err != nil {
		log.Fatal("Gagal Scraping", err)
	}

	err = data_scraping.SaveToJSON(data, "data_scraping/scraped_data.json")
	if err != nil {
		log.Fatal("Gagal Save", err)
	}

	fmt.Println("Data berhasil di-scraping")
}
