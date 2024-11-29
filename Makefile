fmt:
	go fmt ./...

vet:
	go vet ./...

build: fmt vet
	go build -o bin/spinit

unit-tests:
	ginkgo -v --skip-package=tests ./...

integration-tests:
	ginkgo -v --skip-package=cmd ./...

lint:
	golangci-lint run ./...