package config

import "time"

const (
	CSVFilename    = "nepal_job_listings.csv"
	MerojobPages   = 5
	KumarijobPages = 5
	IndeedPages    = 1
	LinkedInPages  = 5
	RequestDelay   = 2 * time.Second
)

var SearchQueries = []string{
	"React Developer",
	"Full Stack Developer",
	"Go Lang Developer",
	"DevOps Engineer",
	"Backend Developer",
	"Frontend Developer",
	"Software Engineer",
}

var NepalLocations = []string{
	"Kathmandu", "Lalitpur", "Bhaktapur", "Pokhara", "Nepal",
}
