#! /bin/bash

set -eu

root_dir=$(git rev-parse --show-toplevel)
cd "$root_dir"
exec go fmt ./...
