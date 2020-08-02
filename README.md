[![Project Status: WIP – Initial development is in progress, but there has not yet been a stable, usable release suitable for the public.](https://www.repostatus.org/badges/latest/wip.svg)](https://www.repostatus.org/#wip)
![Go-Workflow](https://github.com/typical-go/typical-rest-server/workflows/Go/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/typical-go/typical-rest-server)](https://goreportcard.com/report/github.com/typical-go/typical-rest-server)
[![codebeat badge](https://codebeat.co/badges/17e19d4b-6803-4bbb-82bb-e39fe2f1424b)](https://codebeat.co/projects/github-com-typical-go-typical-rest-server-master)
[![codecov](https://codecov.io/gh/typical-go/typical-rest-server/branch/master/graph/badge.svg)](https://codecov.io/gh/typical-go/typical-rest-server)

# typical-rest-server

Opinionated, simple and straight forward Restful server implementation for Golang.

## How to Use 

- Use the project as reference for project layout and code best practice.
- The [pkg](pkg) package is [shared library](#shared-library) that can help in various needs.

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

Typical-Rest encourage [standard go project layout](https://github.com/golang-standards/project-layout)

- [`internal`](internal) Exclusive go source files for the project
  - [`internal/app`](internal/app) main functionality/entry-point
  - [`internal/infra`](internal/infra) infrastructure for the project e.g. config and connection object
  - [`internal/profiler`](internal/profiler) for profiling the project e.g. HealthCheck, PProf, etc
  - [`internal/server`](internal/server) REST-Server actual code
- [`pkg`](pkg) Shareable go source files e.g. helper/utitily Library
- [`api`](api) Any related scripts for API e.g. api-model script (swagger, raml, etc) or client script
- [`databases`](database) Any related scripts for Databases e.g. migration scripts and seed data
- [`tools`](tool) Supporting tool for the project

## Layered Architecture

Typical-Rest encourage [layered architecture](https://en.wikipedia.org/wiki/Multitier_architecture) (as most adoptable architectural pattern) with [SOLID Principle](https://en.wikipedia.org/wiki/SOLID) and [Table-Driven Test](https://github.com/golang/go/wiki/TableDrivenTests)

- Presentation Layer at [`internal/server/controller`](internal/server/controller)
  - Handling HTTP routes
  - Parsing the request
  - Sending response (both success & error)
- Logic Layer at [`internal/server/service`](internal/server/service)
  - Intermediary between controller (end-point) and repository (data)
  - Logic of controller
  - Data Validation
  - DTO (Data Transfer Object) Model
- Data Access Layer at [`internal/server/repository`](internal/server/repository)
  - No logic except operation to database
  - Repository pattern
  - DAO (Data Access Object) Model
  - Database Entity or Business Entity

## Dependency Injection

Typical-Rest encourage [dependency injection](https://en.wikipedia.org/wiki/Dependency_injection) using [uber-dig](https://github.com/uber-go/dig) and annotations (`@ctor` for constructor and `@dtor` for destructor).

```go
// OpenConn open new database connection
// @ctor
func OpenConn() *sql.DB{
}
```

```go
// CloseConn close the database connection
// @dtor
func CloseConn(db *sql.DB){
}
```

## Application Config

Typical-Rest encourage [application config with environment variables](https://12factor.net/config) using [envconfig](https://github.com/kelseyhightower/envconfig) and annotation (`@app-cfg`). 

```go
type (
  // AppCfg application configuration
  // @app-cfg (prefix:"APP")
  AppCfg struct {
    Address string `envconfig:"ADDRESS" default:":8089" required:"true"`
    Debug   bool   `envconfig:"DEBUG" default:"true"`
  }
)
```

## Mocking

Typical-Rest encourage [mocking](https://en.wikipedia.org/wiki/Mock_object) using [gomock](https://github.com/golang/mock) and annotation(`@mock`). 

```go
type(
  // Reader responsible to read
  // @mock
  Reader interface{
    Read() error
  }
)
```

## ORM Hate

Typical-Rest do not encourage Objection-Relation-Mapping/ORM ([ORM Hate](https://martinfowler.com/bliki/OrmHate.html))

## Shared Library

The [`pkg`](pkg) contain useful library for general needs

- [`pkg/typrest`](pkg/typrest) Utility for rest application e.g. health check, error, etc
- [`pkg/dbkit`](pkg/dbkit) Utility for database operation
- [`pkg/dbtxn`](pkg/dbtxn) Utility for database transaction
- [`pkg/dockerrx`](pkg/dockerrx) Docker Recipe Collection to be generated by typical-go
- [`pkg/echotest`](pkg/echotest) Utility for table-driven test echo handler

## Tools

- [`tools/typical-build`](tools/typical-build) Build-Tool for the project [learn more](https://github.com/typical-go/typical-go)
- [`tools/pg-tool`](tools/pg-tool) Tool to prepare the database [learn more](tools/pg-tool/README.md)


## 3rd-Party Library

The project [use go modules](https://blog.golang.org/using-go-modules) to manage package dependency. The compelete library list can be found in [go.mod](go.mod). 

- [Echo Framework](https://echo.labstack.com/): Web Framework
- [dig](https://github.com/uber-go/dig): Dependency Injection
- [gomock](https://github.com/golang/mock): Mock for interface
- [logrus](https://github.com/sirupsen/logrus): Logging
- [testify](https://github.com/stretchr/testify): Test assertion
- [squirrel](https://github.com/Masterminds/squirrel): SQL Query Builder
- [go-sqlmock](https://github.com/DATA-DOG/go-sqlmock): Mock for database connection 

## Golang References

- [Go Documentation](https://golang.org/doc/)
- [Go For Industrial Programming](https://peter.bourgon.org/go-for-industrial-programming/)
- [Uber Go Style Guide](https://github.com/uber-go/guide)
- [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)

## Project Status

The project status is `WIP` (Work in progress) which means the author continously evaluate and improve the project.

## License

This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details