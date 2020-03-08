FROM gcr.io/distroless/base-debian10
WORKDIR /app
COPY ./client/dist/ client/dist/
COPY ./ogp-app /app/ogp-app
COPY ./config/prd.toml /app/prd.toml
COPY ./Koruri-Bold.ttf /app/Koruri-Bold.ttf
CMD ["/app/ogp-app", "-c", "prd.toml"]
