package scraper

import (
	"fmt"
	"net/http"
	"time"

	"github.com/PuerkitoBio/goquery"
)

func FetchDoc(url string) (*goquery.Document, error) {
	client := &http.Client{Timeout: 15 * time.Second}

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp == nil {
		return nil, fmt.Errorf("nil response")
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("status code %d", resp.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}
	if doc == nil {
		return nil, fmt.Errorf("parsed document is nil")
	}

	return doc, nil
}
