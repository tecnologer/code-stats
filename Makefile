.SILENT:
.PHONY: *

OUTPUT_DIR=./dist
RELEASE_BRANCH=main
CURRENT_BRANCH=$(shell git rev-parse --abbrev-ref HEAD)

# Get the latest tag matching the pattern vX.Y.Z
LATEST_VERSION=$(shell git tag --list 'v*' --sort=-version:refname | grep -E '^v[0-9]+\.[0-9]+\.[0-9]+$$' | head -n 1)

# Check if RELEASE_VERSION is provided, else increment the patch version of the latest tag
#
# Example: If the RELEASE_VERSION is not provided and the latest tag is v1.2.3, the next version will be v1.2.4
NEXT_VERSION=$(if $(RELEASE_VERSION),$(RELEASE_VERSION),$(shell echo $(LATEST_VERSION) | sed 's/^v//' | awk -F. '{printf "v%d.%d.%d", $$1, $$2, $$3+1}'))

install: build
	go install
	@printf '\033[32m\033[1m%s\033[0m\n' ":: Install complete"

build: prepare
	go build -ldflags "-X 'main.version=v$(NEXT_VERSION)'" -o "code-stats" main.go
	@printf '\033[32m\033[1m%s\033[0m\n' ":: Build complete"

# prepare the code, running the necessary tools before building it
prepare: install-tools
	go generate ./...

# Install required tools, if not already installed
install-tools: check-gobin
	@$(SHELL) -c 'binpath=$${GOPATH:-$$HOME/go}/bin; \
		if ! [ -x "$$binpath/enumer" ]; then \
				echo -e "\033[1;33m!! enumer not found in $$binpath. installing it...\033[0m"; \
				go install github.com/dmarkham/enumer@latest; \
				echo -e "\033[0;32m:: enumer installed\033[0m"; \
			else \
				echo -e "\033[0;32m:: enumer is installed\033[0m"; \
		fi'

# Check if go/bin is in PATH, if not add it
check-gobin:
	@$(SHELL) -c 'binpath=$${GOPATH:-$$HOME/go}/bin; \
		if ! echo $$PATH | grep -q $$binpath; then \
				echo -e "\033[1;33m!! $$binpath is NOT in PATH. Adding it to your PATH.\033[0m"; \
				export PATH="$$GOPATH/bin:$$PATH"; \
			else \
				echo -e "\033[0;32m:: go/bin is in PATH\033[0m"; \
		fi'


lint:
	golangci-lint run --go 1.23.0 --config .golangci.yml ./...

release:
	echo "Current branch: $(CURRENT_BRANCH)"

	if [ "$(CURRENT_BRANCH)" != "$(RELEASE_BRANCH)" ]; then \
		echo "Switching to $(RELEASE_BRANCH) branch"; \
		git checkout $(RELEASE_BRANCH); \
	fi

	echo "Pulling latest changes from $(RELEASE_BRANCH)"
	git pull origin $(RELEASE_BRANCH)

	echo "Latest version: $(LATEST_VERSION)"

	printf "You are about to create a new tag: $(NEXT_VERSION). Press Enter to continue or type a new version: "; \
	read input_version; \
	if [ -z "$$input_version" ]; then \
		final_version="$(NEXT_VERSION)"; \
	else \
		final_version=$$input_version; \
	fi; \
	echo "Next release version $$final_version"; \
	git tag -a $$final_version -m "Release of $$final_version"; \

	printf "Do you want to push the tag to the remote repository? [Y/n] "; \
	read ans; [ "$$ans" = "y" ] || [ "$$ans" = "Y" ] || [ "$$ans" = "" ]

	echo "New version pushed to the remote repository";
	git push --tags
