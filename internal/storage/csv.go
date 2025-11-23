package storage

import (
	"encoding/csv"
	"go-concurrent-job-scapper/internal/model"
	"os"
	"strings"
)

func SetupCSV(filename string) {

	if _, err := os.Stat(filename); os.IsNotExist(err) {
		f, _ := os.Create(filename)
		defer f.Close()

		w := csv.NewWriter(f)
		w.Write([]string{"Job Title", "Company", "Location", "Job URL",
			"Source", "Date Posted", "Date Scraped", "Status", "Notes"})
		w.Flush()
	}

}

func GetExistingURLs(filename string) map[string]bool {
	m := make(map[string]bool)
	f, err := os.Open(filename)
	if err != nil {
		return m
	}
	defer f.Close()

	r := csv.NewReader(f)
	rows, _ := r.ReadAll()

	for i, row := range rows {
		if i == 0 {
			continue
		}
		if len(row) < 4 {
			continue
		}
		url := strings.TrimSpace(row[3])
		if url == "" {
			continue // ❗ skip empty URLs
		}
		m[url] = true
	}
	return m
}

func AppendJobs(filename string, jobs []model.Job, existing map[string]bool) int {
	f, _ := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, 0644)
	defer f.Close()

	w := csv.NewWriter(f)
	defer w.Flush()

	count := 0

	for _, job := range jobs {
		if existing[job.URL] {
			continue
		}
		w.Write([]string{
			job.Title, job.Company, job.Location, job.URL, job.Source,
			job.DatePosted, job.DateScraped, job.Status, job.Notes,
		})
		existing[job.URL] = true
		count++
	}

	return count
}
