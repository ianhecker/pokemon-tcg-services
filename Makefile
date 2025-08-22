
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

run: build
	@./$(BIN)/$(APP)

tidy:
	@go mod tidy

mocks:
	mockery

test: tidy
	@$(GOTEST) -v -count=1 ./...

docker-build:
	docker build -t pokemon-tcg-services .

docker-run-cardByID:
	docker run --rm -p 8080:8080 pokemon-tcg-services:latest cardByID --port 8080
