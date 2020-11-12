# Apothiki /apoθε:cε:/

This is a very simple docker repository supporting only the Docker Registry HTTP API V2. There is no
authentication layer implemented right now so it can be used behind an HTTP Proxy that will provide that.

Configuration is done with the use of a config.yaml file you can find an example in the root of this repo. There is
also an admin interface running in port 8088.

## Docker

Starting the server with the docker image is very easy:

```bash
$ docker run -p 8080:8080 -v $(pwd)/data:/opt/apothiki/data kdihalas/apothiki:latest
```

## Todo
- Cache mode: use a list of upstream servers and cache layers locally for a period of time
- Mirror mode: replicate one or more upstream repositories
- Tests
- Proper error messages
- Authorization with webhooks