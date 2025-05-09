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
	resp, err := http.Get(baseURL)
	if err != nil {
		return Catalog{}, err
	}
	defer resp.Body.Close()
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return Catalog{}, err
	}

	var catalog Catalog
	doc.Find("h3").Each(func(_ int, hdr *goquery.Selection) {
		rawTitle := hdr.Find("span.mw-headline").Text()
		if rawTitle == "" {
			return
		}

		tbl := hdr.NextAll().Filter("table.list-table").First()
		if tbl.Length() == 0 {
			return
		}

		tierName := cleanTierName(rawTitle)
		tierFolderName := strings.ReplaceAll(tierName, " ", "_")
		tierDir := filepath.Join("data_scraping", tierFolderName)

		elems := []Element{}
		if len(elems) > 0 {
			os.MkdirAll(tierDir, 0755)
			catalog.Tiers = append(catalog.Tiers, Tier{
				Name:     tierName,
				Elements: elems,
			})
		}

		tbl.Find("tr").Each(func(i int, row *goquery.Selection) {
			if i == 0 {
				return
			}
			cols := row.Find("td")
			if cols.Length() < 2 {
				return
			}

			name := cols.Eq(0).Find("a[title]").First().Text()
			if name == "" {
				return
			}

			fileA := cols.Eq(0).Find("a.mw-file-description")
			href, _ := fileA.Attr("href")
			local := ""
			if href != "" {
				fname := strings.ReplaceAll(name, " ", "_") + ".svg"
				local = filepath.Join(tierFolderName, fname)
			}

			recipes := [][]string{}
			cols.Eq(1).Find("ul li").Each(func(_ int, li *goquery.Selection) {
				parts := li.Find("a[title]").Map(func(_ int, a *goquery.Selection) string {
					return a.Text()
				})
				if len(parts) == 2 {
					recipes = append(recipes, []string{parts[0], parts[1]})
				}
			})

			elems = append(elems, Element{
				Name:           name,
				LocalSVGPath:   local,
				OriginalSVGURL: href,
				Recipes:        recipes,
			})
		})

		if len(elems) > 0 {
			catalog.Tiers = append(catalog.Tiers, Tier{
				Name:     tierName,
				Elements: elems,
			})
		} else {
			os.RemoveAll(tierDir)
		}
	})

	return catalog, nil
}

func SaveToJSON(catalog Catalog, filename string) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	return enc.Encode(catalog)
}

func cleanTierName(raw string) string {
	s := strings.TrimPrefix(raw, "Tier ")
	s = strings.TrimSuffix(s, " elements")
	s = strings.TrimSuffix(s, " element")
	return s
}