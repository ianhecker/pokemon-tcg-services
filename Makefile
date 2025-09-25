
APP=pokemon-tcg-services
BIN="bin"
GOTEST=gotest

.PHONY: clean bin build run tidy mocks test

clean:
	@rm -rf $(BIN)/$(APP)
	@rm -rf $(BIN)

bin:
	@mkdir -p $(BIN)

build: bin
	@go build -o $(BIN)/$(APP)

tidy:
	@go mod tidy

mocks:
	mockery

test: tidy
	@$(GOTEST) -v -count=1 ./...

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