
BINDIR      := $(CURDIR)/bin
GOPATH = $(shell go env GOPATH)

BINNAME     ?= x-tracer
AGENT_NAME  ?= x-agent

NS ?= sheenam3
IMAGE  ?= x-agent
VERSION ?= latest
# go option
PKG        := ./...
TAGS       :=
TESTS      := .
TESTFLAGS  :=
LDFLAGS    := -w -s
GOFLAGS    :=
SRC        := $(shell find . -type f -name '*.go' -print)

GIT_COMMIT = $(shell git rev-parse HEAD)
GIT_SHA    = $(shell git rev-parse --short HEAD)
GIT_TAG    = $(shell git describe --tags --abbrev=0 --exact-match 2>/dev/null)
GIT_DIRTY  = $(shell test -n "`git status --porcelain`" && echo "dirty" || echo "clean")

.PHONY: all
all: tracer agent build-image push-image

# ------------------------------------------------------------------------------
#  build

.PHONY: tracer
tracer: $(BINDIR)/$(BINNAME)
$(BINDIR)/$(BINNAME): $(SRC)
	@echo "====    Build x-tracer    ===="
	GO111MODULE=on go build $(GOFLAGS) -tags '$(TAGS)' -ldflags '$(LDFLAGS)' -o $(BINDIR)/$(BINNAME) ./cmd/x-tracer


.PHONY: agent
agent: $(BINDIR)/$(AGENT_NAME)
$(BINDIR)/$(AGENT_NAME): $(SRC)
	@echo "====    Build x-agent    ===="
	GO111MODULE=on go build $(GOFLAGS) -tags '$(TAGS)' -ldflags '$(LDFLAGS)' -o $(BINDIR)/$(AGENT_NAME) ./cmd/x-agent

.PHONY: build-image
build-image:
	docker build --pull=false -f build/Dockerfile -t "x-agent" . --no-cache

.PHONY: push-image
push-image:
	docker tag x-agent $(NS)/$(IMAGE):$(VERSION)
	docker push  $(NS)/$(IMAGE):$(VERSION)
#	docker save x-agent | gzip > x-agent.tar.gz
#	scp x-agent.tar.gz root@node2:~/
#	ssh root@node2 'docker load < x-agent.tar.gz


# ------------------------------------------------------------------------------
#  clean
.PHONY: clean
clean:
	@rm -rf $(BINDIR) ./_dist
	@docker rmi x-agent
	@docker rmi sheenam3/x-agent:latest
#	@rm x-agent.tar.gz
#	@ssh root@node2 'docker rmi x-agent'
#	@ssh root@node2 'rm ~/x-agent.tar.gz' 
