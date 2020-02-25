package main

import (
	"fmt"
	"golang.org/x/net/publicsuffix"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	"golang.org/x/net/html"
)

const (
	baseURL  string = "https://nl.pepper.com"
	maxDepth uint16 = 20
)

var wg   = sync.WaitGroup{}

var g = digraph{
	vertices: map[string]*vertex{},
	lock:     sync.RWMutex{},
}

var client = http.Client{
	Timeout: 30 * time.Second,
}

func main() {
	var depth uint16

	start := time.Now()

	baseTLD, _ := parseTLD(baseURL)

	t := vertex{
		element: baseTLD,
		in:      []*vertex{},
		out:     []*vertex{},
		lock:    sync.RWMutex{},
	}

	g.addVertex(&t)

	wg.Add(1)
	go retrieve(baseURL, depth, &t)
	wg.Wait()

	fmt.Println("operation took", time.Since(start))

	g.serialize("graph.txt")
}

func retrieve(src string, d uint16, p *vertex) {
	defer wg.Done()

	if d == maxDepth {
		return
	}

	response, err := client.Get(src)

	if err != nil {
		return
	}

	fmt.Println(src)

	tokenizer := html.NewTokenizer(response.Body)
	defer response.Body.Close()

	links := map[string]struct{}{}

Loop:
	for {
		token := tokenizer.Next()

		switch token {
		case html.ErrorToken:
			break Loop
		case html.StartTagToken:
			t := tokenizer.Token()

			isAnchor := t.Data == "a"

			if isAnchor {
				for _, attribute := range t.Attr {
					if attribute.Key == "href" {
						base, err := parseBaseURL(attribute.Val)
						if err == nil {
							links[base] = struct{}{}
						}
					}
				}
			}
		}
	}

	for link := range links {
		destTLD, _ := parseTLD(link)
		t := vertex{
			element: destTLD,
			in:      []*vertex{p},
			out:     []*vertex{},
			lock:    sync.RWMutex{},
		}

		if !g.containsDomain(destTLD) {
			g.addVertex(&t)
			p.addOutgoing(&t)

			wg.Add(1)
			go retrieve(link, d+1, &t)
		} else {
			srcTLD, _ := parseTLD(src)
			g.getVertex(destTLD).addIncoming(g.getVertex(srcTLD))
		}
	}
}

func parseBaseURL(u string) (string, error) {
	v, err := url.Parse(u)

	if err == nil && (v.Scheme == "https" || v.Scheme == "http") {
		return v.Scheme + "://" + v.Hostname(), nil
	}
	return "", fmt.Errorf("domain could not be parsed")
}

func parseTLD(u string) (string, error) {
	v, err := url.Parse(u)
	if err == nil {
		// ignore non-http[s] schemes
		if v.Scheme == "http" || v.Scheme == "https" {
			h := v.Hostname()
			sep := "."

			// check if hostname is ip address
			p := strings.Join(strings.Split(h, sep), "")

			if _, err := strconv.Atoi(p); err == nil {
				return h, err
			}

			extension, _ := publicsuffix.PublicSuffix(h)
			i := strings.LastIndex(h, extension)
			if i > 1 { // else invalid url hostname format due to parsing
				r := strings.Split(h[:i-1], sep)   // get domain components
				t := r[len(r)-1] + sep + extension // get root domain + extension
				return t, err
			}
			return "", fmt.Errorf("incorrect URL format")
		}
		return "", fmt.Errorf("wrong scheme")
	}
	return "", fmt.Errorf("could not parse domain")
}
