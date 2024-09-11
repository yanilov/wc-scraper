# WC parser

This repo contains my solution to building a concurrent scraper which generates a top-k list of word frequencies.

## Build and run

execute the following from the root directory of the repo
```sh
go run cmd/scraper/main.go
```

## Future improvements

* Bloom filter to support larger dictionaries
* Benchmarks
* CLI, file-based configuration instead of hardcoded
* Test core scraper flow control