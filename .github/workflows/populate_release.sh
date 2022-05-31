#! /bin/bash

set -eu

root_dir=$(git rev-parse --show-toplevel)
starfig="$root_dir/.build/starfig"
go build -o "$starfig"

if [ ! -f "$starfig" ]; then
    echo "$starfig is not built."
    exit 1
else
    changelog=$("$starfig" build //tool/changelog:changelog)
    changes=$(echo "$changelog" | jq '."//tool/changelog:changelog".changes')
    current_change=$(echo "$changes" | jq .[0])

    major_version=$(echo "$current_change" | jq .version.major)
    minor_version=$(echo "$current_change" | jq .version.minor)
    version="$major_version.$minor_version"

    echo -e "STARFIG_VERSION=${version}" >> $GITHUB_ENV
    echo -e "starfig version: ${version}"

    log=$(echo "$current_change" | jq -r .log)

    echo "STARFIG_CHANGELOG<<EOF" >> $GITHUB_ENV
    echo "$log" >> $GITHUB_ENV
    echo "EOF" >> $GITHUB_ENV

    echo -e "starfig changelog: ${log}"
fi
