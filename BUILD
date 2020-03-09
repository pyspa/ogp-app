load("@io_bazel_rules_docker//container:container.bzl", "container_image")

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
    name = "ogp_app_container",
    base = "@distroless_base_debian10//image",
    directory = "/app",
    files = [
        ":ogpapp_files",
    ],
    cmd = ["/app/ogp-app", "-c", "prd.toml"],
)