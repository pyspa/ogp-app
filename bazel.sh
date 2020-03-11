#!/bin/bash

bazel run //:gazelle -- update-repos -from_file=go.mod
bazel run :ogpapp_push