package scraper

import (
	"fmt"
	"go-concurrent-job-scapper/internal/config"
	"go-concurrent-job-scapper/internal/model"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

func ScrapeLinkedInjob(query, location string, pages int) []model.Job {
	var jobs []model.Job
	fmt.Printf("\n🔍 Scraping LinkedIn for '%s' in %s...\n", query, location)

	for p := 0; p < pages; p++ {
		url := fmt.Sprintf("https://www.linkedin.com/jobs-guest/jobs/api/seeMoreJobPostings/search?keywords=%s&location=%s&start=%d&sortBy=DD",
			strings.ReplaceAll(query, " ", "%20"),
			strings.ReplaceAll(location, " ", "%20"),
			p*25)

		doc, err := FetchDoc(url)
		if err != nil {
			fmt.Printf("  ❌ Error fetching LinkedIn page %d: %v\n", p+1, err)
			time.Sleep(config.RequestDelay)
			continue
		}

		doc.Find("li").Each(func(i int, s *goquery.Selection) {
			titleSel := s.Find("h3.base-search-card__title").First()
			if titleSel.Length() == 0 {
				return
			}
			title := strings.TrimSpace(titleSel.Text())
			company := strings.TrimSpace(s.Find("h4.base-search-card__subtitle").Text())
			if company == "" {
				company = "N/A"
			}
			jobLocation := strings.TrimSpace(s.Find("span.job-search-card__location").Text())
			if jobLocation == "" {
				jobLocation = location
			}
			href, exists := s.Find("a.base-card__full-link").Attr("href")
			if !exists {
				href = "N/A"
			}
			datePosted := "N/A"
			if timeSel := s.Find("time").First(); timeSel.Length() > 0 {
				if dt, ok := timeSel.Attr("datetime"); ok {
					datePosted = dt
				}
			}

			jobs = append(jobs, model.Job{
				Title:       title,
				Company:     company,
				Location:    jobLocation,
				URL:         href,
				Source:      "LinkedIn",
				DatePosted:  datePosted,
				DateScraped: time.Now().Format("2006-01-02"),
				Status:      "Not Applied",
				Notes:       "",
			})
		})

		fmt.Printf("  Page %d: Found %d job cards\n", p+1, len(jobs))
		time.Sleep(3 * time.Second)
	}

	fmt.Printf("✅ LinkedIn: Scraped %d jobs for '%s' in %s\n", len(jobs), query, location)
	return jobs
}
