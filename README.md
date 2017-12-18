# Cron server

Setup cron tasks to publish to google pub sub or request endpoints at specified
intervals

### Install:

See [releases](https://github.com/fixate/drone-secrets/releases/latest) or build repo using `make build` or `make builddev`

TODO:

- Use CI for github releases for multiple targets

### Usage:

```shell
cron-server -m manifest.yml
```

*manifest.yml*

See [https://github.com/fixate/cron-server/blob/master/examples/manifest.yml](https://github.com/fixate/cron-server/blob/master/examples/manifest.yml)

## TODO:

- Tests 
