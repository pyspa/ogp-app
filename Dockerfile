# Dockerfile is for development use only.
# The container image runs in production should be built and delivered
# via Google Cloud Build. You can build your own with `cloud-build-local`
# command as follows (NOTE: Cloud Build local builder only works on Linux):
# $ cloud-build-local --config=cloudbuild.yaml --dryrun=false .
#
# Refer to the Cloud Build official docs for details:
# https://cloud.google.com/cloud-build/docs/build-debug-locally
FROM golang:1.14-buster as app
COPY . ./src/ogp-app
RUN cd src/ogp-app && go get .

FROM node:13.10.1-stretch-slim as assets
COPY client/ client/
RUN cd client && npm ci && npm run build

FROM gcr.io/distroless/base-debian10
# For debug use: you can enter shell with --entrypoint=sh.
# FROM gcr.io/distroless/base-debian10:debug
WORKDIR /app
COPY --from=app /go/bin/ogp-app .
COPY --from=app /go/src/ogp-app/config/ ./
COPY --from=assets client/dist/ client/dist/
COPY Koruri-Bold.ttf .
CMD ["/app/ogp-app", "-c", "/app/prd.toml"]
