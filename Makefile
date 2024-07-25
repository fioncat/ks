COMMIT_ID=$(shell git rev-parse HEAD)
COMMIT_ID_SHORT=$(shell git rev-parse --short HEAD)

TAG=$(shell git describe --tags --abbrev=0 2>/dev/null)

DATE=$(shell date '+%FT%TZ')

# If current commit is tagged, use tag as version, else, use dev-${COMMIT_ID} as version
VERSION=$(shell git tag --points-at ${COMMIT_ID})
VERSION:=$(if $(VERSION),$(VERSION),dev-${COMMIT_ID_SHORT})

.PHONY: version
version:
	@echo ${VERSION}

.PHONY: fmt
fmt:
	@find . -name \*.go -exec goimports -w {} \;

.PHONY: build
build:
	@CGO_ENABLED=0 go build -ldflags="-X 'main.Version=$(VERSION)' -X 'main.Commit=$(COMMIT_ID)' -X 'main.BuildDate=$(DATE)'" -o bin/kser ./kser

.PHONY: install
install:
	@CGO_ENABLED=0 go install -ldflags="-X 'main.Version=$(VERSION)' -X 'main.Commit=$(COMMIT_ID)' -X 'main.BuildDate=$(DATE)'" ./kser

.PHONY: install-check
install-check:
	@go install github.com/fzipp/gocyclo/cmd/gocyclo@latest
	@go install github.com/client9/misspell/cmd/misspell@latest
	@go install github.com/gordonklaus/ineffassign@latest

.PHONY: check
check:
	@echo "==> check ineffassign"
	@ineffassign ./...
	@echo "==> check spell"
	@find . -type f -name '*.go' | xargs misspell -error
	@echo "==> check gocyclo"
	@gocyclo -over 15 .
	@echo "==> go vet"
	@go vet ./...

.PHONY: clean
clean:
	@rm -rf bin
