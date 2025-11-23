package scraper

import (
	"fmt"
	"go-concurrent-job-scapper/internal/config"
	"go-concurrent-job-scapper/internal/model"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

func ScrapeKumarijob(query string, pages int) []model.Job {
	var jobs []model.Job
	baseURL := "https://kumarijob.com"
	fmt.Printf("\n🔍 Scraping Kumarijob.com for '%s'...\n", query)

	for p := 1; p <= pages; p++ {
		url := fmt.Sprintf("%s/search?keywords=%s&page=%d", baseURL, strings.ReplaceAll(query, " ", "+"), p)
		doc, err := FetchDoc(url)
		if err != nil {
			fmt.Printf("  ❌ Error fetching Kumarijob page %d: %v\n", p, err)
			time.Sleep(config.RequestDelay)
			continue
		}

		selection := doc.Find("div.job-list-item, div.job-item")
		selection.Each(func(i int, s *goquery.Selection) {
			titleSel := s.Find("h2, h3, a.job-title").First()
			if titleSel.Length() == 0 {
				return
			}
			title := strings.TrimSpace(titleSel.Text())
			var href string
			if titleSel.Is("a") {
				href, _ = titleSel.Attr("href")
			} else {
				a := titleSel.Find("a").First()
				href, _ = a.Attr("href")
			}
			if href == "" {
				href = "N/A"
			}
			if !strings.HasPrefix(href, "http") && href != "N/A" {
				href = baseURL + href
			}

			company := strings.TrimSpace(s.Find("span.company, div.company-name").First().Text())
			if company == "" {
				company = "N/A"
			}
			location := strings.TrimSpace(s.Find("span.location, i.location").First().Text())
			if location == "" {
				location = "Nepal"
			}
			datePosted := strings.TrimSpace(s.Find("span.date, time").First().Text())
			if datePosted == "" {
				datePosted = "N/A"
			}

			if title != "" && href != "N/A" {
				jobs = append(jobs, model.Job{
					Title:       title,
					Company:     company,
					Location:    location,
					URL:         href,
					Source:      "Kumarijob.com",
					DatePosted:  datePosted,
					DateScraped: time.Now().Format("2006-01-02"),
					Status:      "Not Applied",
					Notes:       "",
				})
			}
		})

		fmt.Printf("  Page %d: Found %d job cards\n", p, selection.Length())
		time.Sleep(config.RequestDelay)
	}

	fmt.Printf("✅ Kumarijob: Scraped %d jobs for '%s'\n", len(jobs), query)
	return jobs

}
