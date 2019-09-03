PWD := $(shell pwd)
GOPATH := $(shell go env GOPATH)
LDFLAGS := $(shell echo "")

GOOS := $(shell go env GOOS)
GOOSALT ?= 'linux'
ifeq ($(GOOS),'darwin')
  GOOSALT = 'mac'
endif

BUILD_LDFLAGS := '$(LDFLAGS)'

all: build

build:
	@echo "building tablef binary to ./tablef"
	@GOPROXY=https://proxy.golang.org GO111MODULE=on GO_FLAGS="" CGO_ENABLED=0 go build -tags kqueue --ldflags $(BUILD_LDFLAGS)

build-linux:
	@echo "building tablef-linux-amd64 to ./tablef-linux-amd64"
	@GOOS=linux GOARCH=amd64 GOPROXY=https://proxy.golang.org GO111MODULE=on GO_FLAGS="" CGO_ENABLED=0 go build -o tablef-linux-amd64 -tags kqueue --ldflags $(BUILD_LDFLAGS)

build-darwin:
	@echo "building tablef-linux-amd64 to ./tablef-darwin-amd64"
	@GOOS=darwin GOARCH=amd64 GOPROXY=https://proxy.golang.org GO111MODULE=on GO_FLAGS="" CGO_ENABLED=0 go build -o tablef-darwin-amd64 -tags kqueue --ldflags $(BUILD_LDFLAGS)

build-windows:
	@echo "building tablef-linux-amd64 to ./tablef-windows-amd64"
	@GOOS=windows GOARCH=amd64 GOPROXY=https://proxy.golang.org GO111MODULE=on GO_FLAGS="" CGO_ENABLED=0 go build -o tablef-windows-amd64 -tags kqueue --ldflags $(BUILD_LDFLAGS)

release: build-linux build-darwin build-windows
	@echo "releasing multi-platform tablef binaries to releases/"
	@mkdir -p releases/
	@mv tablef-linux-amd64	releases/tablef-linux-amd64
	@mv tablef-windows-amd64	releases/tablef-windows-amd64
	@mv tablef-darwin-amd64	releases/tablef-darwin-amd64
	@tar cvzf tablef-multiplatform.tar.gz releases/
