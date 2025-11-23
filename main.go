package main

import (
	"fmt"
	"go-concurrent-job-scapper/internal/config"
	"go-concurrent-job-scapper/internal/model"
	"go-concurrent-job-scapper/internal/scraper"
	"go-concurrent-job-scapper/internal/storage"
	"math/rand"
	"strings"
	"sync"
	"time"
)

func main() {
	fmt.Println(strings.Repeat("=", 70))
	fmt.Println("JOB SCRAPER - NEPAL EDITION")
	fmt.Println(strings.Repeat("=", 70))

	storage.SetupCSV(config.CSVFilename)

	var allJobs []model.Job
	var wg sync.WaitGroup
	jobChan := make(chan []model.Job, 1000)

	for _, q := range config.SearchQueries {
		// Merojob
		wg.Add(1)
		go func(query string) {
			defer wg.Done()
			jobChan <- scraper.ScrapeMerojob(query, config.MerojobPages)
		}(q)

		// Kumarijob
		wg.Add(1)
		go func(query string) {
			defer wg.Done()
			jobChan <- scraper.ScrapeKumarijob(query, config.KumarijobPages)
		}(q)
		// LinkedIN
		for _, loc := range config.NepalLocations {
			wg.Add(1)
			go func(query, l string) {
				defer wg.Done()
				jobChan <- scraper.ScrapeIndeedjob(query, l, config.IndeedPages)
			}(q, loc)
			linkedinSem := make(chan struct{}, 2)
			wg.Add(1)
			go func(query, l string) {
				defer wg.Done()
				linkedinSem <- struct{}{}
				defer func() {
					<-linkedinSem
				}()
				jobChan <- scraper.ScrapeLinkedInjob(query, l, config.LinkedInPages)
				time.Sleep(time.Duration(2+rand.Intn(4)) * time.Second)
			}(q, loc)
		}

	}
	// Close channel
	go func() {
		wg.Wait()
		close(jobChan)
	}()

	for batch := range jobChan {
		allJobs = append(allJobs, batch...)
	}

	existing := storage.GetExistingURLs(config.CSVFilename)
	newCount := storage.AppendJobs(config.CSVFilename, allJobs, existing)

	fmt.Println("Total Jobs:", len(allJobs))
	fmt.Println("New Jobs Added:", newCount)
}
