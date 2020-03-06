FROM busybox as caddy
COPY ./caddy /caddy
RUN chmod +x /caddy

FROM gcr.io/distroless/base-debian10
WORKDIR /app
COPY --from=caddy /caddy /app/caddy
COPY ./exec/exec /app/exec
COPY ./exec/Caddyfile /app/Caddyfile
COPY ./ogp-app /app/ogp-app
COPY ./config/prd.toml /app/prd.toml
COPY ./Koruri-Bold.ttf /app/Koruri-Bold.ttf
EXPOSE 9000
CMD ["/app/exec"]
