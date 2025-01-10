package parsers

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"sync"

	"golang.org/x/net/html"
)

func FetchFavicon(resourceURL string) (string, error) {
	resp, err := http.Get(resourceURL)
	if err != nil {
		return "", fmt.Errorf("failed to fetch URL: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("received non-200 response code: %d", resp.StatusCode)
	}

	doc, err := html.Parse(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to parse HTML: %w", err)
	}

	var baseURL *url.URL
	baseURL, _ = url.Parse(resourceURL)

	var faviconURL string
	var found bool
	var mu sync.Mutex
	wg := &sync.WaitGroup{}

	var traverse func(*html.Node)
	traverse = func(n *html.Node) {
		defer wg.Done()
		if n.Type == html.ElementNode && n.Data == "link" {
			var rel, href string
			for _, attr := range n.Attr {
				if attr.Key == "rel" {
					rel = strings.ToLower(attr.Val)
				}
				if attr.Key == "href" {
					href = attr.Val
				}
			}
			if (rel == "icon" || rel == "shortcut icon") && href != "" {
				mu.Lock()
				if !found {
					faviconURL = href
					found = true
				}
				mu.Unlock()
				return
			}
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			if found {
				return
			}
			wg.Add(1)
			go traverse(c)
		}
	}

	wg.Add(1)
	go traverse(doc)
	wg.Wait()

	if faviconURL == "" {
		return "", fmt.Errorf("favicon not found")
	}

	faviconURLParsed, err := url.Parse(faviconURL)
	if err != nil {
		return "", fmt.Errorf("invalid favicon URL: %w", err)
	}
	if !faviconURLParsed.IsAbs() {
		faviconURL = baseURL.ResolveReference(faviconURLParsed).String()
	}

	return faviconURL, nil
}