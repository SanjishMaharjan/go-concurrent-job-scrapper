package scraper

import (
	"fmt"
	"go-concurrent-job-scapper/internal/config"
	"go-concurrent-job-scapper/internal/model"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

func ScrapeIndeedjob(query, location string, pages int) []model.Job {
	var jobs []model.Job
	baseURL := "https://www.indeed.com/jobs"
	fmt.Printf("\n🔍 Scraping Indeed for '%s' in %s...\n", query, location)

	for p := 0; p < pages; p++ {
		url := fmt.Sprintf("%s?q=%s&l=%s&start=%d&sort=date",
			baseURL,
			strings.ReplaceAll(query, " ", "+"),
			strings.ReplaceAll(location, " ", "+"),
			p*10)

		doc, err := FetchDoc(url)
		if err != nil || doc == nil {
			fmt.Printf("  ❌ Error fetching Indeed page %d: %v\n", p+1, err)
			time.Sleep(config.RequestDelay)
			continue
		}

		cardSel := doc.Find("div.job_seen_beacon")
		if cardSel.Length() == 0 {
			cardSel = doc.Find("td.resultContent")
		}

		cardSel.Each(func(i int, s *goquery.Selection) {
			title := strings.TrimSpace(s.Find("h2.jobTitle, h2").Text())
			if title == "" {
				return
			}
			company := strings.TrimSpace(s.Find("span.companyName").Text())
			if company == "" {
				company = "N/A"
			}
			jobLocation := strings.TrimSpace(s.Find("div.companyLocation").Text())
			if jobLocation == "" {
				jobLocation = location
			}
			href, exists := s.Find("a").Attr("href")
			if !exists || href == "" {
				return
			}
			if !strings.HasPrefix(href, "http") {
				href = "https://www.indeed.com" + href
			}
			datePosted := strings.TrimSpace(s.Find("span.date").Text())
			if datePosted == "" {
				datePosted = "N/A"
			}

			jobs = append(jobs, model.Job{
				Title:       title,
				Company:     company,
				Location:    jobLocation,
				URL:         href,
				Source:      "Indeed",
				DatePosted:  datePosted,
				DateScraped: time.Now().Format("2006-01-02"),
				Status:      "Not Applied",
				Notes:       "",
			})
		})

		fmt.Printf("  Page %d: Found %d job cards\n", p+1, cardSel.Length())
		time.Sleep(config.RequestDelay)
	}

	fmt.Printf("✅ Indeed: Scraped %d jobs for '%s' in %s\n", len(jobs), query, location)
	return jobs
}
