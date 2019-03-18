# Typical Go Server (WIP)

This project aim to provide practical template project for RESTful API Server.

- [x] Encourage [KISS Principle](https://en.wikipedia.org/wiki/KISS_principle)
  - [x] Go Idiomatic
  - [x] [Rails-like](https://guides.rubyonrails.org/getting_started.html#creating-the-blog-application) Project Layout
  - [x] [Separation of concerns](https://en.wikipedia.org/wiki/Separation_of_concerns)
  - [x] [Dependency Injection](https://stackoverflow.com/questions/130794/what-is-dependency-injection)
- [x] Micro-services Friendly
  - [x] Follow [12 Factor App](https://12factor.net/)
  - [x] Environment Variable as configuration (for development, use [direnv](https://direnv.net/))
  - [x] Package Management/Vendoring
  - [x] [Graceful Shutdown](https://12factor.net/disposability)
  - [ ] [Health Check](https://microservices.io/patterns/observability/health-check-api.html)
- [ ] Self-explanatory Project
  - [ ] [Self Testing Code](https://martinfowler.com/bliki/SelfTestingCode.html)
  - [x] [Project README](Project_README.md)
  - [x] Automatically generate cli/command documentation
  - [x] Automatically generate config documentation
  - [x] Makefile
  - [ ] GoDoc
  - [ ] Wiki
- [x] RESTful API
  - [x] CRUD Operation
  - [x] Model Validation
  - [ ] Authentication
  - [ ] CORS
  - [ ] Cache
  - [ ] Pagination
  - [ ] Search API
  - [ ] API Versioning
  - [ ] API Documentation
- [x] Working with Database
  - [x] Postgres Database
  - [x] Data Access Layer/Repository Pattern
  - [ ] Test database
  - [ ] Soft delete
  - [x] Database Administration
    - [x] Create db
    - [x] Drop db
    - [x] Migration
    - [x] Rollback
- [ ] Worker
  - [ ] Job Background Process
- [ ] Internationalization  
- [ ] Misc
  - [ ] Debug/Profiling
  - [ ] Travis CI example
  - [ ] Docker example

## Learn More

Check our [wiki](https://github.com/imantung/typical-go-server/wiki)

## Contributing

_Under Construction_
<!-- FIXME: -->

## License

This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details
