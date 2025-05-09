package data_scraping

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

const baseURL = "https://little-alchemy.fandom.com/wiki/Elements_(Little_Alchemy_2)"

func ScrapeAll() (Catalog, error) {
	// Mengirim request untuk mendapatkan halaman web
	resp, err := http.Get(baseURL)
	if err != nil {
		return Catalog{}, err
	}
	defer resp.Body.Close()

	// Memparsing halaman HTML
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return Catalog{}, err
	}

	var catalog Catalog

	// Mendapatkan <h3> sebagai subjudul dari tier
	doc.Find("h3").Each(func(_ int, hdr *goquery.Selection) {
		rawTitle := hdr.Find("span.mw-headline").Text()
		if rawTitle == "" {
			return
		}

		// Mencari elemen-elemen dari tiap tier
		tbl := hdr.NextAll().Filter("table.list-table").First()
		if tbl.Length() == 0 {
			return
		}

		// Menyimpan elemen ke dalam folder
		tierName := cleanTierName(rawTitle)
		tierFolderName := strings.ReplaceAll(tierName, " ", "_")
		tierDir := filepath.Join("data_scraping", tierFolderName)

		os.MkdirAll(tierDir, 0755)

		elements := []Element{}
		tbl.Find("tr").Each(func(i int, row *goquery.Selection) {
			if i == 0 {
				return
			}

			cols := row.Find("td")
			if cols.Length() < 2 {
				return
			}

			// Mengambil nama elemen
			name := cols.Eq(0).Find("a[title]").First().Text()
			if name == "" {
				return
			}

			// Mengambil resep dari elemen
			recipes := [][]string{}
			cols.Eq(1).Find("li").Each(func(i int, li *goquery.Selection){
				recipe := li.Text()
				if recipe != "" {
					recipe = strings.ReplaceAll(recipe, "+", "") // Hapus tanda + antara 2 elemen resep
					recipes = append(recipes, strings.Fields(recipe))
				}
			})

			// Mengambil link untuk file SVG
			fileA := cols.Eq(0).Find("a.mw-file-description")
			href, _ := fileA.Attr("href")
			local := ""

			// Menyimpan file SVG
			if href != "" {
				fname := strings.ReplaceAll(name, " ", "_") + ".svg"
				local = filepath.Join(tierDir, fname)

				err := downloadSVG(href, local)
				if err != nil {
					return
				}
			}

			elements = append(elements, Element{
				Name:           name,
				LocalSVGPath:   strings.ReplaceAll(local, "\\", "/"),
				OriginalSVGURL: href,
				Recipes:        recipes,
			})
		})

		if len(elements) > 0 {
			// Menambahkan tier ke dalam katalog
			catalog.Tiers = append(catalog.Tiers, Tier{
				Name:     tierName,
				Elements: elements,
			})
		}
	})

	// Menghapus resep yang mengandung elemen dari tier special
	catalog = removeSpecialRecipes(catalog)

	// Menyimpan file JSON
	jsonFilePath := filepath.Join("data_scraping", "scraped_data.json")
	err = SaveToJSON(catalog, jsonFilePath)
	if err != nil {
		return Catalog{}, err
	}

	return catalog, nil
}

func cleanTierName(rawTitle string) string {
	return strings.TrimSpace(rawTitle)
}

func downloadSVG(url, localPath string) error {
	// Mengirimkan request untuk mengunduh SVG
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Membuat file baru untuk menyimpan hasil unduhan
	file, err := os.Create(localPath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Membaca dan menulis konten dari body HTTP ke dalam file
	_, err = file.ReadFrom(resp.Body)
	if err != nil {
		return err
	}

	return nil
}

func SaveToJSON(data interface{}, filePath string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Mengonversi data ke format JSON dan menyimpannya dalam file
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	err = encoder.Encode(data)
	if err != nil {
		return err
	}

	return nil
}

func removeSpecialRecipes(catalog Catalog) Catalog {
	var filteredTiers []Tier
	for _, tier := range catalog.Tiers {
		if tier.Name == "Special element" {
			continue
		}
		var filteredElements []Element
		for _, element := range tier.Elements {
			var filteredRecipes [][]string
			for _, recipe := range element.Recipes {
				if !containsSpecialElement(recipe, catalog) {
					filteredRecipes = append(filteredRecipes, recipe)
				}
			}

			element.Recipes = filteredRecipes
			filteredElements = append(filteredElements, element)
		}

		tier.Elements = filteredElements
		filteredTiers = append(filteredTiers, tier)
	}

	catalog.Tiers = filteredTiers
	return catalog
}

func containsSpecialElement(recipe []string, catalog Catalog) bool {
	for _, ingredient := range recipe {
		for _, tier := range catalog.Tiers {
			// Jika nama tier adalah "Special" dan ada elemen yang cocok, return true
			if tier.Name == "Special element" {
				for _, element := range tier.Elements {
					if element.Name == ingredient {
						return true
					}
				}
			}
		}
	}
	return false
}