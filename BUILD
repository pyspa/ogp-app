load("@io_bazel_rules_docker//container:container.bzl", "container_image", "container_push")

filegroup(
    name = "ogpapp_files",
    srcs = [
        "ogp-app",
        "Koruri-Bold.ttf",
        "config/prd.toml",
    ],
    visibility = ["//visibility:public"],
)
container_image(
    name = "ogpapp_container",
    base = "@distroless_base_debian10//image",
    directory = "/app",
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