package main

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"golang.org/x/net/html"
	"golang.org/x/net/publicsuffix"
)

const (
	baseURL  string = "https://nl.pepper.com"
	maxDepth uint16 = 8
)

var (
	g = digraph{
		vertices: map[string]*vertex{},
		lock:     sync.RWMutex{},
	}
)

var (
	lock = sync.RWMutex{}
	wg   = sync.WaitGroup{}
)

var client = http.Client{
	Timeout: 20 * time.Second,
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

	fmt.Println("Operation took", time.Since(start))

	writeVertices(g.vertices)
}

func retrieve(u string, d uint16, p *vertex) {
	defer wg.Done()

	if d == maxDepth {
		return
	}

	response, err := client.Get(u)

	if err != nil {
		return
	}

	fmt.Println(u)

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
			temp := tokenizer.Token()

			isAnchor := temp.Data == "a"

			if isAnchor {
				for _, attribute := range temp.Attr {
					if attribute.Key == "href" {
						dom, err := parseBaseURL(attribute.Val)
						if err == nil {
							links[dom] = struct{}{}
						}
					}
				}
			}
		}
	}

	for k := range links {
		kTLD, _ := parseTLD(k)
		t := vertex{
			element: kTLD,
			in:      []*vertex{p},
			out:     []*vertex{},
			lock:    sync.RWMutex{},
		}

		if !g.containsDomain(kTLD) {
			g.addVertex(&t)
			p.addOutgoing(&t)

			wg.Add(1)
			go retrieve(k, d+1, &t)
		} else {
			uTLD, _ := parseTLD(u)
			g.getVertex(kTLD).addIncoming(g.getVertex(uTLD))
		}
	}
}

func parseBaseURL(u string) (string, error) {
	v, err := url.Parse(u)

	if err == nil && (v.Scheme == "https" || v.Scheme == "http") {
		return v.Scheme + "://" + v.Hostname(), nil
	}
	return "", fmt.Errorf("Domain could not be parsed")
}

func parseTLD(u string) (string, error) {
	v, err := url.Parse(u)
	if err == nil {
		// ignore non-http[s] schemes
		if v.Scheme == "https" || v.Scheme == "http" {
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
			return "", fmt.Errorf("Incorrect URL format")
		}
		return "", fmt.Errorf("Wrong scheme")
	}
	return "", fmt.Errorf("Could not parse domain")
}

func writeVertices(vertices map[string]*vertex) {
	sep := "\n"

	file, _ := os.OpenFile("urls.txt", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	defer file.Close()

	for _, v := range vertices {
		t := v.customString() + sep
		file.WriteString(t)
	}
}
