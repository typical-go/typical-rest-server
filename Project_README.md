<!-- FIXME: Project Title & Project description -->
# [Project Title]

[One Paragraph of project description goes here]

## Usages

List of Commands
```
serve, s              Serve the clients
database, db          database administration
config-details, conf  Print configuration detail
help, h               Shows a list of commands or help for one command
```

## Configuration


```
+-------------+--------+----------+-----------+
|    NAME     |  TYPE  | REQUIRED |  DEFAULT  |
+-------------+--------+----------+-----------+
| ADDRESS     | string | true     |           |
| DB_NAME     | string | true     |           |
| DB_USER     | string | true     |           |
| DB_PASSWORD | string | true     |           |
| DB_HOST     | string |          | localhost |
| DB_PORT     | int    |          |      5432 |
+-------------+--------+----------+-----------+
```

## Getting Started

### Prerequisites

What things you need to install the software and how to install them
1. Install [go](https://golang.org/) and [set up `GOPATH`](https://github.com/golang/go/wiki/SettingGOPATH)
  ```sh
  # for macOS, required https://brew.sh/
  brew install go

  # create empty directory for GOPATH
  mkdir $HOME/go && cd $HOME/go && mkdir bin pkg src && cd -

  # set up GOPATH, the example using https://ohmyz.sh/
  # change ~/.zshrc to ~/.bashrc if you use default bash
  echo "export GOPATH=$HOME/go
  export PATH=$PATH:$GOPATH/bin
  " >> ~/.zshrc
  source ~/.zshrc
  ```
2. Install [dep](https://github.com/golang/dep) as dependency management tool
  ```sh
  # for macOS
  brew install dep
  ```
3. Install [direnv](https://direnv.net/) for directory level environment variable
  ```sh
  # for macOS
  brew install direnv
  ```

### Clone

<!-- FIXME: Project path, git path and project binary name -->
A step by step that tell you how to get a development environment running
1. Create project directory
  ```sh
  mkdir $GOPATH/src/your/project/path
  ```
2. Change directory to project path then `git clone` the project
  ```sh
  cd $GOPATH/src/your/project/path && git clone git@your/project
  ```

### Build

```sh
make build
```
