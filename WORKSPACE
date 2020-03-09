load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")

# Download the rules_docker repository at release v0.14.1
http_archive(
    name = "io_bazel_rules_docker",
    sha256 = "dc97fccceacd4c6be14e800b2a00693d5e8d07f69ee187babfd04a80a9f8e250",
    strip_prefix = "rules_docker-0.14.1",
    urls = ["https://github.com/bazelbuild/rules_docker/releases/download/v0.14.1/rules_docker-v0.14.1.tar.gz"],
)

load(
    "@io_bazel_rules_docker//repositories:repositories.bzl",
    container_repositories = "repositories",
)

container_repositories()

# This is NOT needed when going through the language lang_image
# "repositories" function(s).
load("@io_bazel_rules_docker//repositories:deps.bzl", container_deps = "deps")

container_deps()

load("@io_bazel_rules_docker//container:container.bzl", "container_pull")

container_pull(
  name = "distroless_base_debian10",
  registry = "gcr.io",
  repository = "distroless/base-debian10",
  # 'tag' is also supported, but digest is encouraged for reproducibility.
  # Find the SHA256 digest value from the detials page of prebuilt containers.
  # https://console.cloud.google.com/gcr/images/distroless/GLOBAL/base-debian10
  digest = "sha256:732acc54362badaa64d9c01619020cf96ce240b97e2f1390d2a44cc22b9ba6a3",
)

# for debug
container_pull(
  name = "distroless_base_debian10_debug",
  registry = "gcr.io",
  repository = "distroless/base-debian10",
  tag = "debug",
  # 'tag' is also supported, but digest is encouraged for reproducibility.
  # Find the SHA256 digest value from the detials page of prebuilt containers.
  # https://console.cloud.google.com/gcr/images/distroless/GLOBAL/base-debian10
  digest = "sha256:8ca4526452afe5d03f53c41c76c4ddb079734eb99913aff7069bfd0d72457726",
)
