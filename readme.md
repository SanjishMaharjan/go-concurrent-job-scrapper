# Go Concurrent Job Scraper

This is a powerful, concurrent web scraper written in Go, designed to gather job listings from popular job boards in Nepal. It efficiently scrapes multiple sites simultaneously, aggregates the data, and stores it in a structured CSV file, avoiding duplicate entries.

## Features

- **Concurrent Scraping**: Utilizes Go's concurrency features (goroutines and channels) to scrape multiple job sites and search queries simultaneously for maximum speed.
- **Multiple Job Boards**: Scrapes from several popular job platforms:
  - Merojob.com
  - Kumarijob.com
  - LinkedIn
  - Indeed
- **Configurable Search**: Easily customize search queries, locations, and the number of pages to scrape in the configuration file.
- **Duplicate Prevention**: Checks for existing job URLs in the output file to ensure only new listings are added.
- **CSV Output**: Saves all scraped job listings into a clean, easy-to-use CSV file named `nepal_job_listings.csv`.

## Prerequisites

- [Go](https://golang.org/doc/install) version 1.24 or higher.

## Installation & Usage

1.  **Clone the repository:**

    ```bash
    git clone https://github.com/sanjishmaharjan/go-concurrent-job-scrapper.git
    cd go-concurrent-job-scrapper
    ```

2.  **Install dependencies:**
    The project uses Go Modules. Dependencies will be downloaded automatically when you build or run the project.

3.  **Run the scraper:**
    Execute the main application from your terminal.

    ```bash
    go run main.go
    ```

    The scraper will start, displaying its progress in the console.

    ```
    ======================================================================
    JOB SCRAPER - NEPAL EDITION
    ======================================================================

    Scraping Merojob for 'React Developer'...

    🔍 Scraping Kumarijob.com for 'React Developer'...

    🔍 Scraping Indeed for 'React Developer' in Kathmandu...
    ...
    ✅ LinkedIn: Scraped 125 jobs for 'Software Engineer' in Nepal
    Total Jobs: 1530
    New Jobs Added: 95
    ```

## Configuration

You can customize the scraper's behavior by modifying the constants in `internal/config/config.go`.

```go
package config

const (
	CSVFilename    = "nepal_job_listings.csv"
	MerojobPages   = 5 // Number of pages to scrape from Merojob
	KumarijobPages = 5 // Number of pages to scrape from Kumarijob
	IndeedPages    = 1
	LinkedInPages  = 5
)

// SearchQueries defines the job titles to search for.
var SearchQueries = []string{
	"React Developer",
	"Full Stack Developer",
	"Go Lang Developer",
	"DevOps Engineer",
	// Add more queries here
}

// NepalLocations defines the locations to search within.
var NepalLocations = []string{
	"Kathmandu", "Lalitpur", "Bhaktapur", "Pokhara", "Nepal",
}
```

## Output

The scraper generates a CSV file named `nepal_job_listings.csv` in the root directory. This file contains the aggregated job listings with the following columns:

- `Job Title`
- `Company`
- `Location`
- `Job URL`
- `Source` (e.g., Merojob.com, LinkedIn)
- `Date Posted`
- `Date Scraped`
- `Status` (Defaults to "Not Applied")
- `Notes` (Empty by default)

### Example `nepal_job_listings.csv`:

```csv
Job Title,Company,Location,Job URL,Source,Date Posted,Date Scraped,Status,Notes
Front-End Developer (Junior),Softbenz Infosys,"Kathmandu, Bāgmatī, Nepal",https://np.linkedin.com/jobs/view/front-end-developer-junior-at-softbenz-infosys-4340211861,LinkedIn,2025-11-12,2025-11-23,Not Applied,
Senior Full Stack Developer,WeFlow Agency,"Kathmandu, Bāgmatī, Nepal",https://np.linkedin.com/jobs/view/senior-full-stack-developer-at-weflow-agency-4323073922,LinkedIn,2025-11-18,2025-11-23,Not Applied,
DevOps Engineer,UXCam,"Kathmandu, Bāgmatī, Nepal",https://np.linkedin.com/jobs/view/devops-engineer-at-uxcam-4341131170,LinkedIn,2025-11-17,2025-11-23,Not Applied,
```

## Project Structure

```
.
├── go.mod
├── go.sum
├── main.go                 # Main application entry point
├── nepal_job_listings.csv  # Output file with scraped data
└── internal/
    ├── config/
    │   └── config.go       # Search queries, locations, and other settings
    ├── model/
    │   └── job.go          # Defines the Job data structure
    ├── scraper/
    │   ├── fetch.go        # Handles HTTP requests and HTML parsing
    │   ├── indeedscraper.go
    │   ├── kumariscraper.go
    │   ├── linkedinscraper.go
    │   └── meroscraper.go  # Scraping logic for each job board
    └── storage/
        └── csv.go          # Functions for reading and writing to the CSV file
```
