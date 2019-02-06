<!-- FIXME: Project Title & Project description -->
# [Project Title]

[One Paragraph of project description goes here]

## Usages

### Configuration

### Usages

## Getting Started

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes. See deployment for notes on how to deploy the project on a live system.

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

### Installing

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
3. Download all dependencies
  ```sh
  dep ensure
  ```
4. Build and run the project
  ```sh
  # rename `my-project` to your project binary name
  go build -o my-project
  ./my-project
  ```


## Running the tests

Explain how to run the automated tests for this system

### Break down into end to end tests

Explain what these tests test and why

```
Give an example
```

### And coding style tests

Explain what these tests test and why

```
Give an example
```

## Deployment

Notes about how to deploy this on a live system

## Authors

* **Iman Tunggono** - *Initial work* - [imantung](https://github.com/imantung)

See also the list of [contributors](https://github.com/your/project/contributors) who participated in this project.

## License

This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details

## Acknowledgments

<!-- FIXME: Add some acknowledgements or remove the section -->
