# typical-rest-server

[![Project Status: WIP – Initial development is in progress, but there has not yet been a stable, usable release suitable for the public.](https://www.repostatus.org/badges/latest/wip.svg)](https://www.repostatus.org/#wip)
![Go-Workflow](https://github.com/typical-go/typical-rest-server/workflows/Go/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/typical-go/typical-rest-server)](https://goreportcard.com/report/github.com/typical-go/typical-rest-server)
[![codebeat badge](https://codebeat.co/badges/17e19d4b-6803-4bbb-82bb-e39fe2f1424b)](https://codebeat.co/projects/github-com-typical-go-typical-rest-server-master)
[![codecov](https://codecov.io/gh/typical-go/typical-rest-server/branch/master/graph/badge.svg)](https://codecov.io/gh/typical-go/typical-rest-server)


The aims of project is provide simple, straight-forward and easy-to-adopt Rest-Server implementation in Go.

## How to Use 
 
- Use the project as reference for project layout and code best practice.
- The `pkg` package is general library that can help in various needs

## Prerequisite

- [Go](https://golang.org/doc/install) 
- [Docker-Compose](https://docs.docker.com/compose/install/)

## Quick Start

The project using [typical-go](https://github.com/typical-go/typical-go) as its build-tool.

```bash
./typicalw docker up   # equivalent with `docker-compose up -d` (if infrastructure not up)
./typicalw pg reset    # drop, create and migrate postgres database (if database not ready)
./typicalw mock        # generate mock (if require mock)

./typicalw test        # run test 

./typicalw run         # run the application
```

## Project Layout

- `internal` Exclusive go source files for the project
  - `internal/app` main functionality/entry-point
  - `internal/infra` infrastructure for the project e.g. config and connection object
  - `internal/profiler` for profiling the project e.g. HealthCheck, PProf, etc
  - `internal/server` REST-Server actual code
- `pkg` Shareable go source files e.g. helper/utitily Library
- `api` Any related scripts for API e.g. api-model script (swagger, raml, etc) or client script
- `databases` Any related scripts for Databases e.g. migration scripts and seed data
- `tools` External program that used by the project
- `typical` Build-Tool descriptor

## Layered Architecture

- `controller`
  - Handling HTTP routes
  - Parsing the request
  - Sending response (both success & error)
- `service`
  - Intermediary between controller (end-point) and repository (data)
  - Logic of controller
  - Data Validation
  - DTO (Data Transfer Object) Model
- `repository`
  - Data access layer 
  - No logic except operation to database
  - Repository pattern
  - DAO (Data Access Object) Model
  - Database Entity or Business Entity

## Annotations

The project utilize Typical-Go annotation for code generation
- `@ctor` functions will be provided in [dig](https://github.com/uber-go/dig) service-locator
- `@dtor` functions will be execute when program end
- `@mock` interfaces is target for mock generation

## Configurations

Config in environment as [12factor](https://12factor.net/config) recommentation. The build-tool (typical-go) will load values in `.env` file and generated it if not available.

## Pg-Tool

Simple CLI to prepare the database [learn more](tools/pg-tool/README.md)

```bash
./typicalw pg
```

## Pkg Library

- `dbkit` Utility for database operation
- `dbtxn` Utility for database transaction
- `dockerrx` Docker Recipe Collection to be generated by typical-go
- `echokit` Utility for controller e.g. create HealthCheck api or interface for echo.Server
- `echotest` Utility for table-driven test echo handler
- `errvalid` wrapper for validation error
- `pgcmd` Command collection for postgres
- `rediscmd` Command collection for redis

## 3rd Party Library

- [Echo Framework](https://echo.labstack.com/): Web Framework
- [dig](https://github.com/uber-go/dig): Dependency Injection
- [gomock](https://github.com/golang/mock): Mock for interface
- [logrus](https://github.com/sirupsen/logrus): Logging
- [testify](https://github.com/stretchr/testify): Test assertion
- [squirrel](https://github.com/Masterminds/squirrel): SQL Query Builder
- [go-sqlmock](https://github.com/DATA-DOG/go-sqlmock): Mock for database connection 


## Concepts

- [SOLID](https://en.wikipedia.org/wiki/SOLID)
- [Dependency Injection](https://en.wikipedia.org/wiki/Dependency_injection) 
- [Table-Driven Test](https://github.com/golang/go/wiki/TableDrivenTests)
- [Mock Object](https://en.wikipedia.org/wiki/Mock_object) 
- [ORM-Hate](https://martinfowler.com/bliki/OrmHate.html)

## Golang References

- [Go Documentation](https://golang.org/doc/)
- [Go For Industrial Programming](https://peter.bourgon.org/go-for-industrial-programming/)
- [Uber Go Style Guide](https://github.com/uber-go/guide)
- [Standard Go Project Layout](https://github.com/golang-standards/project-layout/)
- [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)

## Project Status

The project status is `WIP` (Work in progress), means the author plan to continously improve and add more best-practices.

## License

This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details