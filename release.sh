#!/usr/bin/env bash

set -e

if [ "$#" -ne 1 ]; then
  echo "Usage: ./release.sh <version>"
  exit 1
fi

version=$1
outdir=releases/$version
bin=kubectl-images
linux_dist=kubectl-images_linux_amd64
darwin_dist=kubectl-images_darwin_amd64
windows_dist=kubectl-images_windows_amd64

echo $outdir

if [ ! -d $outdir ]; then
  mkdir -p $outdir
fi

cd cmd
GOOS=linux go build -ldflags="-s -w" -o ../$outdir/$linux_dist
GOOS=darwin go build -ldflags="-s -w" -o ../$outdir/$darwin_dist
GOOS=windows go build -ldflags="-s -w" -o ../$outdir/$windows_dist
cd ..

cp LICENSE $outdir
cd $outdir
cp $linux_dist $bin && tar cfz $linux_dist.tar.gz LICENSE $bin
cp $darwin_dist $bin && tar cfz $darwin_dist.tar.gz LICENSE $bin
cp $windows_dist $bin && tar cfz $windows_dist.tar.gz  LICENSE $bin
rm $bin

echo "Please update this file accordingly:"
echo "https://github.com/kubernetes-sigs/krew-index/blob/master/plugins/images.yaml"

echo
echo "VERSION: $version"

echo
echo "SHA256 SUMS:"
echo "------------"
sha256sum $linux_dist.tar.gz
sha256sum $darwin_dist.tar.gz
sha256sum $windows_dist.tar.gz
