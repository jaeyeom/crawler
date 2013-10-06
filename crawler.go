// Package crawler implements a crawler that reads a list of URLs and
// crawls the URL periodically.
package main

import (
	"bufio"
	"flag"
	"log"
	"os"
	"sync"

	"github.com/jaeyeom/crawler/fetcher"
	"github.com/jaeyeom/gofiletable/table"
)

var (
	tablePath = flag.String("table_path", "", "path to the backend table")
	urlsPath  = flag.String("urls_path", "", "path to the list of URLs")
)

// Create a new crawler and returns a channel that receives URLs.
func NewCrawler(tbl *table.Table, fetcherImpl fetcher.Fetcher, wg *sync.WaitGroup) chan string {
	c := make(chan string)
	go func() {
		if wg != nil {
			defer wg.Done()
		}
		for url := range c {
			fetchInfo, err := fetcherImpl.Fetch(url)
			if err != nil {
				log.Println("Crawl error:", url)
				// TODO: Handle error
				continue
			}
			if fetchInfo == nil {
				log.Println("Status is not 200:", url)
			}
			if fetchInfo.Key == nil {
				log.Println("URL parsing failed:", url)
			}
			tbl.Put([]byte(fetchInfo.Key.String()), fetchInfo.Contents)
		}
	}()
	return c
}

func main() {
	flag.Parse()
	f, err := os.Open(*urlsPath)
	if err != nil {
		log.Println(err)
		return
	}
	defer f.Close()
	urls := bufio.NewScanner(f)
	defaultFetcher := fetcher.DefaultFetcher{}
	tbl, err := table.Create(table.TableOption{
		BaseDirectory: *tablePath,
		KeepSnapshots: true,
	})
	if err != nil {
		log.Println(err)
		return
	}
	var wg sync.WaitGroup
	wg.Add(1)
	crawler := NewCrawler(tbl, defaultFetcher, &wg)
	for urls.Scan() {
		crawler <- urls.Text()
	}
	close(crawler)
	wg.Wait()
}
