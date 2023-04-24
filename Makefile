APP_VERSION=$(shell hack/version.sh)
GO_BUILD_CMD= CGO_ENABLED=0 go build -ldflags="-X main.appVersion=$(APP_VERSION)"

BINARY_NAME=openshift-network-validator
BUILD_DIR=build

.PHONY: all
all: clean lint test build-all package-all

.PHONY: lint
lint:
	@echo "Linting code..."
	@go vet ./...

.PHONY: test
test:
	@echo "Running tests..."
	@go test ./...

.PHONY: pre-build
pre-build:
	@mkdir -p $(BUILD_DIR)

.PHONY: build-linux
build-linux: pre-build
	@echo "Building Linux binary..."
	GOOS=linux GOARCH=amd64 $(GO_BUILD_CMD) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64

.PHONY: build-osx
build-osx: pre-build
	@echo "Building OSX binary..."
	GOOS=darwin GOARCH=amd64 $(GO_BUILD_CMD) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64


.PHONY: build-windows
build-windows: pre-build
	@echo "Building Windows binary..."
	GOOS=windows GOARCH=amd64 $(GO_BUILD_CMD) -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe

.PHONY: build build-all
build-all: build-linux build-osx

.PHONY: package-linux
package-linux:
	@echo "Packaging Linux binary..."
	tar -C $(BUILD_DIR) -zcf $(BUILD_DIR)/$(BINARY_NAME)-$(APP_VERSION)-linux-amd64.tar.gz $(BINARY_NAME)-linux-amd64

.PHONY: package-osx
package-osx:
	@echo "Packaging OSX binary..."
	tar -C $(BUILD_DIR) -zcf $(BUILD_DIR)/$(BINARY_NAME)-$(APP_VERSION)-darwin-amd64.tar.gz $(BINARY_NAME)-darwin-amd64

.PHONY: package-all
package-all: package-linux package-osx

.PHONY: docker
docker:
	docker build --force-rm -t $(BINARY_NAME) .

.PHONY: build-in-docker
build-in-docker: docker
	docker rm -f $(BINARY_NAME) || true
	docker create --name $(BINARY_NAME) $(BINARY_NAME)
	docker cp '$(BINARY_NAME):/opt/' $(BUILD_DIR)
	docker rm -f $(BINARY_NAME)

.PHONY: clean
clean:
	@echo "Cleaning..."
	@rm -Rf $(BUILD_DIR)