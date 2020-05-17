#!/bin/bash
compile() {
	os=$1
	arch=$2
	LINK_FLAGS="-w -s -buildid=rscgo-production -v -extldflags=-static -X main.version=$tag"
	BUILD_FLAGS="-trimpath -smallframes -pack -buildid=rscgo-production -v -complete -nolocalimports -v"

	EXECUTABLE="./bin/game-$os-$arch"
	if [[ $os == 'windows' ]]; then
		EXECUTABLE="$EXECUTABLE.exe"
	fi
	
#	if [[ $arch == 'amd64' && $os == 'linux' || $arch == 'amd64' && $os == 'freebsd' ]]; then
#		echo 'Position-independent executable files supported on target OS; enabling...'
#		LINK_FLAGS="-buildmode pie $LINK_FLAGS"
#	fi
	CGO_ENABLED=0 CC=gcc GOOS=$os GOARCH=$arch go build -o "$EXECUTABLE" -gcflags="$BUILD_FLAGS" -tags=netgo -ldflags="$LINK_FLAGS -extld=ldd" pkg/game/server.go	
}
listTargets() {
	echo "Available targets (os/arch):"
	echo `go tool dist list`
	exit
}

if [[ $1 == 'all' ]]; then
	for os in `go tool dist list |sed -e 's/\// /g' |cut -d' ' -f1 |sort |uniq`; do
		for arch in `go tool dist list |grep $os |sed 's/\// /g' |cut -d' ' -f2`; do
			compile $os $arch
			
		done
		OS=`echo $tuple |cut -f1 -d' '`
		ARCH=`echo $tuple |cut -f2 -d' '`
		echo $OS $ARCH
#		for ARCH in `go tool dist list |sed -e 's/\// /g' |cut -f2 -d' ' |sort |uniq`; do
#		done
	done
	echo "Install a target build profile to default game-server binary path (./bin/game)? "
	
	read -p 'Specify platform tuple to install to default game binary (blank for linux-amd64): ' -i 'linux-amd64' ANSWER
	#mv "bin/$ANSWER" "bin/game"
	cp bin/game-$ANSWER bin/game
	exit
fi

unset OS
for arch in `go tool dist list|sed -e 's/\// /g' |cut -f1 -d' ' |sort |uniq`; do
	if [[ $1 == $arch ]]; then
		OS=$arch
		break
	fi
done
for os in `go tool dist list|sed -e 's/\// /g' |cut -f2 -d' ' |sort |uniq`; do
	if [[ $2 == $os ]]; then
		ARCH=$os
		break
	fi
done
if [[ -z $OS || -z $ARCH ]]; then
	echo "Could not find target tuple.  Must provide target OS and architecture as arguments, e.g: ./build.sh linux amd64"
	listTargets
	exit
fi
#if [[ -z $ARCH ]]; then
#	echo "Could not find matching target architecture.  Must provide target OS and architecture as arguments, e.g: ./build.sh linux amd64"
#	listTargets
#	exit
#fi
compile $OS $ARCH
if [[ -z $3 ]]; then
	exit
fi
cp bin/game-$3 bin/game