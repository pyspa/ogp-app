FROM busybox as caddy
COPY ./caddy /caddy
RUN chmod +x /caddy

FROM gcr.io/distroless/base-debian10
WORKDIR /app
COPY --from=caddy /caddy /app/caddy
COPY ./ogp-app /app/ogp-app
COPY ./config/prd.toml /app/prd.toml
COPY ./Koruri-Bold.ttf /app/Koruri-Bold.ttf
EXPOSE 2015
#CMD ["/app/ogp-app", "-c", "/app/prd.toml", "&&", "/app/caddy", "reverse-proxy", "--to", "localhost:8080"]
#CMD ["/app/caddy", "run"]
CMD ["/app/ogp-app", "-c", "/app/prd.toml"]