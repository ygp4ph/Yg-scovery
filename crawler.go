package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"sync"
	"time"

	"github.com/fatih/color"
)

// Config représente les paramètres de configuration du crawler.
type Config struct {
	TargetURL    string
	MaxDepth     int
	OnlyInternal bool
	OnlyExternal bool
	OutputPath   string
}

// Crawler représente un crawler web.
type Crawler struct {
	Config  Config
	Client  *http.Client
	Visited map[string]bool
	mu      sync.Mutex
	Results []string
}

// New crée un nouveau crawler.
func New(cfg Config) *Crawler {
	return &Crawler{
		Config: cfg,
		Client: &http.Client{
			Timeout: 10 * time.Second,
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			},
		},
		Visited: make(map[string]bool),
	}
}

// Start lance le crawler.
func (c *Crawler) Start() error {
	parsed, err := url.Parse(c.Config.TargetURL)
	if err != nil {
		return err
	}
	norm := parsed.String()
	c.mu.Lock()
	c.Visited[norm] = true
	c.mu.Unlock()
	return c.crawl(norm, 0)
}

// crawl effectue la recherche.
func (c *Crawler) crawl(rawURL string, depth int) error {
	if depth > c.Config.MaxDepth {
		return nil
	}
	parsed, err := url.Parse(rawURL)
	if err != nil {
		return err
	}

	resp, err := c.Client.Get(rawURL)
	if err != nil {
		fmt.Printf("[%s] %s: %v\n", color.RedString("ERR"), rawURL, err)
		return nil
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	for _, link := range Extract(string(body)) {
		res, err := parsed.Parse(link)
		if err != nil {
			continue
		}
		abs := res.String()

		c.mu.Lock()
		if c.Visited[abs] {
			c.mu.Unlock()
			continue
		}
		c.Visited[abs] = true
		c.mu.Unlock()

		if res.Host != parsed.Host {
			if !c.Config.OnlyInternal {
				fmt.Printf("[%s] %s\n", color.CyanString("EXT"), abs)
				c.mu.Lock()
				c.Results = append(c.Results, abs)
				c.mu.Unlock()
			}
		} else {
			if !c.Config.OnlyExternal {
				fmt.Printf("[%s] %s\n", color.GreenString("INT"), abs)
				c.mu.Lock()
				c.Results = append(c.Results, abs)
				c.mu.Unlock()
			}
			c.crawl(abs, depth+1)
		}
	}
	return nil
}

// SaveJSON sauvegarde les résultats dans un fichier JSON.
func (c *Crawler) SaveJSON() error {
	if c.Config.OutputPath == "" {
		return nil
	}
	type Export struct {
		Target  string   `json:"target"`
		Results []string `json:"results"`
		Count   int      `json:"count"`
	}
	data := Export{
		Target:  c.Config.TargetURL,
		Results: c.Results,
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
