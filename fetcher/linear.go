package fetcher

import (
	"log"
)

// ScheduleLinearFetch fetches each url in urls using the fetcherImpl and emits fetchInfo.
func ScheduleLinearFetch(fetcherImpl Fetcher, urls <-chan string, fetchInfo chan<- *FetchInfo) {
	defer close(fetchInfo)
	for url := range urls {
		info, err := fetcherImpl.Fetch(url)
		if err != nil {
			log.Println("Crawl error:", url)
			// TODO: Handle error
			continue
		}
		if info == nil {
			log.Println("Status is not 200:", url)
			continue
		}
		if info.Key == nil {
			log.Println("URL parsing failed:", url)
			continue
		}
		fetchInfo <- info
	}
}
