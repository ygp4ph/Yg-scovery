package main

import (
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/fatih/color"
)

type Config struct {
	TargetURL    string
	MaxDepth     int
	OnlyInternal bool
	OnlyExternal bool
}

type Crawler struct {
	Config  Config
	Client  *http.Client
	Visited map[string]bool
	mu      sync.Mutex
}

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

func (c *Crawler) Start() error { return c.crawl(c.Config.TargetURL, 0) }

func (c *Crawler) crawl(rawURL string, depth int) error {
	if depth > c.Config.MaxDepth {
		return nil
	}
	parsed, err := url.Parse(rawURL)
	if err != nil {
		return err
	}
	norm := parsed.String()

	c.mu.Lock()
	if c.Visited[norm] {
		c.mu.Unlock()
		return nil
	}
	c.Visited[norm] = true
	c.mu.Unlock()

	resp, err := c.Client.Get(norm)
	if err != nil {
		fmt.Printf("[%s] %s: %v\n", color.RedString("ERR"), norm, err)
		return nil
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	for _, link := range Extract(string(body)) { // Call Extract directly
		res, err := parsed.Parse(link)
		if err != nil {
			continue
		}
		abs := res.String()
		if res.Host != parsed.Host {
			if !c.Config.OnlyInternal {
				fmt.Printf("[%s] %s\n", color.CyanString("EXT"), abs)
			}
		} else {
			if !c.Config.OnlyExternal {
				fmt.Printf("[%s] %s\n", color.GreenString("INT"), abs)
			}
			c.crawl(abs, depth+1)
		}
	}
	return nil
}
