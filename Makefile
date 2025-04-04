export GO ?= go
export CGO_ENABLED ?= 0
export GOARCH ?= $(shell "$(GO)" env GOARCH)
export GOOS ?= $(shell "$(GO)" env GOOS)
export GOPATH ?= $(shell "$(GO)" env GOPATH)
export GOROOT ?= $(shell "$(GO)" env GOROOT)

BUILD_FLAGS ?= -ldflags "-s -w"
OUTPUT ?= kube-proxy-healthz

.PHONY: build clean

build:
	"$(GO)" build $(BUILD_FLAGS) -o "$(OUTPUT)" main.go

clean:
	rm -f "$(OUTPUT)"
