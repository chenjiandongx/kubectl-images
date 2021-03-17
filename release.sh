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
GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o ../$outdir/$darwin_amd64_dist
GOOS=darwin GOARCH=arm64 go build -ldflags="-s -w" -o ../$outdir/$darwin_arm64_dist
GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o ../$outdir/$windows_dist
cd ..

cp LICENSE $outdir
cd $outdir
cp $linux_amd64_dist $bin && tar cfz $linux_amd64_dist.tar.gz LICENSE $bin
cp $linux_arm_dist $bin && tar cfz $linux_arm_dist.tar.gz LICENSE $bin
cp $darwin_amd64_dist $bin && tar cfz $darwin_amd64_dist.tar.gz LICENSE $bin
cp $darwin_arm64_dist $bin && tar cfz $darwin_arm64_dist.tar.gz LICENSE $bin
cp $windows_dist $bin && tar cfz $windows_dist.tar.gz LICENSE $bin
rm $bin

echo "Please update this file accordingly:"
echo "https://github.com/kubernetes-sigs/krew-index/blob/master/plugins/images.yaml"

echo
echo "VERSION: $version"

echo
echo "SHA256 SUMS:"
echo "------------"
sha256sum $linux_amd64_dist.tar.gz
sha256sum $linux_arm_dist.tar.gz
sha256sum $darwin_amd64_dist.tar.gz
sha256sum $darwin_arm64_dist.tar.gz
sha256sum $windows_dist.tar.gz
