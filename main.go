package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/mohamidsaiid/gophercises/sitemap-builder/link-parser"
	"github.com/mohamidsaiid/gophercises/sitemap-builder/queue"
)

type empty struct{}

// XML sitemap structs

type UrlSet struct {
	XMLName xml.Name `xml:"urlset"`
	Xmlns   string   `xml:"xmlns,attr"`
	Urls    []Url    `xml:"url"`
}

type Url struct {
	Loc string `xml:"loc"`
}

func main() {
	urlFlag := flag.String("url", "http://localhost:3000/", "specify which website you want to get its site map")
	flag.Parse()

	q := queue.Queue{
		Len:  1,
		Data: []string{*urlFlag},
	}
	visited := map[string]empty{}

	visited[q.Front()] = empty{}
	for !q.Empty() {
		current := q.Pop()
		fmt.Println(current)
		neighpors := get(current)
		for _, next := range neighpors {
			if _, ok := visited[next]; ok {
				continue
			}
			q.Push(next)
			visited[next] = empty{}
		}
	}
	// Collect URLs for XML
	var urlList []Url
	for lnk := range visited {
		urlList = append(urlList, Url{Loc: lnk})
	}
	urlSet := UrlSet{
		Xmlns: "http://www.sitemaps.org/schemas/sitemap/0.9",
		Urls:  urlList,
	}
	file, err := os.Create("sitemap.xml")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	enc := xml.NewEncoder(file)
	enc.Indent("", "  ")
	file.WriteString(xml.Header)
	if err := enc.Encode(urlSet); err != nil {
		panic(err)
	}
	fmt.Println("Sitemap written to sitemap.xml")
}

func get(siteUrl string) []string {
	resp, err := http.Get(siteUrl)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	reqUrl := resp.Request.URL
	baseLink := &url.URL{
		Scheme: reqUrl.Scheme,
		Host:   reqUrl.Host,
	}

	res, err := link.Parse(resp.Body)
	if err != nil {
		return nil
	}
	return filter(baseLink, res)
}

func filter(baseLnk *url.URL, links []link.Link) []string {
	l := []string{}
	for _, val := range links {
		l = append(l, val.Href)
	}

	res := []string{}
	for _, val := range l {
		if strings.HasPrefix(val, baseLnk.String()) && baseLnk.String() != val {
			res = append(res, val)
		} else if strings.HasPrefix(val, "/") && baseLnk.String() != val {
			val = baseLnk.String() + val
			res = append(res, val)
		}
	}
	return res
}
/*func filter(baseLnk *url.URL, links []link.Link) []string {
	res := []string{}
	seen := map[string]struct{}{}

	for _, val := range links {
		parsed, err := url.Parse(val.Href)
		if err != nil {
			continue
		}
		absURL := baseLnk.ResolveReference(parsed)
		// Only allow same host, no query, and http(s) scheme
		if absURL.Host != baseLnk.Host {
			continue
		}
		if absURL.Scheme != "http" && absURL.Scheme != "https" {
			continue
		}
		if absURL.RawQuery != "" {
			continue
		}
		// Remove fragment
		absURL.Fragment = ""
		// Avoid duplicates and self-link
		urlStr := absURL.String()
		if urlStr == baseLnk.String() {
			continue
		}
		if _, ok := seen[urlStr]; !ok {
			res = append(res, urlStr)
			seen[urlStr] = struct{}{}
		}
	}
	return res
}*/