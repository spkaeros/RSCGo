#!/bin/bash
CGO_ENABLED=0
LINK_FLAGS="-w -s -buildid=rscgo-production -installsuffix rscgo -v"
BUILD_FLAGS="-trimpath $PWD -smallframes -pack -buildid=rscgo-production -v -complete -p github.com/spkaeros/rscgo/ -nolocalimports -m"

unset GOARCH
for arch in `go tool dist list|sed -e 's/\// /g' |cut -f2 -d' ' |sort |uniq`; do
	if [[ $1 == $arch ]]; then
		export GOARCH=$1
		break
	fi
done
if [[ -z $GOARCH ]]; then
	echo "Could not find matching target CPU architecture.  Must provide target OS and architecture as arguments, e.g: ./build.sh linux amd64"
	listTargets
fi
unset GOOS
for os in `go tool dist list|sed -e 's/\// /g' |cut -f1 -d' ' |sort |uniq`; do
	if [[ $2 == $os ]]; then
		export GOOS=$2
		break
	fi
done
if [[ -z $GOOS ]]; then
	echo "Could not find matching target OS.  Must provide target OS and architecture as arguments, e.g: ./build.sh linux amd64"
	listTargets
	exit
fi
EXECUTABLE="./bin/game-$GOOS-$GOARCH"

if [[ $GOOS == 'windows' ]]; then
	EXECUTABLE="$EXECUTABLE.exe"
fi

if [[ $GOOS == 'linux' || $GOOS == 'freebsd' ]]; then
	echo 'Position-independent executable files supported on target OS; setting buildmode to pie in build flags...'
	LINK_FLAGS="-buildmode pie $LINK_FLAGS"
fi
go build -o="$EXECUTABLE" -gcflags="$BUILD_FLAGS" -ldflags="$LINK_FLAGS" pkg/game/server.go

exit

listTargets() {
	echo "Available targets (os/arch):"
	echo `go tool dist list`
	exit
}