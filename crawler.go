// Package crawler implements a crawler that reads a list of URLs and
// crawls the URL periodically.
package main

import (
	"bufio"
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/jaeyeom/gofiletable/table"
)

var (
	tablePath = flag.String("table_path", "", "path to the backend table")
	urlsPath  = flag.String("urls_path", "", "path to the list of URLs")
)

type Fetcher interface {
	Fetch(url string) (body []byte, err error)
}

type DefaultFetcher struct {
}

// Create a new crawler and returns a channel that receives URLs.
func NewCrawler(tbl *table.Table, fetcher Fetcher, wg *sync.WaitGroup) chan string {
	c := make(chan string)
	go func() {
		if wg != nil {
			defer wg.Done()
		}
		for url := range c {
			body, err := fetcher.Fetch(url)
			if err != nil {
				log.Println("Crawl error:", url)
				// TODO: Handle error
				continue
			}
			tbl.Put([]byte(url), body)
		}
	}()
	return c
}

// Fetch fetches the page from URL.
func (f DefaultFetcher) Fetch(url string) (body []byte, err error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		// TODO: Handle the case.
		return []byte{}, nil
	}
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return
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
	fetcher := DefaultFetcher{}
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
	crawler := NewCrawler(tbl, fetcher, &wg)
	for urls.Scan() {
		crawler <- urls.Text()
	}
	close(crawler)
	wg.Wait()
}
