# typical-rest-server

[![Project Status: WIP – Initial development is in progress, but there has not yet been a stable, usable release suitable for the public.](https://www.repostatus.org/badges/latest/wip.svg)](https://www.repostatus.org/#wip)

Example of typical and scalable RESTful API Server for Go
- Classic Architecture
  - `controller`: for handling routes, request parsing and sending response
  - `service`: logic of controller, validation, intermediary between controller (end-point) and repository (data).
  - `repository`: data access layer using repository pattern
- Interface Segregation Principle (ISP) between layers
- Dependency Injection
- Profiler
  - HealthCheck end-point
- Mock using gomock 
- Database Utility
  - Create/Drop Database
  - Migrate/Rollback Database
  - Seed Database
  - Console
- Table-Driven Test

## Prerequisite

Install [Go](https://golang.org/doc/install) (It is recommend to install via [Homebrew](https://brew.sh/) `brew install go`)

## Quick Start

The project using [typical-go](https://github.com/typical-go/typical-go) as its build-tool.

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
```
