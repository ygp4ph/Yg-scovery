package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/fatih/color"
)

func main() {
	var u string
	var d int
	var onlyExternal, onlyInternal, h bool
	var output string
	// Définition des flags et leurs valeurs par défaut
	flag.StringVar(&u, "u", "", "")
	flag.StringVar(&u, "url", "", "")
	flag.IntVar(&d, "d", 3, "")
	flag.IntVar(&d, "depth", 3, "")
	flag.BoolVar(&onlyExternal, "e", false, "")
	flag.BoolVar(&onlyExternal, "ext", false, "")
	flag.BoolVar(&onlyInternal, "i", false, "")
	flag.BoolVar(&onlyInternal, "int", false, "")
	flag.StringVar(&output, "o", "", "")
	flag.StringVar(&output, "output", "", "")
	flag.BoolVar(&h, "h", false, "")
	flag.BoolVar(&h, "help", false, "")

	// hehe
	banner := func() {
		color.Cyan(`
_____.___.                  _________                                      
\__  |   | ____            /   _____/ ____  _______  __ ___________ ___.__.
 /   |   |/ ___\   ______  \_____  \_/ ___\/  _ \  \/ // __ \_  __ <   |  |
 \____   / /_/  > /_____/  /        \  \__(  <_> )   /\  ___/|  | \/\___  |
 / ______\___  /          /_______  /\___  >____/ \_/  \___  >__|   / ____|
 \/     /_____/                   \/     \/                \/       \/      v1.0.0
 `)
	}

	flag.Usage = func() {
		banner()
		fmt.Fprintf(os.Stderr, "\nUSAGE: %s [flags]\n\nFLAGS:\n  -u, --url\tTarget URL\n  -d, --depth\tMax recursion (default 3)\n  -e, --ext\tExternal links only\n  -i, --int\tInternal links only\n  -o, --output\tOutput file (JSON)\n  -h, --help\tShow help\n", os.Args[0])
	}
	flag.Parse()

	if h {
		flag.Usage()
		os.Exit(0)
	}

	banner()
	if u == "" {
		color.Red("[ERR] -u <url> required")
		fmt.Println("Use -h for help")
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
	if output != "" {
		color.Blue("[INF] Output will be saved to %s", output)
	}

	cfg := Config{
		TargetURL:    u,
		MaxDepth:     d,
		OnlyInternal: onlyInternal,
		OnlyExternal: onlyExternal,
		OutputPath:   output,
	}

	c := New(cfg)
	if err := c.Start(); err != nil {
		log.Fatalf("%s %v", color.RedString("[FATAL] Crawler failed:"), err)
	}

	if output != "" {
		if err := c.SaveJSON(); err != nil {
			color.Red("[ERR] Failed to save output: %v", err)
		} else {
			color.Green("[INF] Saved results to %s", output)
		}
	}
}
