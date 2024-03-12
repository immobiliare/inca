# Inca <a href="#inca"><img align="left" width="100px" src="https://github.com/immobiliare/inca/blob/main/.github/icon.png"></a>

[![pipeline status](https://github.com/immobiliare/inca/actions/workflows/test.yml/badge.svg)](https://github.com/immobiliare/inca/actions/workflows/test.yml)

Inca stands for INternal CA, and it's primary aim is self-explained: handling certificate-wise flows with regards to a local and private CA.
On the flip side, its ambitious vocation is to eliminate all the complexity on maintaining a PKI within a company.

### Alternatives

Alternatives to Inca that don't have all the required features:

- [OpenXPKI](https://www.openxpki.org/)
- [EJBCA](https://www.ejbca.org/)
- [step-ca](https://github.com/smallstep/certificates)
- [Locksmith](https://github.com/kenmoini/locksmith)
- [Certbot](https://certbot.eff.org/) - The recommended LetsEncrypt client
- [Lego](https://github.com/go-acme/lego) - Let's Encrypt client and ACME library written in Go

### Internal CA

Given a CA keypair, Inca exposes a set of endpoints usable to interact with the aforementioned CA to issue, revoke, extend valid certificates.

### Proxying to other CAs

If configured to do so, Inca can proxy the already mentioned requests to external providers (e.g. Let's Encrypt), providing a simple and common interface for certificates regardless of their origin.

### Storing certificates

Inca does not only issue certificates, it caches and stores them on a configurable storage (e.g. locally on filesystem, on S3), reusing them if asked to. 

### Foreign certificates

Through the webgui, Inca allows for certificates to be manually imported, if issued via a third-party flow.

### Self-explanatory API endpoints

Given it's minimal semantic, Inca is super easy to integrate into third-party tools, as obtaining a valid certificate is as easy as `curl https://inca.domain.tld/whatever-cn.domain.tld`.

[![Inca homepage](https://github.com/immobiliare/inca/blob/main/.github/sample-1.png)](#inca)

[![Inca detail](https://github.com/immobiliare/inca/blob/main/.github/sample-2.png)](#inca)

## Table of Contents

- [Install](#install)
- [Usage](#usage)
  - [Custom installation](#custom-installation)
  - [Generate certificates](#generate-certificates)
- [Changelog](#changelog)
- [Contributing](#contributing)
- [Documentation](#documentation)
- [Powered apps](#powered-apps)
- [Support](#support)

## Install

Either

```sh
go build
go install
inca --help
```

or

```sh
docker run -it -v --network host ${PWD}/inca.yml:/etc/inca:ro \
  ghcr.io/immobiliare/inca:latest
```

## Usage

If you're `curl`-ninja enough:

```sh
alias curl='curl -H "Authorization:Bearer REDACTED"'
# fetch certificate
curl https://inca.domain.tld/domain.tld.farm -o whatever.domain.tld.pem
# fetch certificate with further params
curl https://inca.domain.tld/whatever-with-details.domain.tld?alt=whatever2.domain.tld&duration=2y
# fetch key
curl https://inca.domain.tld/whatever.domain.tld/key -o whatever.domain.tld.key
# remove certificate
curl -X DELETE https://inca.domain.tld/whatever.domain.tld
```

Otherwise, just open Inca on a browser.

#### Custom installation

```sh
inca gen -n domain.tld -o /etc/inca.d
cat >/etc/inca <<EOF
bind: :80
providers:
  - type: local
    crt: /etc/inca.d/crt.pem
    key: /etc/inca.d/key.pem
storage:
  type: fs
  path: /etc/inca.d
acl:
  REDACTED:
    - ^nice.domain.tld$
    - .*.notsonice.domain.tld$
EOF
inca server
```

#### Generate certificates

```sh
curl -H "Authorization:Bearer REDACTED" http://localhost:80/crt.domain.tld -o crt.domain.tld.pem
curl -H "Authorization:Bearer REDACTED" http://localhost:80/crt.domain.tld/key -o crt.domain.tld.key
```

## Changelog

See [changelog](./CHANGELOG.md).

## Contributing

See [contributing](./CONTRIBUTING.md).

## Documentation

See [configuration](../documentation/configuration.md).

## Powered apps

Inca was created by ImmobiliareLabs, the technology department of [Immobiliare.it](https://www.immobiliare.it), the #1 real estate company in Italy.

**If you are using Inca [drop us a message](mailto:opensource@immobiliare.it)**.

## Support

Made with ❤️ by [ImmobiliareLabs](https://github.com/immobiliare) and all the [contributors](./CONTRIBUTING.md#contributors)

If you have any question on how to use Inca, bugs and enhancement please feel free to reach us out by opening a [GitHub Issue](https://github.com/immobiliare/inca/issues).
