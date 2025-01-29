# Session Server

This is a web server that is a front end to the session deployment system. The public instance is hosted on
`https://session.luhack.uk`.

## Building

To build the server, you will need to have the following installed:

- go
- make
- tar

You should be able to just clone this repo and run `make` to build the server. This will give you a `release.tar.gz`
that you can extract anywhere and then use as below.

## Usage

The binary is called `session-server` and requires a configuration file in the same directory as the binary called
`config.yml`.

```yaml
server:
  host: "localhost:8080"

  domain: "session.luhack.uk"
  protocol: "https"

session:
  title: "Demo Session"
  backendMap: "backend-map.yml"

security:
  jwtSecret: "change_me"
  server: "https://auth.luhack.uk"
```

- The `host` field is the address that the server will listen on.
- The `domain` field is the domain that the public domain that users will access the server on.
- The `protocol` field is the protocol that public users will access the server on.
- The `title` field is the title of the session that will be displayed on the front page.
- The `backendMap` field is the path to the backend map file, this expects the standard format of the backend map file.
- The `jwtSecret` field is the secret that the server will use to sign JWTs.
- The `server` field is the address of the auth server that the server will use to verify JWTs.
 
