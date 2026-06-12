# kys-scraper

URL scraping service for [Know Your Stash](https://github.com/yourname/kys).
Accepts a URL and returns structured data extracted from the page, with
domain-specific scrapers for supported sites and a generic HTML fallback.

## Supported sources

| Domain          | Result types                  |
|-----------------|-------------------------------|
| `dc.fandom.com` | `issues`, `issue series`      |

## Requirements

- Go 1.26+
- Docker (optional)

## Running locally

```bash
go mod tidy
go run ./cmd/server/main.go
```

## Running with Docker

```bash
docker compose up --build
```

## Environment variables

| Variable         | Default | Description                    |
|------------------|---------|--------------------------------|
| `PORT`           | `8081`  | Port the server listens on     |
| `ALLOWED_ORIGIN` | `*`     | CORS allowed origin            |
| `GIN_MODE`       | `debug` | Set to `release` in production |

## Adding a new scraper

1. Add a result struct in `internal/scraper/results/`
2. Implement the classifier in the domain scraper file
3. Register it in `internal/registry/registry.go`
