package main

import (
	"regexp"
	"strings"
)

// regex pour extraire les liens
var (
	urlRegex  = regexp.MustCompile(`https?://[a-zA-Z0-9\-\.]+\.[a-zA-Z]{2,}(?:/[^"'\s<>` + "`" + `]*)?`)
	pathRegex = regexp.MustCompile(`["'](\.?\.?/[^"'\s<>` + "`" + `]+)["']`)
	attrRegex = regexp.MustCompile(`(href|src)=["']([^"']+)["']`)
)

// Extract extrait les liens d'une page web.
func Extract(content string) []string {
	seen := make(map[string]bool)
	var found []string
	add := func(s string) {
		if !seen[s] && len(s) > 1 && !strings.ContainsAny(s, "\n ") {
			found = append(found, s)
			seen[s] = true
		}
	}

	for _, m := range urlRegex.FindAllString(content, -1) {
		add(m)
	}
	for _, m := range pathRegex.FindAllStringSubmatch(content, -1) {
		if len(m) > 1 {
			add(m[1])
		}
	}
	for _, m := range attrRegex.FindAllStringSubmatch(content, -1) {
		if len(m) > 2 {
			add(m[2])
		}
	}
	return found
}
