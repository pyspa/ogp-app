load("@build_bazel_rules_nodejs//:index.bzl", "nodejs_binary")

nodejs_binary(
    name = "vue_cli_service",
    entry_point = "@npm//:node_modules/@vue/cli-service/bin/vue-cli-service.js",
    node_modules = "@npm//:node_modules",
)

filegroup(
    name = "ogpapp_client_src",
    srcs = glob([
        "src/**/*",
        "public/*"
    ]),
)

genrule(
    name = "ogpapp_client",
    srcs = [
        ":ogpapp_client_src",
    ],
    outs = ["dist"],
    cmd = "$(location vue_cli_service) build",
    tools = [":vue_cli_service"],
)