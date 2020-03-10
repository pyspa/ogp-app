load("@io_bazel_rules_go//go:def.bzl", "go_binary", "go_library", "go_test")
load(
    "@io_bazel_rules_docker//container:container.bzl",
    "container_image",
    "container_push",
    "container_layer",
)
load("@bazel_tools//tools/build_defs/pkg:pkg.bzl", "pkg_tar")
load("@bazel_gazelle//:def.bzl", "gazelle")

# gazelle:prefix github.com/pyspa/ogp-app
gazelle(name = "gazelle")

go_library(
    name = "go_default_library",
    srcs = [
        "app.go",
        "main.go",
        "middleware.go",
        "observability.go",
    ],
    importpath = "github.com/pyspa/ogp-app",
    visibility = ["//visibility:private"],
    deps = [
        "@com_github_burntsushi_toml//:go_default_library",
        "@com_github_golang_freetype//:go_default_library",
        "@com_github_golang_freetype//truetype:go_default_library",
        "@com_github_google_uuid//:go_default_library",
        "@com_github_gorilla_mux//:go_default_library",
        "@com_github_rs_zerolog//:go_default_library",
        "@com_github_rs_zerolog//log:go_default_library",
        "@com_google_cloud_go//compute/metadata:go_default_library",
        "@com_google_cloud_go//profiler:go_default_library",
        "@in_gopkg_natefinch_lumberjack_v2//:go_default_library",
        "@org_golang_x_image//font:go_default_library",
        "@org_golang_x_image//math/fixed:go_default_library",
    ],
)

go_binary(
    name = "ogpapp_binary",
    embed = [":go_default_library"],
    cgo = False,
    out = "ogp-app",
    visibility = ["//visibility:public"],
)

go_test(
    name = "go_default_test",
    srcs = ["app_test.go"],
    embed = [":go_default_library"],
)

filegroup(
    name = "ogpapp_files",
    srcs = [
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

container_image(
    name = "ogpapp_container",
    base = "@distroless_base_debian10_debug//image",
    tars = ["ogpapp_client"],
    directory = "/app",
    workdir = "/app",
    files = [
        ":ogpapp_binary",
        ":ogpapp_files",
    ],
    cmd = [
        "/app/ogp-app",
        "-c",
        "/app/prd.toml",
    ],
)

container_push(
    name = "ogpapp_push",
    format = "Docker",
    image = ":ogpapp_container",
    registry = "gcr.io",
    repository = "pyspa-bot/ogp-app",
    tag = "bazel",
)
