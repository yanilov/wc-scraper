# WC parser

This repo contains my solution to building a concurrent scraper which generates a top-k list of word frequencies.

## Run the CLI

Execute the following from the root directory of the repo.\

```sh
go run cmd/scraper/main.go 
```

You can also view help on commnands and flags (currently only the root command exists)

```sh
go run cmd/scraper/main.go --help
```

**Note:** when running in a terminal, the above will print stderr to terminal as well. To avoid seeing stderr at all, run

```sh
 go run cmd/scraper/main.go 2>/dev/null
 ```

### Specify configuration

To specify configuration file, overriding the default, run

 ```sh
  go run cmd/scraper/main.go --config ./cmd/scraper/scrapeconfig.yaml
 ```

 See the provided `cmd/scraper/config.yaml` for example
 **Note:** the `.job_loader.page_cutoff` parameter is set to 3, to avoid getting accidentally blocked during development.\
 For a real-world use case, either omit the attribute or set it to 0.

### Specify log level

To set the log level, add a `GO_LOG=<level>` envvar when executing, e.g.
 
 ```sh
 GO_LOG=debug go run cmd/scraper/main.go
 ```
 
 When log level is unspecified, the default is "info".

 ## Run Tests
 
 Execute the following from the root directory of the repo.
 ```sh
go test ./...
 ```



## Future improvements

* Supprort multi-node architecture for large-scale scraping
* Bloom filter to support larger dictionaries
* Benchmarks
* Test core scraper flow control