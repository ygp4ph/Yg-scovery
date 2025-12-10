# Yg-scovery

![License](https://img.shields.io/badge/license-MIT-blue)
![Go Version](https://img.shields.io/badge/go-1.20%2B-cyan)

**Yg-scovery** is a powerful, recursive content discovery tool written in Go. It goes beyond simple URL fuzzing by analyzing the source code of the target page (HTML, JavaScript, JSON) to extract hidden endpoints, static assets, and API routes.

Designed with a modern, colored CLI interface inspired by the best security tools.

## Features

- 🕸️ **Recursive Crawling**: Automatically traverses found internal links to discover deep content.
- 🔍 **Source Code Analysis**: Extracts URLs and paths from HTML, JS scripts, and JSON responses using robust regex patterns.
- 🎨 **Visual Output**: Beautiful, colored output with distinct tags for (`[INT]`, `[EXT]`, `[FOUND]`).
- ⚡ **Fast & Lightweight**: Built with Go's efficiency.
- 🎯 **Smart Filtering**: Filter results to show only internal (`-i`) or external (`-e`) links.

## Installation

Ensure you have Go installed (1.20+).

```bash
git clone https://github.com/ygp4ph/Yg-scovery.git
cd Yg-scovery
go build -o yg-scovery .
```

## Usage

```bash
./yg-scovery -u <TARGET_URL> [flags]
```

### Flags

| Flag | Alias | Description |
|------|-------|-------------|
| `-u` | `--url` | Target URL to scan (Required) |
| `-d` | `--depth` | Maximum recursion depth (default 3) |
| `-i` | `--int` | Show only internal links |
| `-e` | `--ext` | Show only external links |
| `-h` | `--help` | Show help manual |

### Examples

**Basic Recursive Scan:**
```bash
./yg-scovery -u https://example.com
```

**Internal Only with Depth 5:**
```bash
./yg-scovery -u https://example.com -d 5 -i
```

## Screenshot / Output

```bash

```

## Contributing

Contributions are welcome! Feel free to open issues or submit pull requests.

## License

This project is licensed under the MIT License.
