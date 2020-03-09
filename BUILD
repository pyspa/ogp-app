load("@io_bazel_rules_docker//container:container.bzl",
    "container_image", "container_push", "container_layer")

filegroup(
    name = "ogpapp_files",
    srcs = [
        "ogp-app",
        "Koruri-Bold.ttf",
        "config/prd.toml",
    ],
)

container_layer(
    name = "ogpapp_client_layer",
    files = glob(["client/dist/*"]),
    directory = "/app/client/dist",
)

container_image(
    name = "ogpapp_container",
    base = "@distroless_base_debian10//image",
    layers = ["ogpapp_client_layer"],
    directory = "/app",
    workdir = "/app",
    files = [
        ":ogpapp_files",
    ],
    cmd = ["/app/ogp-app", "-c", "/app/prd.toml"],
)

container_push(
    name = "ogpapp_push",
    format = "Docker",
    image = ":ogpapp_container",
    registry = "gcr.io",
    repository = "pyspa-bot/ogp-app",
    tag = "bazel",
)