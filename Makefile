.PHONY: build

VERSION=v0.2.0

PLATFORMS := linux/amd64 darwin/amd64 darwin/arm64 windows/amd64

temp = $(subst /, ,$@)
os = $(word 1, $(temp))
arch = $(word 2, $(temp))

BINARY_NAME=fastxviz

build:
	go build

build-all-platforms: $(PLATFORMS)
	mv bin/fastxviz-$(VERSION)-windows-amd64 bin/fastxviz-$(VERSION)-windows-amd64.exe

gzip:
	gzip bin/*

$(PLATFORMS):
	GOOS=$(os) GOARCH=$(arch) go build -o bin/$(BINARY_NAME)-$(VERSION)-$(os)-$(arch) .

clean:
	rm -rf bin
