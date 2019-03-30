<!-- FIXME: Project Title & Project description -->
# [Project Title]

[One Paragraph of project description goes here]

## Getting Started

### Prerequisites

1. Install [go](https://golang.org/) and [set up `GOPATH`](https://github.com/golang/go/wiki/SettingGOPATH)
2. Install [dep](https://github.com/golang/dep) as dependency management tool
3. Install [direnv](https://direnv.net/) for directory level environment variable

### Clone

<!-- FIXME: Project path, git path and project binary name -->

1. Create project directory
  ```sh
  mkdir $GOPATH/src/your/project/path
  ```
2. Change directory to project path then `git clone` the project
  ```sh
  cd $GOPATH/src/your/project/path && git clone git@your/project
  ```

## Installing

Use `Make` to clean, build, test, generate mock class, etc.
- `make all`: Install dependencies and build the binary
- `make dep`: Install dependencies
- `make build`: Build the binary
- `make test`: Running test
- `make test-report`: Running test and show coverage profile
- `make clean`: Clean build files
- `make clean-all`: Clean build files and all dependency
- `make mock`: Generate mock class

## Usages

### Show a list of Commands
`[BINARY] help`
```
NAME:
   [Name] - API for [Usage]

USAGE:
   typical-go-server [global options] command [command options] [arguments...]

VERSION:
   0.0.1

COMMANDS:
     serve, s      Run the server
     database, db  Database Administration
     config, conf  Configuration
     help, h       Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h     show help
   --version, -v  print the version
```

### Run the server
`[BINARY] serve`

### Database Administration

- `[BINARY] db create`: Create New Database
- `[BINARY] db drop`: Drop Database
- `[BINARY] db migrate`: Migrate Database
- `[BINARY] db rollback`: Rollback Database

### Print configuration details

`[BINARY] config`
```
+-------------+--------+----------+-----------+
|    NAME     |  TYPE  | REQUIRED |  DEFAULT  |
+-------------+--------+----------+-----------+
| APP_ADDRESS | string | true     |           |
| DB_NAME     | string | true     |           |
| DB_USER     | string | true     |           |
| DB_PASSWORD | string | true     |           |
| DB_HOST     | string |          | localhost |
| DB_PORT     | int    |          |      5432 |
+-------------+--------+----------+-----------+
```
