![Name](http://gitlab.rete.farm/pepita/guideline/docs/raw/master/ReadmeRepository/images/immobiliare-labs.png)

# InCA

<Description here>

## Table of Contents

- [Compatibility](#compatibility)
- [Install](#install)
- [Usage](#usage)
- [Requirements](#requirements)
- [Issues](#issues)
- [Contributing](#contributing)
- [Changelog](#changelog)
- [Reference team](#reference-team)
- [Credits](#credits)

## Compatibility

| Version | Status     | Go compatibility  |
|---      |---         |---                |
| latest  | maintained | >=1.18            |

## Install

### Build

Run `make build`

#### Build docker image

The docker image is built using mod vendoring (go >= 1.11 required) and multi-stage.

<insert instructions to build docker image here, es. Run `docker build -t registry.ekbl.net/sre/prunum:latest .`>

### Run

Run `make run` or manually with `./<insert command here>`
`-h` for inline help

#### Docker

Just run the docker image (see above): `docker run -v ${PWD}/configuration.yml:/app/configuration.yml -p 8080:8080 registry.ekbl.net/sre/prunum:latest`

Follow the above instructions for envvars or command line options.

## Usage

After running the http server (see below) send POST request like:

```
curl -v \
  -H "Content-Type: application/json" \
  -X POST \
  -d '{"host":"www.immobiliare.it","path":"/path/to/purge","layer":"front", "dry-run":true}' \
  http://localhost:8081/clear
```

The `layer` key can assume the following values:

- `front`: clear resource only on front servers
- `back`: clear resource only on back servers

**VERY IMPORTANT: the host/path MUST be relative to the cache layer. In other words assets to be purged on the front layer MUST be specified with the CDN host/path while assets to be purged on the back layer must be specidied with the mediaserver host/path**.

Eg.

- Correct :`'{"host": "pic.indomio.es", "path": "/image/1016776006/xxl.jpg", "layer": "front" ... }`
- Correct: `'{"host": "media.indomio.es", "path": "/image/1176892371/m-c.jpg", "layer": "back" ... }`
- **WRONG**: `'{"host": "media.indomio.es", "path": "/image/1176892371/m-c.jpg", "layer": "front" ... }`
  
See the `configuration.yml.example` file for server pools definition.

An example response could be:

### Configuration

Configuration file, listening port and address and server's log level can be passed both as flag or as environment vars (flag options overrides envvar). See the `-h` option for more info.

Configuration for ATS back/front servers and other stuff (to be implemented) is written on the `configuration.yml` file.

Configuration can be hot-reloaded sending SIGHUP to the main process.

### Response

Given the fact this software cannot know if the resource is actually cached or not into ATS servers, it always returns OK status code (200) immediately. The iteration over the ATS servers cluster is done asynchronously in the background.

## Requirements

- Go >= 1.15

## Issues

Open issue on this repository.

- [X] Add concurrency
  - [X] Add option to select if use concurrency or not
- [X] Add authentication (IP acl)
- [X] Better request parsing (return error on malformed requests)
- [X] Better doc
- [ ] Unit testing
- [ ] Add middleware for request limit (deny massive cache purge)
- [ ] Add metrics (send to statsd)
- [X] Automatic configuration reload on SIGHUP
- [X] Unified versioning for application/git tags/docker tags at build time

## Contributing

See [contributing](./CONTRIBUTING.md).


## Changelog

See [changelog](./CHANGELOG.md).

## Reference team

Immobiliare SRE

## Credits

Acknowledgment, repository links that were inspiring, etc. etc.
