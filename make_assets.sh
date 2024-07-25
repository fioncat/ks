#!/bin/bash

set -e

VERSION="$1"
COMMIT="$2"
BUILD_DATE="$3"

printf "Begin to build assets, version: %s, commit: %s, build_date: %s\n" "$VERSION" "$COMMIT" "$BUILD_DATE"

targets=( \
	"darwin_amd64" \
	"darwin_arm64" \
	"linux_amd64" \
	"linux_arm64" \
)

mkdir -p assets

for target in "${targets[@]}"; do
	echo "Build target: ${target}"
	IFS='_' read -r -a OS_ARCH <<< "$target"
	CGO_ENABLED=0 GOOS="${OS_ARCH[0]}" GOARCH="${OS_ARCH[1]}" go build -ldflags="-X 'main.Version=${VERSION}' -X 'main.Commit=${COMMIT}' -X 'main.BuildDate=${BUILD_DATE}'" -o bin/${target}/kser ./kser
	tar -czf assets/ks-${target}.tar.gz -C bin/${target} kser
done
