
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
