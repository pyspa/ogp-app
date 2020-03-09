load("@io_bazel_rules_docker//container:container.bzl",
    "container_image", "container_push", "container_layer")
load("@bazel_tools//tools/build_defs/pkg:pkg.bzl", "pkg_tar")

filegroup(
    name = "ogpapp_files",
    srcs = [
        "ogp-app",
        "Koruri-Bold.ttf",
        "config/prd.toml",
    ],
)

filegroup(
    name = "ogpapp_client_js",
    srcs = glob([
        "client/dist/js/*.js",
        "client/dist/js/*.js.map",
    ]),
)

filegroup(
    name = "ogpapp_client_css",
    srcs = glob(["client/dist/css/*.css"]),
)

filegroup(
    name = "ogpapp_client_img",
    srcs = glob([
        "client/dist/img/*",
        "client/dist/img/icons/*.png",
    ]),
)

pkg_tar(
    name = "ogpapp_client",
    strip_prefix = "client/dist",
    srcs = glob(["client/dist/*"]) + [
        ":ogpapp_client_js",
        ":ogpapp_client_css",
        ":ogpapp_client_img",
    ],
    package_dir = "/client/dist",
)

container_layer(
    name = "ogpapp_client_layer",
    files = glob(["client/dist/*"]),
    directory = "/app/client/dist",
)

container_image(
    name = "ogpapp_container",
    base = "@distroless_base_debian10_debug//image",
    tars = ["ogpapp_client"],
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