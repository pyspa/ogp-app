FROM busybox
RUN chmod +x /workspace/caddy

FROM gcr.io/distroless/base-debian10
WORKDIR /app
COPY /workspace/caddy /app/caddy
COPY /workspace/ogp-app /app/ogp-app
COPY /workspace/config/prd.toml /app/prd.toml
COPY /workspace/Koruri-Bold.ttf /app/Koruri-Bold.ttf
CMD ["/app/ogp-app", "-c", "/app/prd.toml", "&&", "caddy", "reverse-proxy", "--from", "ogp.app", "--to", "localhost:8080"]