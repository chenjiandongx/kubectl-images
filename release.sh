#!/usr/bin/env bash

set -e

if [ "$#" -ne 1 ]; then
  echo "Usage: ./release.sh <version>"
  exit 1
fi

version=$1
outdir=releases/$version
bin=kubectl-images

linux_amd64_dist=kubectl-images_linux_amd64
linux_arm_dist=kubectl-images_linux_arm
linux_arm64_dist=kubectl-images_linux_arm64
darwin_amd64_dist=kubectl-images_darwin_amd64
darwin_arm64_dist=kubectl-images_darwin_arm64
windows_dist=kubectl-images_windows_amd64

echo $outdir

if [ ! -d $outdir ]; then
  mkdir -p $outdir
fi

cd cmd
GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o ../$outdir/$linux_amd64_dist
GOOS=linux GOARCH=arm go build -ldflags="-s -w" -o ../$outdir/$linux_arm_dist
GOOS=linux GOARCH=arm64 go build -ldflags="-s -w" -o ../$outdir/$linux_arm64_dist
GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o ../$outdir/$darwin_amd64_dist
GOOS=darwin GOARCH=arm64 go build -ldflags="-s -w" -o ../$outdir/$darwin_arm64_dist
GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o ../$outdir/$windows_dist
cd ..

cp LICENSE $outdir
cd $outdir
cp $linux_amd64_dist $bin && tar cfz $linux_amd64_dist.tar.gz LICENSE $bin
cp $linux_arm_dist $bin && tar cfz $linux_arm_dist.tar.gz LICENSE $bin
cp $linux_arm64_dist $bin && tar cfz $linux_arm64_dist.tar.gz LICENSE $bin
cp $darwin_amd64_dist $bin && tar cfz $darwin_amd64_dist.tar.gz LICENSE $bin
cp $darwin_arm64_dist $bin && tar cfz $darwin_arm64_dist.tar.gz LICENSE $bin
cp $windows_dist $bin && tar cfz $windows_dist.tar.gz LICENSE $bin
rm $bin

echo

linux_amd64_hash=$(echo `sha256sum $linux_amd64_dist.tar.gz` | awk '{print $1}')
linux_arm_hash=$(echo `sha256sum $linux_arm_dist.tar.gz` | awk '{print $1}')
linux_arm64_hash=$(echo `sha256sum $linux_arm64_dist.tar.gz` | awk '{print $1}')
darwin_amd64_hash=$(echo `sha256sum $darwin_amd64_dist.tar.gz` | awk '{print $1}')
darwin_arm64_hash=$(echo `sha256sum $darwin_arm64_dist.tar.gz` | awk '{print $1}')
windows_hash=$(echo `sha256sum $windows_dist.tar.gz` | awk '{print $1}')

cat <<EOF
apiVersion: krew.googlecontainertools.github.com/v1alpha2
kind: Plugin
metadata:
  name: images
spec:
  version: $1
  homepage: https://github.com/chenjiandongx/kubectl-images
  shortDescription: Show container images used in the cluster.
  description: |
    This plugin shows container images used in the Kubernetes cluster in a
    table view. You can show all images or show images used in a specified
    namespace.
  platforms:
  - selector:
      matchLabels:
        os: darwin
        arch: amd64
    files:
      - from: kubectl-images
        to: .
      - from: LICENSE
        to: .
    uri: https://github.com/chenjiandongx/kubectl-images/releases/download/$1/kubectl-images_darwin_amd64.tar.gz
    sha256: $darwin_amd64_hash
    bin: kubectl-images
  - selector:
      matchLabels:
        os: darwin
        arch: arm64
    files:
      - from: kubectl-images
        to: .
      - from: LICENSE
        to: .
    uri: https://github.com/chenjiandongx/kubectl-images/releases/download/$1/kubectl-images_darwin_arm64.tar.gz
    sha256: $darwin_arm64_hsah
    bin: kubectl-images
  - selector:
      matchLabels:
        os: linux
        arch: amd64
    files:
      - from: kubectl-images
        to: .
      - from: LICENSE
        to: .
    uri: https://github.com/chenjiandongx/kubectl-images/releases/download/$1/kubectl-images_linux_amd64.tar.gz
    sha256: $linux_amd64_hash
    bin: kubectl-images
  - selector:
      matchLabels:
        os: linux
        arch: arm64
    files:
      - from: kubectl-images
        to: .
      - from: LICENSE
        to: .
    uri: https://github.com/chenjiandongx/kubectl-images/releases/download/$1/kubectl-images_linux_arm64.tar.gz
    sha256: $linux_arm64_hash
    bin: kubectl-images
  - selector:
      matchLabels:
        os: linux
        arch: arm
    files:
      - from: kubectl-images
        to: .
      - from: LICENSE
        to: .
    uri: https://github.com/chenjiandongx/kubectl-images/releases/download/$1/kubectl-images_linux_arm.tar.gz
    sha256: $linux_arm_hash
    bin: kubectl-images
  - selector:
      matchLabels:
        os: windows
        arch: amd64
    files:
      - from: kubectl-images
        to: .
      - from: LICENSE
        to: .
    uri: https://github.com/chenjiandongx/kubectl-images/releases/download/$1/kubectl-images_windows_amd64.tar.gz
    sha256: $windows_hash
    bin: kubectl-images
EOF
