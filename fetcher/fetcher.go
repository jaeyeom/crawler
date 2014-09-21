// Package fetcher defines interface of fetchers and implements a
// default fetcher.
package fetcher

import (
	"io/ioutil"
	"net/http"
	"net/url"

	"code.google.com/p/goprotobuf/proto"
)

// Fetcher is an interface that implements Fetch function.
type Fetcher interface {
	Fetch(rawurl string) (fetchInfo *FetchInfo, err error)
}

// DefaultFetcher is a straightforward fetcher implementation that
// uses http.Get().
type DefaultFetcher struct {
}

// Fetch fetches the page from URL and fills fetchInfo.
func (f DefaultFetcher) Fetch(rawurl string) (fetchInfo *FetchInfo, err error) {
	info := &FetchInfo{}
	var key *url.URL
	key, err = url.Parse(rawurl)
	if err != nil {
		return
	}
	info.Url = proto.String(key.String())
	resp, err := http.Get(rawurl)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	info.StatusCode = proto.Int32(int32(resp.StatusCode))
	info.Contents, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	fetchInfo = info
	return
}
