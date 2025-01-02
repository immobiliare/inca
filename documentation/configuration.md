# Configuration

To ensure proper functionality, this project relies on configuration specified in a YAML file. Follow the steps below to configure the project appropriately.

The configuration document is divided based on four main keys: `acl`, `providers`, `sentry`, `storage`.

## ACL

ACLs are mappings that, starting from the definition of the authentication token, specify an array of regular expressions. Through the mechanism of pattern matching, it is possible to check on which zones a specific token has permission to operate.

```
acl:
    04490F74210AFF3D0A49FB6280E11E09:
    - .*
    64EB0010394F9CD5C1F2C6AF452C8D52:
    - ^domain.dev.zone$
    - ^domain.stage.zone$
    - ^domain.zone$
```

# Providers

Local and letsencrypt are the two currently implemented providers. Local should be used to manage certificates for internal certificate authorities, while letsencrypt for certificates that need to be in the public WebPKI.

This is the expected structure of the local provider:

```
providers:
-   crt: /etc/inca.d/<zone>/crt.pem
    key: /etc/inca.d/<zone>/key.pem
    type: local
```

This is the expected structure of the letsencrypt provider:

```
providers:
-   ca: https://letsencrypt.org/certs/lets-encrypt-r3.pem
    email: <...>
    key: /etc/inca.d/letsencrypt.org/key.pem
    targets:
    -   challenge:
            aws_access_key_id: <...>
            aws_assume_role_arn: <...>
            aws_hosted_zone_id: <...>
            aws_profile: <...>
            aws_region: <...>
            aws_sdk_load_config: <...>
            aws_secret_access_key: <...>
            id: route53
        domain: <zone>
    type: letsencrypt
    -   challenge:
            id: webroot
        domain: <zone>
    type: letsencrypt
```

While using the webroot challenge (an http-01 challenge), it is the administrators' responsibility to route HTTP traffic for the prefix `/.well-known/acme-challenge/` to Inca.
Any challenge ID other than webroot is interpreted as a provider capable of handling a dns-01 challenge.

The variables required for configuring the DNS provider are those of Lego, the library used by Inca: [https://go-acme.github.io/lego/dns/](https://go-acme.github.io/lego/dns/).

## Sentry

It is possible to associate a string specifying a DSN with the sentry key.

```
sentry: https://<project_key>@<sentry_address>/<org_id>
```

## Storage

Filesystem and s3 are the two currently implemented interfaces. This is the persistence mechanism for the certificates.

This is the expected structure of the s3 storage:

```
storage:
    access: <...>
    endpoint: <...>
    region: <...>
    secret: <...>
    type: s3
```

This is the expected structure of the fs storage:

```
storage:
    path: <...>
    type: fs
```

This is the expected structure of the postgresql storage:

```
storage:
    host: <...>
    port: <...>
    user: <...>
    password: <...>
    dbname: <...>
    type: postgresql
```
