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
			recipes := make([][]string, 0)
			cols.Eq(1).Find("ul li").Each(func(_ int, li *goquery.Selection) {
				parts := li.Find("a[title]").Map(func(_ int, a *goquery.Selection) string {
				return a.Text()
				})

				if len(parts) > 0 {
					parts = mergeMultiWordElements(parts, catalog) // kasus untuk elemen dengan lebih dari satu kata
					recipes = append(recipes, parts)
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
			catalog.Tiers = append(catalog.Tiers, Tier{
				Name:     tierName,
				Elements: elements,
			})
		}
	})

	// Menghapus resep yang terdapat tier special
	catalog = removeSpecialRecipes(catalog)

	// Validasi resep agar tidak mengandung elemen dari myths dan monster
	validateRecipes(&catalog)

	// Menyimpan file JSON
	jsonFilePath := filepath.Join("data_scraping", "scraped_data.json")
	err = SaveToJSON(catalog, jsonFilePath)
	if err != nil {
		return Catalog{}, err
	}

	return catalog, nil
}

func cleanTierName(raw string) string {
	s := raw
	s = strings.TrimPrefix(s, "Tier ")
	s = strings.TrimSuffix(s, " elements")
	s = strings.TrimSuffix(s, " element")
	return s
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
		if tier.Name == "Special" {
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
			if tier.Name == "Special" {
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

func mergeMultiWordElements(parts []string, catalog Catalog) []string {
	var result []string
  
	for _, part := range parts {
		// Memeriksa apakah elemen terdiri lebih dari satu kata
		if isMultiWordElement(part, catalog) {
			result = append(result, part) 
		} else {
			result = append(result, part) 
		}
	}
  
	return result
}

func isMultiWordElement(element string, catalog Catalog) bool {
	for _, tier := range catalog.Tiers {
		for _, elem := range tier.Elements {
			if elem.Name == element && strings.Contains(element, " ") {
			return true
			}
		}
	}
	return false
}

func validateRecipes(catalog *Catalog) {
	for i := range catalog.Tiers {
		for j := range catalog.Tiers[i].Elements {
			validRecipes := make([][]string, 0)

			for _, recipe := range catalog.Tiers[i].Elements[j].Recipes {
				validRecipe := true
				for _, part := range recipe {
					if !isValidElement(part, *catalog) {
						validRecipe = false
						break
					}
				}

				if validRecipe {
					validRecipes = append(validRecipes, recipe)
				}
			}

			catalog.Tiers[i].Elements[j].Recipes = validRecipes
		}
	}
}

func isValidElement(element string, catalog Catalog) bool {
	for _, tier := range catalog.Tiers {
		for _, el := range tier.Elements {
			if el.Name == element {
			return true
			}
		}
	}
	return false 
}