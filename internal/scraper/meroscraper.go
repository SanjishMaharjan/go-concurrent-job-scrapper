package scraper

import (
	"fmt"
	"go-concurrent-job-scapper/internal/config"
	"go-concurrent-job-scapper/internal/model"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

func ScrapeMerojob(query string, pages int) []model.Job {
	var jobs []model.Job
	baseURL := "https://merojob.com"

	fmt.Printf("\nScraping Merojob for '%s'...\n", query)

	for p := 1; p <= pages; p++ {
		url := fmt.Sprintf("%s/search/?q=%s&page=%d", baseURL, strings.ReplaceAll(query, " ", "+"), p)

		doc, err := FetchDoc(url)
		if err != nil {
			fmt.Printf("Error fetching page %d: %v \n", p, err)
			time.Sleep(config.RequestDelay)
			continue
		}
		doc.Find("div.card-body").Each(func(_ int, s *goquery.Selection) {
			titleSel := s.Find("h1.job-title,a.text-dark").First()
			title := strings.TrimSpace(titleSel.Text())
			if title == "" {
				return
			}
			href, _ := titleSel.Attr("href")
			if !strings.HasPrefix(href, "http") {
				href = baseURL + href
			}

			jobs = append(jobs, model.Job{
				Title:       title,
				Company:     strings.TrimSpace(s.Find("p.company-name").Text()),
				Location:    strings.TrimSpace(s.Find("span.text-muted").Text()),
				URL:         href,
				Source:      "Merojob.com",
				DatePosted:  "N/A",
				DateScraped: time.Now().Format("2006-01-02"),
				Status:      "Not Applied",
			})
		})
		time.Sleep(config.RequestDelay)
	}
	return jobs

}
