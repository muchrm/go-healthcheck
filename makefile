repo = muchrm/go-healthcheck
commit = latest
name = krud

check: test vet ## Runs all tests
install:
	go install -v .
build:
	go build main.go
docker:
	docker build -f Dockerfile -t $(repo):$(commit) .
vet: ## Run the vet tool
	go vet $(shell go list ./... | grep -v /vendor/)

clean: ## Clean up build artifacts
	go clean

test: ## Run the  tests
	echo "" > coverage.txt
	for d in $(shell go list ./... | grep -v vendor); do \
		go test -race -coverprofile=profile.out -covermode=atomic $$d || exit 1; \
		[ -f profile.out ] && cat profile.out >> coverage.txt && rm profile.out; \
	done
