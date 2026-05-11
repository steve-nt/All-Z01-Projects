package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"golang.org/x/net/html"
)

type mirrorConfig struct {
	baseURL      *url.URL
	outputDir    string
	rejectTypes  []string
	excludeDirs  []string
	convertLinks bool
	rateLimit    int64
	visited      map[string]bool
	mu           sync.Mutex
}

// mirrorSite downloads an entire website recursively.
func mirrorSite(rawURL string, rejectTypes, excludeDirs []string, convertLinks bool, rateLimit int64, w io.Writer) error {
	baseURL, err := url.Parse(rawURL)
	if err != nil {
		return fmt.Errorf("invalid URL: %v", err)
	}

	outputDir := baseURL.Host
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output dir: %v", err)
	}

	mc := &mirrorConfig{
		baseURL:      baseURL,
		outputDir:    outputDir,
		rejectTypes:  rejectTypes,
		excludeDirs:  excludeDirs,
		convertLinks: convertLinks,
		rateLimit:    rateLimit,
		visited:      make(map[string]bool),
	}

	fmt.Fprintf(w, "Mirror: saving to directory '%s'\n", outputDir)
	mc.crawl(rawURL, w)
	fmt.Fprintf(w, "\nMirror complete.\n")
	return nil
}

func (mc *mirrorConfig) crawl(rawURL string, w io.Writer) {
	mc.mu.Lock()
	if mc.visited[rawURL] {
		mc.mu.Unlock()
		return
	}
	mc.visited[rawURL] = true
	mc.mu.Unlock()

	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return
	}

	// Only crawl same host
	if parsedURL.Host != mc.baseURL.Host {
		return
	}

	// Check excluded dirs
	for _, excl := range mc.excludeDirs {
		excl = strings.TrimSuffix(excl, "/")
		if strings.HasPrefix(parsedURL.Path, excl+"/") || parsedURL.Path == excl {
			return
		}
	}

	// Check rejected file types
	if mc.isRejected(parsedURL.Path) {
		return
	}

	// Fetch the page
	resp, err := http.Get(rawURL)
	if err != nil {
		fmt.Fprintf(w, "  failed to fetch %s: %v\n", rawURL, err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Fprintf(w, "  skip %s (status %s)\n", rawURL, resp.Status)
		return
	}

	contentType := resp.Header.Get("Content-Type")
	isHTML := strings.Contains(contentType, "text/html")
	isCSS := strings.Contains(contentType, "text/css")

	localPath := mc.urlToLocalPath(parsedURL)

	// Ensure directory exists
	dir := filepath.Dir(localPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		fmt.Fprintf(w, "  failed to create dir %s: %v\n", dir, err)
		return
	}

	fmt.Fprintf(w, "  saving: %s\n", localPath)

	var links []string

	if isHTML || isCSS {
		// Read body into memory so we can parse it
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return
		}

		if isHTML {
			links = extractHTMLLinks(body, parsedURL, mc.baseURL)
			if mc.convertLinks {
				body = convertHTMLLinks(body, parsedURL, mc.baseURL, mc.outputDir)
			}
		} else if isCSS {
			links = extractCSSLinks(body, parsedURL, mc.baseURL)
			if mc.convertLinks {
				body = convertCSSLinks(body, parsedURL, mc.baseURL, mc.outputDir)
			}
		}

		if err := os.WriteFile(localPath, body, 0644); err != nil {
			fmt.Fprintf(w, "  failed to write %s: %v\n", localPath, err)
		}
	} else {
		// Binary file: write directly
		f, err := os.Create(localPath)
		if err != nil {
			fmt.Fprintf(w, "  failed to create %s: %v\n", localPath, err)
			return
		}
		defer f.Close()

		var reader io.Reader = resp.Body
		if mc.rateLimit > 0 {
			reader = &rateLimitedReader{r: resp.Body, rateLimit: mc.rateLimit}
		}
		io.Copy(f, reader)
	}

	// Recursively crawl extracted links
	var wg sync.WaitGroup
	for _, link := range links {
		link := link
		mc.mu.Lock()
		already := mc.visited[link]
		mc.mu.Unlock()
		if already {
			continue
		}

		// Check if link should be skipped
		parsedLink, err := url.Parse(link)
		if err != nil {
			continue
		}
		if parsedLink.Host != mc.baseURL.Host {
			continue
		}
		if mc.isRejected(parsedLink.Path) {
			continue
		}
		skip := false
		for _, excl := range mc.excludeDirs {
			excl = strings.TrimSuffix(excl, "/")
			if strings.HasPrefix(parsedLink.Path, excl+"/") || parsedLink.Path == excl {
				skip = true
				break
			}
		}
		if skip {
			continue
		}

		wg.Add(1)
		go func(l string) {
			defer wg.Done()
			mc.crawl(l, w)
		}(link)
	}
	wg.Wait()
}

func (mc *mirrorConfig) isRejected(path string) bool {
	for _, ext := range mc.rejectTypes {
		ext = strings.TrimSpace(ext)
		if !strings.HasPrefix(ext, ".") {
			ext = "." + ext
		}
		if strings.HasSuffix(strings.ToLower(path), strings.ToLower(ext)) {
			return true
		}
	}
	return false
}

func (mc *mirrorConfig) urlToLocalPath(u *url.URL) string {
	path := u.Path
	if path == "" || path == "/" {
		path = "/index.html"
	} else if strings.HasSuffix(path, "/") {
		path = path + "index.html"
	} else if !strings.Contains(filepath.Base(path), ".") {
		path = path + "/index.html"
	}
	return filepath.Join(mc.outputDir, filepath.FromSlash(path))
}

func extractHTMLLinks(body []byte, pageURL, baseURL *url.URL) []string {
	var links []string
	doc, err := html.Parse(strings.NewReader(string(body)))
	if err != nil {
		return links
	}

	var walk func(*html.Node)
	walk = func(n *html.Node) {
		if n.Type == html.ElementNode {
			var attr string
			switch n.Data {
			case "a", "link":
				attr = "href"
			case "img", "script":
				attr = "src"
			}
			if attr != "" {
				for _, a := range n.Attr {
					if a.Key == attr && a.Val != "" {
						resolved := resolveURL(a.Val, pageURL)
						if resolved != "" && strings.HasPrefix(resolved, baseURL.Scheme+"://"+baseURL.Host) {
							links = append(links, resolved)
						}
					}
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			walk(c)
		}
	}
	walk(doc)
	return links
}

func extractCSSLinks(body []byte, pageURL, baseURL *url.URL) []string {
	var links []string
	content := string(body)
	// Find url(...) patterns in CSS
	for {
		idx := strings.Index(content, "url(")
		if idx == -1 {
			break
		}
		content = content[idx+4:]
		end := strings.Index(content, ")")
		if end == -1 {
			break
		}
		urlStr := strings.Trim(content[:end], `"' `)
		content = content[end+1:]

		resolved := resolveURL(urlStr, pageURL)
		if resolved != "" && strings.HasPrefix(resolved, baseURL.Scheme+"://"+baseURL.Host) {
			links = append(links, resolved)
		}
	}
	return links
}

func resolveURL(href string, base *url.URL) string {
	if href == "" || strings.HasPrefix(href, "#") || strings.HasPrefix(href, "javascript:") || strings.HasPrefix(href, "mailto:") {
		return ""
	}
	u, err := url.Parse(href)
	if err != nil {
		return ""
	}
	resolved := base.ResolveReference(u)
	resolved.Fragment = ""
	resolved.RawQuery = ""
	return resolved.String()
}

func convertHTMLLinks(body []byte, pageURL, baseURL *url.URL, outputDir string) []byte {
	content := string(body)
	doc, err := html.Parse(strings.NewReader(content))
	if err != nil {
		return body
	}

	var walk func(*html.Node)
	walk = func(n *html.Node) {
		if n.Type == html.ElementNode {
			var attrName string
			switch n.Data {
			case "a", "link":
				attrName = "href"
			case "img", "script":
				attrName = "src"
			}
			if attrName != "" {
				for i, a := range n.Attr {
					if a.Key == attrName && a.Val != "" {
						resolved := resolveURL(a.Val, pageURL)
						if resolved != "" && strings.HasPrefix(resolved, baseURL.Scheme+"://"+baseURL.Host) {
							u, err := url.Parse(resolved)
							if err == nil {
								localPath := localPathFromURL(u, outputDir)
								relPath, err := filepath.Rel(filepath.Dir(outputDir), localPath)
								if err == nil {
									n.Attr[i].Val = relPath
								}
							}
						}
					}
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			walk(c)
		}
	}
	walk(doc)

	var sb strings.Builder
	html.Render(&sb, doc)
	return []byte(sb.String())
}

func convertCSSLinks(body []byte, pageURL, baseURL *url.URL, outputDir string) []byte {
	content := string(body)
	var result strings.Builder
	for {
		idx := strings.Index(content, "url(")
		if idx == -1 {
			result.WriteString(content)
			break
		}
		result.WriteString(content[:idx+4])
		content = content[idx+4:]
		end := strings.Index(content, ")")
		if end == -1 {
			result.WriteString(content)
			break
		}
		quote := ""
		urlStr := content[:end]
		if len(urlStr) > 0 && (urlStr[0] == '"' || urlStr[0] == '\'') {
			quote = string(urlStr[0])
			urlStr = strings.Trim(urlStr, `"' `)
		}
		content = content[end:]

		resolved := resolveURL(urlStr, pageURL)
		if resolved != "" && strings.HasPrefix(resolved, baseURL.Scheme+"://"+baseURL.Host) {
			u, err := url.Parse(resolved)
			if err == nil {
				localPath := localPathFromURL(u, outputDir)
				relPath, err := filepath.Rel(filepath.Dir(outputDir), localPath)
				if err == nil {
					urlStr = relPath
				}
			}
		}
		result.WriteString(quote + urlStr + quote)
	}
	return []byte(result.String())
}

func localPathFromURL(u *url.URL, outputDir string) string {
	path := u.Path
	if path == "" || path == "/" {
		path = "/index.html"
	} else if strings.HasSuffix(path, "/") {
		path = path + "index.html"
	} else if !strings.Contains(filepath.Base(path), ".") {
		path = path + "/index.html"
	}
	return filepath.Join(outputDir, filepath.FromSlash(path))
}
