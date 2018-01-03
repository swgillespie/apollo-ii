PROJECT := apollo-ii
GIT_DESCRIBE := $(shell git describe --tags --dirty 2>/dev/null)
GIT_COMMIT := $(shell git rev-parse HEAD)

build:
	go install -ldflags "-X github.com/swgillespie/apollo-ii/pkg/engine.Version=${GIT_DESCRIBE}"

test:
	go test github.com/swgillespie/apollo-ii/pkg/engine

test-verbose:
	go test github.com/swgillespie/apollo-ii/pkg/engine -v

test-bench:
	go test github.com/swgillespie/apollo-ii/pkg/engine -bench .

test-cover:
	go test github.com/swgillespie/apollo-ii/pkg/engine -coverprofile=coverage-engine.out
	go tool cover -html=coverage-engine.out
	rm coverage-engine.out