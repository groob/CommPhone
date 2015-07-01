#!/bin/bash

VERSION="0.0.1"

echo "Building commphone version $VERSION"

mkdir -p pkg

build() {
  EXT=
  echo -n "=> $1-$2: "
  if [ $1 == "windows" ]; then
	  EXT=".exe"
  fi
  GOOS=$1 GOARCH=$2 go build -o pkg/commphone-$1-$2"${EXT}" -ldflags "-X main.Version $VERSION" 
  du -h pkg/commphone-$1-$2"${EXT}"
}

build "darwin" "amd64"
build "linux" "amd64"
build "windows" "amd64"
