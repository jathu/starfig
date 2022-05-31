#! /bin/bash

set -eu

run_coverage=$(echo "$@" | jq '.coverage.value')
open_coverage=$(echo "$@" | jq '."open".value')

root_dir=$(git rev-parse --show-toplevel)
test_path="$root_dir/..."

if [[ "$run_coverage" = "true" ]]; then
    out_file="/tmp/starfig.cover.out"
    go test -coverprofile="$out_file" "$test_path"
    html_file="/tmp/starfig.cover.html"

    if [[ "$open_coverage" = "true" ]]; then
        go tool cover -html="$out_file" -o "$html_file"
        open "$html_file"
    fi
else
    exec go test "$test_path"
fi
