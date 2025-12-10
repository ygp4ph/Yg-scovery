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
	flag.StringVar(&u, "u", "", "")
	flag.StringVar(&u, "url", "", "")
	flag.IntVar(&d, "d", 3, "")
	flag.IntVar(&d, "depth", 3, "")
	flag.BoolVar(&onlyExternal, "e", false, "")
	flag.BoolVar(&onlyExternal, "ext", false, "")
	flag.BoolVar(&onlyInternal, "i", false, "")
	flag.BoolVar(&onlyInternal, "int", false, "")
	flag.BoolVar(&h, "h", false, "")
	flag.BoolVar(&h, "help", false, "")

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
		fmt.Fprintf(os.Stderr, "\nUSAGE: %s [flags]\n\nFLAGS:\n  -u, --url\tTarget URL\n  -d, --depth\tMax recursion (default 3)\n  -e, --ext\tExternal links only\n  -i, --int\tInternal links only\n  -h, --help\tShow help\n", os.Args[0])
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
	if onlyExternal && onlyInternal { // Updated variable names
		color.Red("[ERR] Cannot use both --ext and --int flags together.\n") // Updated error message
		os.Exit(1)
	}

	color.Green("[INF] Scanning %s (Depth: %d)", u, d)
	if onlyExternal {
		color.Yellow("[INF] Filter: External links only")
	}
	if onlyInternal {
		color.Yellow("[INF] Filter: Internal links only")
	}

	cfg := Config{ // Use local Config struct
		TargetURL:    u,
		MaxDepth:     d,
		OnlyInternal: onlyInternal, // Updated variable name
		OnlyExternal: onlyExternal, // Updated variable name
	}

	c := New(cfg) // Use local New function

	if err := c.Start(); err != nil {
		log.Fatalf("%s %v", color.RedString("[FATAL] Crawler failed:"), err) // Updated error message
	}
}
