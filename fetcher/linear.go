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
		if *info.StatusCode != 200 {
			log.Println("Status is not 200:", url)
			continue
		}
		fetchInfo <- info
	}
}
