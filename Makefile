VERSION := $(shell cat ./VERSION)
COMMIT_HASH := $(shell git rev-parse HEAD 2>/dev/null || true)
BUILD_TIME := $(shell date +%s)
BINNAME := rhcos-shell-greeter

# Used during all builds
LDFLAGS := -X main.version=${VERSION} -X main.commitHash=${COMMIT_HASH} -X main.buildTime=${BUILD_TIME}

BIN_DIR ?= /usr/bin

.PHONY: help build clean deps install lint static

help:
	@echo "Targets:"
	@echo " - build: Build the target binary"
	@echo " - static: Build a static binary"
	@echo " - clean: Clean up after build"
	@echo " - deps: Install required tool and dependencies for building"
	@echo " - install: Install build results to the system"
	@echo " - lint: Run golint"
	@echo ""
	@echo "Variables:"
	@echo " - PREFIX: The root location to install. This prepends to all *_DIR variables. Set to: ${PREFIX}"
	@echo " - BIN_DIR: The directory that houses binaries. Set to: ${BIN_DIR}"
	@echo " - VERSION: Generally not overridden. The output of the VERSION file. Set to: ${VERSION}"
	@echo " - COMMIT_HASH: Generally not overridden. The git hash the code was built from. Set to: ${COMMIT_HASH}"
	@echo " - BUILD_TIME: Generally not overridden. The unix time of the build. Set to: ${BUILD_TIME}"

build: clean
	go build -ldflags '${LDFLAGS}' -o ${BINNAME} main.go
	strip ${BINNAME}

static: clean
	CGO_ENABLED=0 go build -ldflags '${LDFLAGS} -w -extldflags "-static"' -a -o ${BINNAME} main.go

clean:
	rm -f ${BINNAME}

deps:
	go get -u github.com/golang/dep/cmd/dep
	dep ensure -v

install: clean build
	install -d ${PREFIX}${BIN_DIR}
	install --mode 755 ${BINNAME} ${PREFIX}${BIN_DIR}/${BINNAME}

lint:
	go get -u github.com/golang/lint/golint
	golint .
