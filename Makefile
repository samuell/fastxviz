.PHONY: build

PLATFORMS := linux/amd64 darwin/amd64 darwin/arm64 windows/amd64

temp = $(subst /, ,$@)
os = $(word 1, $(temp))
arch = $(word 2, $(temp))

BINARY_NAME=fastxviz

build:
	go build

build-all-platforms: $(PLATFORMS)
	mv bin/fastxviz-windows-amd64 bin/fastxviz-windows-amd64.exe

gzip:
	gzip bin/*

$(PLATFORMS):
	GOOS=$(os) GOARCH=$(arch) go build -o bin/$(BINARY_NAME)-$(os)-$(arch) .

