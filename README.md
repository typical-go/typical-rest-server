# Typical Go Server (WIP)

This project aim to provide starter/example project for API Service in Go.

[KISS](https://en.wikipedia.org/wiki/KISS_principle) is the core principle while follows [The Twelve-Factor App](https://12factor.net/) to make sure the service will be easy to ship and adopt another best practices for microservices.
- [x] Go Idiomatic
- [x] [Rails-like](https://guides.rubyonrails.org/getting_started.html#creating-the-blog-application) Project Layout
- [x] Package Management/Vendoring
- [ ] [Self Testing Code](https://martinfowler.com/bliki/SelfTestingCode.html)
- [x] Separation of concern/Modular/[Dependency Injection](https://stackoverflow.com/questions/130794/what-is-dependency-injection)
- [x] Postgres Database
- [x] [Graceful Shutdown](https://12factor.net/disposability)
- [x] Database Administration (Create/Drop/Migrate/Rollback)
- [x] Data Access Layer/Repository Pattern
- [x] [README Documentation](Typical_README.md)
- [ ] [Health Check](https://microservices.io/patterns/observability/health-check-api.html)
- [ ] Test Coverage Checking
- [ ] Simple authentication
- [ ] Generated API Doc
- [ ] API Versioning
- [ ] Cache
- [ ] Worker
- [ ] Profiling
- [ ] Makefile
- [ ] Travis CI Release example
- [ ] Docker file


### How to Use

_Under Construction_
<!-- FIXME: -->

### Project Layout

_Under Construction_
<!-- FIXME: -->

### Library Overview
- [urfave/cli](https://github.com/urfave/cli): A simple, fast, and fun package for building command line apps in Go
- [labstack/echo](https://github.com/labstack/echo): High performance, minimalist Go web framework
- [kelseyhightower/envconfig](https://github.com/kelseyhightower/envconfig): Managing configuration data from environment variables
- [stretchr/testify](https://github.com/stretchr/testify): A toolkit with common assertions and mocks that plays nicely with the standard library
- [uber-go/dig](https://github.com/uber-go/dig): A reflection based dependency injection toolkit for Go.
- [lib/pq](https://github.com/lib/pq): Pure Go Postgres driver for database/sql
- [imantung/go-helper](https://github.com/imantung/go-helper): Helper library collection for golang
- [Masterminds/squirrel](https://github.com/Masterminds/squirrel): Fluent SQL generation for golang
- [DATA-DOG/go-sqlmock](https://github.com/DATA-DOG/go-sqlmock): Sql mock driver for golang to test database interactions


### Contributing

_Under Construction_
<!-- FIXME: -->

### Similar Projects

- [gobuffalo](https://gobuffalo.io/): Web Development eco-system
- [go-swagger](https://goswagger.io/): OpenAPI implementation on Go
- [go-bootstrap](http://go-bootstrap.io/): Generates a lean and mean Go web project.
- [qiangxue/golang-restful-starter-kit](github.com/qiangxue/golang-restful-starter-kit): A RESTful application boilerplate in Go (golang) taking best practices and utilizing best available packages and tools

### Authors

* **[imantung](https://github.com/imantung)** - *Initial work* -

See also the list of [contributors](https://github.com/your/project/contributors) who participated in this project.

## License

This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details
