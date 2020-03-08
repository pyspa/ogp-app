# ogp-app

ogp app generates `ogp:image` with the text user assigned.

## Prerequisite for development

* Go 1.13+ (for [go modules](https://github.com/golang/go/wiki/Modules))

TODO: To be documented

## Production Environment

Currently, [ogp.app](https://ogp.app) is running on Google Compute Engine.

## Run docker image

`$EXTERNAL_ADDR` is exposed host name and port for the ingress to the proxy. (eg. `localhost:9000`, `foo.com:443`)