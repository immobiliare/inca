![Inca](http://gitlab.rete.farm/pepita/guideline/docs/raw/master/.backstage/docs/ReadmeRepository/images/immobiliare-labs.png)

# Inca

Inca is an INternal CA manager for local CAs as well as external ones.

## Table of Contents

- [Inca](#inca)
  - [Table of Contents](#table-of-contents)
  - [Compatibility](#compatibility)
  - [Install](#install)
  - [Usage](#usage)
    - [Bootstrap](#bootstrap)
    - [Generate certificates](#generate-certificates)

## Compatibility

| Version | Status     | Go compatibility |
| ------- | ---------- | ---------------- |
| latest  | maintained | >= 1.18          |

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
  registry.ekbl.net/sistemi/inca:latest
```

## Usage

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
