#! /bin/bash

set -eu

should_run=$(echo "$@" | jq '.run.value')
proxy_arguments=( $(echo "$@" | jq -rc '."proxy-arguments".value[]') )

root_dir=$(git rev-parse --show-toplevel)
destination="$root_dir/.build/starfig"

go build -o "$destination" "$root_dir"

if [[ "$should_run" = "true" ]]; then
    STARFIG_DEBUG=1 exec "$destination" "${proxy_arguments[@]}"
fi
