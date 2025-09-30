APP = pokemon-tcg-services
BIN = "bin"

TEST_CMD ?= $(shell if command -v gotest >/dev/null 2>&1; then echo "gotest"; else echo "go test"; fi)
TEST_PACKAGES ?= $(shell go list ./... | grep -vE 'mocks|cmd')

.PHONY: clean bin build run tidy

clean:
	@rm -rf $(BIN)/$(APP)
	@rm -rf $(BIN)
	@rm coverage.out

bin:
	@mkdir -p $(BIN)

build: bin
	@go build -o $(BIN)/$(APP)

tidy:
	@go mod tidy

.PHONY: mocks test coverage view-coverage

mocks:
	mockery

test: mocks tidy
	@$(TEST_CMD) -v -count=1 ./...

coverage:
	@go test -coverprofile=coverage.out $(TEST_PACKAGES)

view-coverage: coverage
	@go tool cover -html=coverage.out

.PHONY: docker-build run-card-pricer hello-world

docker-build:
	docker build --no-cache -t pokemon-tcg-services .

run-card-pricer: docker-build
	docker run --rm -p 8080:8080 \
	 --env-file .env \
	pokemon-tcg-services:latest \
	card-pricer --port 8080

# Requires run-card-pricer to be running
hello-world:
	./sh/hello-world.sh