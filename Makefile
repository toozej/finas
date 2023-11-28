# Set sane defaults for Make
SHELL = bash
.DELETE_ON_ERROR:
MAKEFLAGS += --warn-undefined-variables
MAKEFLAGS += --no-builtin-rules

# Set default goal such that `make` runs `make help`
.DEFAULT_GOAL := help

# Build info
BUILDER = $(shell whoami)@$(shell hostname)
NOW = $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")

# Version control
VERSION = $(shell git describe --tags --dirty --always)
COMMIT = $(shell git rev-parse --short HEAD)
BRANCH = $(shell git rev-parse --abbrev-ref HEAD)

# Linker flags
PKG = $(shell head -n 1 go.mod | cut -c 8-)
VER = $(PKG)/version
LDFLAGS = -s -w \
	-X $(VER).Version=$(or $(VERSION),unknown) \
	-X $(VER).Commit=$(or $(COMMIT),unknown) \
	-X $(VER).Branch=$(or $(BRANCH),unknown) \
	-X $(VER).BuiltAt=$(NOW) \
	-X $(VER).Builder=$(BUILDER)
	
OS = $(shell uname -s)
ifeq ($(OS), Linux)
	OPENER=xdg-open
else
	OPENER=open
endif

.PHONY: all vet test cover run release-test release sign verify release-verify install get-cosign-pub-key docker-login pre-commit-install pre-commit-run pre-commit pre-reqs update-golang-version docs docs-generate docs-serve clean help wiptest

all: vendor update-deps vet pre-commit clean test cover build sign verify run ## Run default workflow using locally installed Golang toolchain
release-verify: release sign verify ## Release and verify using locally installed Golang toolchain
pre-reqs: pre-commit-install ## Install pre-commit hooks and necessary binaries

wiptest: build ## TODO remove wiptest
	$(CURDIR)/finas --debug helloworld
	rm -f htpasswdFileName && $(CURDIR)/finas --debug htpasswd username password > htpasswdFileName
	$(CURDIR)/finas --debug helloworld "some shit goes here"

get-cosign-pub-key: ## Get finas Cosign public key from GitHub
	test -f $(CURDIR)/finas.pub || curl --silent https://raw.githubusercontent.com/toozej/finas/main/finas.pub -O

update-deps: ## Run `go get -t -u ./...` to update Go module dependencies
	go get -t -u ./...

vet: ## Run `go vet` using locally installed golang toolchain
	go vet $(CURDIR)/...

vendor: ## Run `go mod vendor` using locally installed golang toolchain
	go mod vendor

test: ## Run `go test` using locally installed golang toolchain
	go test -coverprofile c.out -v $(CURDIR)/...
	@echo -e "\nStatements missing coverage"
	@grep -v -e " 1$$" c.out

cover: ## View coverage report in web browser
	go tool cover -html=c.out

build: ## Run `go build` using locally installed golang toolchain
	CGO_ENABLED=0 go build -ldflags="$(LDFLAGS)" $(CURDIR)/cmd/finas/

run: ## Run locally built binary
	$(CURDIR)/finas help
	$(CURDIR)/finas --debug helloworld

release-test: ## Build assets and test goreleaser config using locally installed golang toolchain and goreleaser
	goreleaser check
	goreleaser build --rm-dist --snapshot

release: test docker-login ## Release assets using locally installed golang toolchain and goreleaser
	if test -e $(CURDIR)/finas.key && test -e $(CURDIR)/.env; then \
		export `cat $(CURDIR)/.env | xargs` && goreleaser release --rm-dist; \
	else \
		echo "no cosign private key found at $(CURDIR)/finas.key. Cannot release."; \
	fi

sign: test ## Sign locally installed golang toolchain and cosign
	if test -e $(CURDIR)/finas.key && test -e $(CURDIR)/.env; then \
		export `cat $(CURDIR)/.env | xargs` && cosign sign-blob --key=$(CURDIR)/finas.key --output-signature=$(CURDIR)/finas.sig $(CURDIR)/finas; \
	else \
		echo "no cosign private key found at $(CURDIR)/finas.key. Cannot release."; \
	fi

verify: get-cosign-pub-key ## Verify locally compiled binary
	# cosign here assumes you're using Linux AMD64 binary
	cosign verify-blob --key $(CURDIR)/finas.pub --signature $(CURDIR)/finas.sig $(CURDIR)/finas

install: build verify ## Install compiled binary to local machine
	sudo cp $(CURDIR)/finas /usr/local/bin/finas
	sudo chmod 0755 /usr/local/bin/finas
	mkdir ~/.config/finas/ && cp -r $(CURDIR)/config/*.json ~/.config/finas/

pre-commit: pre-commit-install pre-commit-run ## Install and run pre-commit hooks

pre-commit-install: ## Install pre-commit hooks and necessary binaries
	# golangci-lint
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	# goimports
	go install golang.org/x/tools/cmd/goimports@latest
	# gosec
	go install github.com/securego/gosec/v2/cmd/gosec@latest
	# staticcheck
	go install honnef.co/go/tools/cmd/staticcheck@latest
	# go-critic
	go install github.com/go-critic/go-critic/cmd/gocritic@latest
	# structslop
	go install github.com/orijtech/structslop/cmd/structslop@latest
	# shellcheck
	command -v shellcheck || sudo dnf install -y ShellCheck || sudo apt install -y shellcheck || brew install shellcheck
	# checkmake
	go install github.com/mrtazz/checkmake/cmd/checkmake@latest
	# goreleaser
	go install github.com/goreleaser/goreleaser@latest
	# syft
	go install github.com/anchore/syft/cmd/syft@latest
	# cosign
	go install github.com/sigstore/cosign/cmd/cosign@latest
	# go-licenses
	go install github.com/google/go-licenses@latest
	# go vuln check
	go install golang.org/x/vuln/cmd/govulncheck@latest
	# install and update pre-commits
	pre-commit install
	pre-commit autoupdate

pre-commit-run: ## Run pre-commit hooks against all files
	pre-commit run --all-files
	# manually run the following checks since their pre-commits aren't working or don't exist
	go-licenses report github.com/toozej/finas/cmd/finas
	govulncheck ./...

update-golang-version: ## Update to latest Golang version across the repo
	@VERSION=`curl -s "https://go.dev/dl/?mode=json" | jq -r '.[0].version' | sed 's/go//' | cut -d '.' -f 1,2`; \
	echo "Updating Golang to $$VERSION"; \
	./scripts/update_golang_version.sh $$VERSION

docs: docs-generate docs-serve ## Generate and serve documentation

docs-generate:
	docker build -f $(CURDIR)/Dockerfile.docs -t toozej/finas:docs . 
	docker run --rm --name finas-docs -v $(CURDIR):/package -v $(CURDIR)/docs:/docs toozej/finas:docs

docs-serve: ## Serve documentation on http://localhost:9000
	docker run -d --rm --name finas-docs-serve -p 9000:3080 -v $(CURDIR)/docs:/data thomsch98/markserv
	$(OPENER) http://localhost:9000/docs.md
	@echo -e "to stop docs container, run:\n"
	@echo "docker kill finas-docs-serve"

clean: ## Remove any locally compiled binaries
	rm -f $(CURDIR)/finas

help: ## Display help text
	@grep -E '^[a-zA-Z_-]+ ?:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
