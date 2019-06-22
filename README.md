# auto-deploy

Auto-deploy is a web application written in Go that facilitates deployment via git web hooks. Given a configuration file that contains git repository details, auto-deploy creates and listens on custom urls that can be called by the git repository management system (github, gitlab, etc). Upon invocation of the web hook, a script is executed (as defined in the configuration), which can do anything but is generally used to update the project repository and run build steps.

## Installation

## Development

### Setup

Require Go `>=1.12.1` installed with GVM, also see [TODO](connecting to an oracle database using Go):

- `make deps` : get all dependencies, uses go modules to get exact dependencies.
- `make run` : compile and execute
- `make test` : run all tests
- `make test-unit` : run unit tests
- `make test-integration` : run integration tests

See `Makefile` for more options.

To get a development server up and running, we require ag (`sudo apt install silversearcher-ag`) and entr ([http://eradman.com/entrproject/](http://eradman.com/entrproject/)). Run `run-dev.sh` from the root folder to start a nodemon like process that restarts itself if any changes are detected to go code.

### Managing dependencies

Dependencies are stored under `$GOPATH/pkg/mod`.

- `go mod tidy`
- `go mod vendor`
