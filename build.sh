#!/bin/sh

targets="linux/amd64 linux/arm64 linux/386 windows/amd64 darwin/amd64 darwin/arm64"

for target in $targets; do
    GOOS=$(echo $target | cut -d'/' -f1)
    GOARCH=$(echo $target | cut -d'/' -f2)
    output_name="plot-${GOOS}-${GOARCH}"
    if [ "$GOOS" = "windows" ]; then
        output_name="${output_name}.exe"
    fi
    echo "building for $GOOS/$GOARCH..."
    GOOS=$GOOS GOARCH=$GOARCH go build -o bin/$output_name
done
