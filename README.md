# typical-rest-server

[![Project Status: WIP – Initial development is in progress, but there has not yet been a stable, usable release suitable for the public.](https://www.repostatus.org/badges/latest/wip.svg)](https://www.repostatus.org/#wip)

Yet Another REST-Server Template/Example Project for Go

## How to Use 
 
- The project status is `WIP` (Work in progress), means the author plan to continously improve and add more best-practices.
- Use the project as reference for project layout and code best practice.
- The `pkg` package is general library that can help in various needs
- Feel free to adopt and contributing 


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
  - Entity/DAO (Data Access Object) Model

## Annotations

The project utilize Typical-Go annotation for code generation
- `@ctor` functions will be provided in [dig](https://github.com/uber-go/dig) service-locator
- `@dtor` functions will be execute when program end
- `@mock` interfaces is target for mock generation

## Configurations

Config in environment as [12factor](https://12factor.net/config) recommentation. The build-tool (Typical-Go) will load values in `.env` file and generated it if not available.

## Pg-Tool

Simple CLI to prepare the database [learn more](tools/pg-tool/README.md)

```bash
./typicalw pg
```

## Concepts

- [Dependency inversion principle](https://en.wikipedia.org/wiki/SOLID)
- [Dependency Injection](https://en.wikipedia.org/wiki/Dependency_injection) ([dig](https://github.com/uber-go/dig))
- [Table-Driven Test](https://github.com/golang/go/wiki/TableDrivenTests)
- [Mock Object](https://en.wikipedia.org/wiki/Mock_object) ([gomock](https://github.com/golang/mock))
- [ORM-Hate](https://martinfowler.com/bliki/OrmHate.html)

## Golang References

- [Go Documentation](https://golang.org/doc/)
- [Go For Industrial Programming](https://peter.bourgon.org/go-for-industrial-programming/)
- [Uber Go Style Guide](https://github.com/uber-go/guide)
- [Standard Go Project Layout](https://github.com/golang-standards/project-layout/)
