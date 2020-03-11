load("@build_bazel_rules_nodejs//:index.bzl", "nodejs_binary", "npm_package_bin")
load("@npm//@vue/cli-service:index.bzl", "vue_cli_service")

#nodejs_binary(
#    name = "vue_cli_service",
#    entry_point = "@npm//:node_modules/@vue/cli-service/bin/vue-cli-service.js",
#    node_modules = "@npm//:node_modules",
#)

filegroup(
    name = "ogpapp_client_src",
    srcs = glob([
        "public/*",
        "src/assets/*",
        "src/components/*.vue",
        "src/router/*.ts",
        "src/store/*.ts",
        "src/views/*.vue",
        "src/*",
        "vue.config.js",
        "tsconfig.json",
        "babel.config.js",
        ".browserlistrc",
        ".eslintrc.js",
    ]),
)

# `npm_package_bin` enables to run arbitrary executables under `node_modules/.bin`
# ref. https://github.com/bazelbuild/rules_nodejs/blob/1.4.1/internal/node/npm_package_bin.bzl#L70-L132
vue_cli_service(
    name = "ogpapp_client_build",
    args = [
        "build",
    ],
    data = [
        ":ogpapp_client_src",
        "@npm//:node_modules",
    ],
    output_dir = True,
)