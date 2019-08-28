# auto-deploy

Auto-deploy is a web application written in Go that facilitates deployment via git web hooks.

Given a configuration file that contains project root folders and scripts, auto-deploy creates custom urls (web hooks). These web hooks can be configured to be called by any git repository management system (github, gitlab, etc). Upon invocation of the web hook, the projects deployment script is executed (as defined in the configuration).

Note that the deployment script is something that is provided by you and is thus not restricted to doing deployment, it can do whatever you'd like it to do. Auto-deploy can be thought of as a poor man's continuous deployment (CD) tool and was, in fact created in order to not have to run a heavy CI/CD installation in order to be able to use web hooks.

See the installation section for an example configuration and deployment script.

## Installation

Installation is as simple as compiling the binary (`auto-deploy.exec`) by running `make build` in the root folder. Thereafter either copy `config-sample.json` to `config.json` and update the contents as required before executing the binary.

### `sample-config.json`

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

Alternatively, if you have access to a VPS with a domain you can [roll your own ngrok in 15 minutes](https://zach.codes/roll-your-own-ngrok/). The following is a summary of the article using nginx running on a debian-based linux machine:

```nginx
upstream tunnel {
  server 127.0.0.1:8888;
}

server {
  server_name proxy.zach.codes;
  
  location / {
    proxy_set_header  X-Real-IP  $remote_addr;
    proxy_set_header  X-Forwarded-For $proxy_add_x_forwarded_for;
    proxy_set_header  Host $http_host;
    proxy_redirect off;

    proxy_pass http://tunnel;
  }
}
```

Create an `A` record for `proxy.zach.codes` to point to your machine and then run `sudo certbot --nginx -d proxy.zach.codes` to get https setup.

Finally run `ssh -R 8888:localhost:9111 root@proxy.zach.codes` where `9111` is the local service port. Your local service should now be accessible via `https://proxy.zach.codes`.

### Managing dependencies

Dependencies are stored under `$GOPATH/pkg/mod`.

- `go mod tidy`
- `go mod vendor`
