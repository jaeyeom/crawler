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

// LogError reads errors from cerr and write logs.
func LogError(cerr <-chan error) {
	for err := range cerr {
		log.Println(err)
	}
}

// FeedFromFile reads a file from path and emits urls.
func FeedFromFile(path string, urls chan<- string, cerr chan<- error) {
	defer close(urls)
	f, err := os.Open(path)
	if err != nil {
		cerr <- err
		return
	}
	defer f.Close()
	urlScanner := bufio.NewScanner(f)
	for urlScanner.Scan() {
		urls <- urlScanner.Text()
	}
}

// WriteTable writes fetchInfo to the table.
func WriteTable(tbl *table.Table, fetchInfo <-chan *fetcher.FetchInfo, wg *sync.WaitGroup) {
	if wg != nil {
		defer wg.Done()
	}
	for info := range fetchInfo {
		tbl.Put([]byte(info.Key.String()), info.Contents)
	}
}

func main() {
	flag.Parse()
	cerr := make(chan error)
	defer close(cerr)
	go LogError(cerr)

	urls := make(chan string)
	go FeedFromFile(*urlsPath, urls, cerr)

	defaultFetcher := fetcher.DefaultFetcher{}
	fetchInfo := make(chan *fetcher.FetchInfo)
	go fetcher.ScheduleLinearFetch(defaultFetcher, urls, fetchInfo)

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
	WriteTable(tbl, fetchInfo, &wg)
	wg.Wait()
}
