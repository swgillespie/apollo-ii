PROJECT := apollo-ii
GIT_DESCRIBE := $(shell git describe --tags --dirty 2>/dev/null)
GIT_COMMIT := $(shell git rev-parse HEAD)

build:
	go install -ldflags "-X github.com/swgillespie/apollo-ii/pkg/engine.Version=${GIT_DESCRIBE}"

test:
	go test github.com/swgillespie/apollo-ii/pkg/engine

test-bench:
	go test github.com/swgillespie/apollo-ii/pkg/engine -bench .

test-cover:
	go test github.com/swgillespie/apollo-ii/pkg/engine -coverprofile=c.out
	go tool cover -html=c.out
	rm c.out