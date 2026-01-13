GOLANGCI_LINT_VERSION := 2.6.0

./bin/golangci-lint:
	mkdir -p ./bin
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b "./bin" "v${GOLANGCI_LINT_VERSION}";

lint:
	$(MAKE) ./bin/golangci-lint
	./bin/golangci-lint run ./...