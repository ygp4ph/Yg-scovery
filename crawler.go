package main

import (
	"bufio"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"runtime"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/fatih/color"
)

// Config holds configuration parameters for the crawler.
type Config struct {
	TargetURL    string
	MaxDepth     int
	OnlyInternal bool
	OnlyExternal bool
	OutputPath   string
	Verbose      bool
	ShowTree     bool
}

// Crawler represents the main crawler instance with its configuration and state.
type Crawler struct {
	Config     Config
	Client     *http.Client
	FastClient *http.Client // Client rapide pour HEAD requests
	Visited    sync.Map
	Results    []string
	resultsMu  sync.Mutex
	wg         sync.WaitGroup
	validCache sync.Map // Cache de validation des liens
	semaphore  chan struct{}
}

// New creates and initializes a new Crawler instance with the given configuration.
func New(cfg Config) *Crawler {
	workers := runtime.NumCPU() * 4
	if workers < 16 {
		workers = 16
	}

	transport := &http.Transport{
		TLSClientConfig:     &tls.Config{InsecureSkipVerify: false}, // Default to secure
		MaxIdleConns:        100,
		MaxIdleConnsPerHost: 10,
		MaxConnsPerHost:     20,
		IdleConnTimeout:     30 * time.Second,
		DisableKeepAlives:   false,
	}

	return &Crawler{
		Config: cfg,
		Client: &http.Client{
			Timeout:   60 * time.Second,
			Transport: transport,
		},
		FastClient: &http.Client{
			Timeout:   30 * time.Second,
			Transport: transport,
		},
		semaphore: make(chan struct{}, workers),
	}
}

// Start initiates the crawling process starting from the target URL.
func (c *Crawler) Start() error {
	parsed, err := url.Parse(c.Config.TargetURL)
	if err != nil {
		return err
	}
	norm := parsed.String()

	// Initial check for certificate errors
	if err := c.checkConnection(norm); err != nil {
		return err
	}

	c.Visited.Store(norm, true)

	if err := c.crawl(norm, 0); err != nil {
		return err
	}
	c.wg.Wait()
	return nil
}

func (c *Crawler) checkConnection(targetURL string) error {
	// Try HEAD first
	err := c.doRequest(targetURL, "HEAD")
	if err == nil {
		return nil
	}

	// If HEAD failed with SSL error that was fixed by prompt, doRequest would have retried and succeeded or failed.
	// If it failed with something else (like method allowed or timeout), try GET.
	// We only fallback to GET if the error is NOT a user-aborted SSL check.
	if strings.Contains(err.Error(), "aborted by user") {
		return err
	}

	// Fallback to GET
	if errGet := c.doRequest(targetURL, "GET"); errGet != nil {
		return fmt.Errorf("connection failed (HEAD: %v, GET: %v)", err, errGet)
	}
	return nil
}

func (c *Crawler) doRequest(url, method string) error {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return err
	}

	resp, err := c.FastClient.Do(req)
	if err != nil {
		errStr := strings.ToLower(err.Error())
		if strings.Contains(errStr, "x509") || strings.Contains(errStr, "certificate") || strings.Contains(errStr, "tls") || strings.Contains(errStr, "authority") {
			// Check if we already enabled insecure mode to avoid double prompting
			tr := c.FastClient.Transport.(*http.Transport)
			if tr.TLSClientConfig.InsecureSkipVerify {
				return err // Already insecure, yet failing on SSL? Real error.
			}

			if promptErr := c.promptInsecure(); promptErr != nil {
				return promptErr
			}
			// Retry request with insecure client
			reqRetry, errRetry := http.NewRequest(method, url, nil)
			if errRetry != nil {
				return errRetry
			}
			resp, err = c.FastClient.Do(reqRetry)
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return fmt.Errorf("page not found (404)")
	}
	if resp.StatusCode >= 400 && resp.StatusCode != http.StatusMethodNotAllowed {
		return fmt.Errorf("target returned status %d %s", resp.StatusCode, http.StatusText(resp.StatusCode))
	}
	// If MethodNotAllowed, we return error so fallback can try GET (if we were doing HEAD)
	if resp.StatusCode == http.StatusMethodNotAllowed {
		return fmt.Errorf("method not allowed")
	}

	return nil
}

func (c *Crawler) promptInsecure() error {
	fmt.Printf("%s The target has an invalid/self-signed certificate.\n", color.YellowString("[!]"))
	fmt.Print("Do you want to proceed anyway? [Y/n]: ")

	reader := bufio.NewReader(os.Stdin)
	response, _ := reader.ReadString('\n')
	response = strings.TrimSpace(response)
	response = strings.ToLower(response)

	if response == "" || response == "y" || response == "yes" {
		c.enableInsecure()
		return nil
	}
	return fmt.Errorf("aborted by user: certificate verification failed")
}

func (c *Crawler) enableInsecure() {
	transport := &http.Transport{
		TLSClientConfig:     &tls.Config{InsecureSkipVerify: true},
		MaxIdleConns:        100,
		MaxIdleConnsPerHost: 10,
		MaxConnsPerHost:     20,
		IdleConnTimeout:     30 * time.Second,
		DisableKeepAlives:   false,
	}
	c.Client.Transport = transport
	c.FastClient.Transport = transport
	color.Yellow("[WRN] SSL verification disabled")
}

func (c *Crawler) crawl(rawURL string, depth int) error {
	if depth >= c.Config.MaxDepth {
		return nil
	}
	parsed, err := url.Parse(rawURL)
	if err != nil {
		return err
	}

	resp, err := c.Client.Get(rawURL)
	if err != nil {
		if c.Config.Verbose {
			fmt.Printf("[%s] %s: %v\n", color.RedString("ERR"), rawURL, err)
		}
		return nil
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	links := Extract(string(body))
	validLinks := c.validateLinksParallel(links, parsed)

	for _, linkInfo := range validLinks {
		abs := linkInfo.url
		isExternal := linkInfo.isExternal

		if _, loaded := c.Visited.LoadOrStore(abs, true); loaded {
			continue
		}

		if isExternal {
			if !c.Config.OnlyInternal {
				fmt.Printf("[%s] %s\n", color.CyanString("EXT"), abs)
				c.addResult(abs)
			}
		} else {
			if !c.Config.OnlyExternal {
				fmt.Printf("[%s] %s\n", color.GreenString("INT"), abs)
				c.addResult(abs)
			}

			c.wg.Add(1)
			go func(url string, d int) {
				defer c.wg.Done()
				c.semaphore <- struct{}{}
				defer func() { <-c.semaphore }()
				c.crawl(url, d+1)
			}(abs, depth)
		}
	}
	return nil
}

type linkInfo struct {
	url        string
	isExternal bool
}

func (c *Crawler) validateLinksParallel(links []string, baseURL *url.URL) []linkInfo {
	results := make(chan linkInfo, len(links))
	var wg sync.WaitGroup

	for _, link := range links {
		wg.Add(1)
		go func(l string) {
			defer wg.Done()
			c.semaphore <- struct{}{}
			defer func() { <-c.semaphore }()

			res, err := baseURL.Parse(l)
			if err != nil {
				return
			}
			abs := res.String()
			isExternal := res.Host != baseURL.Host

			if c.Config.OnlyInternal && isExternal {
				return
			}
			if c.validateLink(abs) {
				results <- linkInfo{
					url:        abs,
					isExternal: isExternal,
				}
			}
		}(link)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	var validated []linkInfo
	for li := range results {
		validated = append(validated, li)
	}
	return validated
}

func (c *Crawler) validateLink(u string) bool {
	if cached, ok := c.validCache.Load(u); ok {
		return cached.(bool)
	}

	req, err := http.NewRequest("HEAD", u, nil)
	if err != nil {
		c.validCache.Store(u, false)
		return false
	}

	resp, err := c.FastClient.Do(req)
	if err != nil {
		if c.Config.Verbose {
			fmt.Printf("[%s] %s: %v\n", color.RedString("ERR"), u, err)
		}
		c.validCache.Store(u, false)
		return false
	}
	defer resp.Body.Close()

	valid := resp.StatusCode >= 200 && resp.StatusCode < 400
	c.validCache.Store(u, valid)
	return valid
}

func (c *Crawler) addResult(url string) {
	c.resultsMu.Lock()
	c.Results = append(c.Results, url)
	c.resultsMu.Unlock()
}

// SaveJSON exports the crawling results (and tree if enabled) to a JSON file.
func (c *Crawler) SaveJSON() error {
	if c.Config.OutputPath == "" {
		return nil
	}
	type Export struct {
		Target  string    `json:"target"`
		Results []string  `json:"results"`
		Tree    *treeNode `json:"tree,omitempty"`
		Count   int       `json:"count"`
	}

	var tree *treeNode
	if c.Config.ShowTree {
		tree = c.buildTree()
	}

	data := Export{
		Target:  c.Config.TargetURL,
		Results: c.Results,
		Tree:    tree,
		Count:   len(c.Results),
	}
	file, err := os.Create(c.Config.OutputPath)
	if err != nil {
		return err
	}
	defer file.Close()
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(data)
}

type treeNode struct {
	Name     string               `json:"name"`
	Children map[string]*treeNode `json:"children,omitempty"`
}

func newTreeNode(name string) *treeNode {
	return &treeNode{
		Name:     name,
		Children: make(map[string]*treeNode),
	}
}

// PrintTree outputs the internal directory structure tree to stdout.
func (c *Crawler) PrintTree() {
	if !c.Config.ShowTree {
		return
	}
	fmt.Printf("\n%s\n%s\n", color.MagentaString("=== Site Tree ==="), c.Config.TargetURL)

	root := c.buildTree()
	c.printRecursive(root, "")
}

func (c *Crawler) printRecursive(node *treeNode, prefix string) {
	keys := make([]string, 0, len(node.Children))
	for k := range node.Children {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for i, name := range keys {
		isLast := i == len(keys)-1
		connector := "├── "
		if isLast {
			connector = "└── "
		}
		fmt.Printf("%s%s%s\n", prefix, connector, name)

		newPrefix := prefix + "│   "
		if isLast {
			newPrefix = prefix + "    "
		}
		c.printRecursive(node.Children[name], newPrefix)
	}
}

func (c *Crawler) buildTree() *treeNode {
	rootURL, _ := url.Parse(c.Config.TargetURL)
	root := newTreeNode("/")

	urls := append([]string{c.Config.TargetURL}, c.Results...)
	for _, uStr := range urls {
		u, err := url.Parse(uStr)
		if err != nil || u.Host != rootURL.Host {
			continue
		}

		path := u.Path
		if path == "" {
			path = "/"
		}

		suffix := ""
		if u.RawQuery != "" {
			suffix = "?" + u.RawQuery
		}

		parts := strings.Split(path, "/")
		current := root

		for i, part := range parts {
			if part == "" {
				continue
			}
			name := part
			if i == len(parts)-1 {
				name += suffix
			}
			if _, exists := current.Children[name]; !exists {
				current.Children[name] = newTreeNode(name)
			}
			current = current.Children[name]
		}

		if path == "/" && suffix != "" {
			name := "?" + u.RawQuery
			if _, exists := root.Children[name]; !exists {
				root.Children[name] = newTreeNode(name)
			}
		}
	}
	return root
}
