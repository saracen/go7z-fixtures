PLATFORMS := linux/386 linux/amd64 linux/arm linux/arm64 linux/ppc64 linux/ppc64le linux/mips linux/mipsle linux/mips64 linux/mips64le

temp = $(subst /, ,$@)
os = $(word 1, $(temp))
arch = $(word 2, $(temp))

release: $(PLATFORMS)

$(PLATFORMS):
	mkdir -p binaries
	GOOS=$(os) GOARCH=$(arch) go build -o 'executables/hello-world-$(os)-$(arch)' -ldflags="-s -w" helloworld.go

helloworld: $(PLATFORMS)

random_binary:
	if [ ! -f random_01.bin ]; then dd if=/dev/urandom of=random_01.bin bs=2MiB count=1; fi
	if [ ! -f random_02.bin ]; then dd if=/dev/urandom of=random_02.bin bs=2MiB count=1; fi

testdata: random_binary helloworld
	# empty archive
	7z a archives/empty.7z *.nothing

	# all executables
	7z a archives/executables.7z executables/*

	# all executables, BCJ only, lzma, solid
	7z a archives/executables-bcj2.7z executables/* -m0=BCJ -m1=LZMA:d=21 -ms -mmt

	# all executables, BCJ2 only, lzma, solid
	7z a archives/executables-bcj2.7z executables/* -m0=BCJ2 -m1=LZMA:d=21 -ms -mmt

	# 386, amd64 executable, BCJ2
	7z a archives/executables-bcj2-386-amd64.7z executables/hello-world-linux-386 executables/hello-world-linux-amd64
	
	# random.bin, delta filter
	7z a archives/delta.7z random_01.bin random_02.bin -m0=delta:2 -m1=delta:2

	# random.bin, delta filter, lzma
	7z a archives/delta-lzma.7z random_01.bin random_02.bin -m0=delta:2 -m1=delta:2 -m2=lzma:fb128

	# random.bin, ppmd, bzip2, deflate, copy
	7z a archives/ppmd.7z random_01.bin -m0=ppmd
	7z a archives/bzip2.7z random_01.bin -m0=bzip2
	7z a archives/deflate.7z random_01.bin -m0=deflate
	7z a archives/copy.7z random_01.bin -m0=copy
	7z a archives/ppmd-bzip2-deflate-copy.7z random_01.bin random_02.bin -m0=ppmd -m1=bzip2 -m2=deflate -m3=copy