# brim

A terminal-based search interface with vim-style navigation, built on a self-hosted [SearXNG](https://github.com/searxng/searxng) instance.

## Why

Modern search results are noisy. `brim` gives you a fast, keyboard-driven TUI to query SearXNG (which aggregates 70+ engines) and open results directly in your browser — no ads, no tracking, no leaving the terminal.

## Requirements

- Go 1.25+
- Docker (for running SearXNG locally)

## Setup

### 1. Run SearXNG locally

```bash
mkdir -p ~/searxng && cd ~/searxng
cat > settings.yml << 'EOF'
use_default_settings: true
server:
  secret_key: "CHANGE_ME"
  limiter: false
search:
  formats:
    - html
    - json
EOF

docker run -d --name searxng \
  -p 8080:8080 \
  -v ~/searxng/settings.yml:/etc/searxng/settings.yml:ro \
  -e SEARXNG_BASE_URL=http://localhost:8080/ \
  --restart unless-stopped \
  docker.io/searxng/searxng:latest
```

Verify it's working:

```bash
curl -s "http://localhost:8080/search?q=test&format=json" | head
```

### 2. Build brim

```bash
git clone https://github.com/simar-malhotra09/brim.git
cd brim
go build .
./brim
```

Or run directly:

```bash
go run .
```

## Usage

| Key       | Action                          |
|-----------|---------------------------------|
| (type)    | Enter a query                   |
| `enter`   | Run search / open selected URL  |
| `j` / `k` | Move down / up in results       |
| `esc`     | Back to search input            |
| `ctrl+c`  | Quit                            |

## Configuration

By default, `brim` talks to `http://localhost:8080`. To point it elsewhere, edit `provider/provider.go` and change `defaultBaseURL`.

## Project Structure

```
brim/
├── main.go           # TUI (Bubble Tea)
├── provider/         # SearXNG client
└── utils/            # cross-platform helpers (URL opener, etc.)
```

## License

MIT
