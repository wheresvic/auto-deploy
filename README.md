# auto-deploy

Auto-deploy is a web application written in Go that facilitates deployment via git web hooks.

Given a configuration file that contains git repository details, auto-deploy creates and listens on custom urls that can be called by the git repository management system (github, gitlab, etc). Upon invocation of the web hook, a script is executed (as defined in the configuration). See the installation section for an example configuration and deployment script.

## Installation

### `config.json`

```json

```

### sample `deploy.sh`

Note that this deployment script assumes that the project git repository is also cloned on the production server.

```sh

```

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

### Testing Webhooks

The current webhook api is available under `/api/webhooks/:project-slug/:scm-type`. 

Note you can use [serveo.net](https://serveo.net) to obtain a public url which can be used for testing webhooks, e.g.: `ssh -R 80:localhost:9111 serveo.net` where `9111` is the configured port number. The final url will then look something like: `https://optimus.serveo.net/api/webhooks/basic-blockchain/github`.

### Managing dependencies

Dependencies are stored under `$GOPATH/pkg/mod`.

- `go mod tidy`
- `go mod vendor`
