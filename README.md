# typical-rest-server

[![Project Status: WIP – Initial development is in progress, but there has not yet been a stable, usable release suitable for the public.](https://www.repostatus.org/badges/latest/wip.svg)](https://www.repostatus.org/#wip)

Example of typical and scalable RESTful API Server for Go

## Prerequisite

Install [Go](https://golang.org/doc/install) (It is recommend to install via [Homebrew](https://brew.sh/) `brew install go`)

## Quick Start

```bash
# equivalent with `docker-compose up -d` (if infrastructure not up)
./typicalw docker up 

# drop, create and migrate postgres database (if database not ready)
./typicalw pg reset 

# generate mock (if require mock)
./typicalw mock 

# run test 
./typicalw test

# run the application
./typicalw run 

# release the distribution
./typicalw release 
```

## TODO

TypicalGo support java-like annotation for code generation before build
- [x] `[constructor]` to add constructor in service-locator 
- [x] `[mock]` is target to be mock for `./typicalw mock`
- [ ] `[api(method=<METHOD>, path=<PATH>)]` to generate route function (TODO)
- [ ] `[cacheable]` to generate cached function (TODO)

