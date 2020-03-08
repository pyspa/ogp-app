FROM golang:1.14-buster as app
COPY . ./src/ogp-app
RUN cd src/ogp-app && go get .

FROM node:13.10.1-stretch-slim as assets
COPY client/ client/
RUN cd client && npm i && npm run build

FROM gcr.io/distroless/base-debian10
# FROM debian:buster-slim
WORKDIR /app
COPY --from=app /go/bin/ogp-app .
COPY --from=app /go/src/ogp-app/config/ ./
COPY --from=assets client/dist/ client/dist/
COPY Koruri-Bold.ttf .
CMD ["/app/ogp-app", "-c", "prd.toml"]
