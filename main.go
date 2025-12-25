package main

import (
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"

	"github.com/fatih/color"
)

var Version = "v2.1.0"

func main() {
	var (
		u                          string
		d                          int
		onlyExternal, onlyInternal bool
		output                     string
		h, verbose, showVersion    bool
		tree                       bool
	)

	flag.StringVar(&u, "u", "", "Target URL")
	flag.StringVar(&u, "url", "", "Target URL")
	flag.IntVar(&d, "d", 3, "Max recursion depth")
	flag.IntVar(&d, "depth", 3, "Max recursion depth")
	flag.BoolVar(&onlyExternal, "e", false, "External links only")
	flag.BoolVar(&onlyExternal, "ext", false, "External links only")
	flag.BoolVar(&onlyInternal, "i", false, "Internal links only")
	flag.BoolVar(&onlyInternal, "int", false, "Internal links only")
	flag.StringVar(&output, "o", "", "Output file (JSON)")
	flag.StringVar(&output, "output", "", "Output file (JSON)")
	flag.BoolVar(&tree, "t", false, "Show internal links tree")
	flag.BoolVar(&tree, "tree", false, "Show internal links tree")
	flag.BoolVar(&h, "h", false, "Show help")
	flag.BoolVar(&h, "help", false, "Show help")
	flag.BoolVar(&verbose, "v", false, "Show errors")
	flag.BoolVar(&verbose, "verbose", false, "Show errors")
	flag.BoolVar(&showVersion, "version", false, "Show version")

	banner := func() {
		color.Cyan(`
   __  ______ _      ______________ _   _____  _______  __
  / / / / __ `+"`"+`/_____/ ___/ ___/ __ \ | / / _ \/ ___/ / / /
 / /_/ / /_/ /_____(__  ) /__/ /_/ / |/ /  __/ /  / /_/ / 
 \__, /\__, /     /____/\___/\____/|___/\___/_/   \__, /  
/____//____/                                     /____/   %s
 `, Version)
	}

	flag.Usage = func() {
		banner()
		fmt.Fprintf(os.Stderr, "\nUSAGE: %s [flags]\n\nFLAGS:\n  -u, --url\tTarget URL\n  -d, --depth\tMax recursion (default 3)\n  -e, --ext\tExternal links only\n  -i, --int\tInternal links only\n  -t, --tree\tShow internal links tree\n  -o, --output\tOutput file (JSON)\n  -v, --verbose\tShow errors\n  --version\tShow version\n  -h, --help\tShow help\n", os.Args[0])
	}
	flag.Parse()

	if h {
		flag.Usage()
		os.Exit(0)
	}

	if showVersion {
		fmt.Printf("yg-scovery %s\n", Version)
		os.Exit(0)
	}

	banner()
	if u == "" {
		color.Red("[ERR] -u <url> required")
		fmt.Println("Use -h for help")
		os.Exit(1)
	}
	if !strings.HasPrefix(u, "http://") && !strings.HasPrefix(u, "https://") {
		u = "https://" + u
	}
	if _, err := url.Parse(u); err != nil {
		color.Red("[ERR] Invalid URL: %v", err)
		os.Exit(1)
	}
	if onlyExternal && onlyInternal {
		color.Red("[ERR] Conflict: -e and -i")
		os.Exit(1)
	}

	color.Green("[INF] Scanning %s (Depth: %d)", u, d)
	if onlyExternal {
		color.Yellow("[INF] Filter: External links only")
	}
	if onlyInternal {
		color.Yellow("[INF] Filter: Internal links only")
	}
	if tree {
		color.Magenta("[INF] Tree view enabled (Internal links)")
	}
	if output != "" {
		color.Blue("[INF] Output will be saved to %s", output)
	}

	cfg := Config{
		TargetURL:    u,
		MaxDepth:     d,
		OnlyInternal: onlyInternal,
		OnlyExternal: onlyExternal,
		OutputPath:   output,
		Verbose:      verbose,
		ShowTree:     tree,
	}

	c := New(cfg)
	if err := c.Start(); err != nil {
		log.Fatalf("%s %v", color.RedString("[FATAL] Crawler failed:"), err)
	}

	if tree {
		c.PrintTree()
	}

	if output != "" {
		if err := c.SaveJSON(); err != nil {
			color.Red("[ERR] Failed to save output: %v", err)
		} else {
			color.Green("[INF] Saved results to %s", output)
		}
	}
}
