#!/usr/bin/env bash

set -e

if [ "$#" -ne 1 ]; then
  echo "Usage: ./release.sh <version>"
  exit 1
fi

version=$1
outdir=releases/$version
linux_dist=kubectl-images_linux_amd64
darwin_dist=kubectl-images_darwin_amd64
windows_dist=kubectl-images_windows_amd64

if [ -d $outdir ]; then
  mkdir -p $outdir
fi

cp LICENSE $outdir

cd cmd
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o ../$outdir/$linux_dist
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o ../$outdir/$darwin_dist
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o ../$outdir/$windows_dist
cd ..

tar cfz $outdir/$linux_dist.tar.gz -C $outdir --transform="flags=r;s|$linux_dist|kubectl-images|" $linux_dist LICENSE
tar cfz $outdir/$darwin_dist.tar.gz -C $outdir --transform="flags=r;s|$darwin_dist|kubectl-images|" $darwin_dist LICENSE
tar cfz $outdir/$windows_dist.tar.gz -C $outdir --transform="flags=r;s|$windows_dist|kubectl-images|" $windows_dist LICENSE

echo "Please update this file accordingly:"
echo "https://github.com/kubernetes-sigs/krew-index/blob/master/plugins/images.yaml"

echo
echo "VERSION: $version"

echo
echo "SHA256 SUMS:"
echo "------------"
sha256sum $outdir/$linux_dist.tar.gz
sha256sum $outdir/$darwin_dist.tar.gz
sha256sum $outdir/$windows_dist.tar.gz
