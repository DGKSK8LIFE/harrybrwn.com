GO_ENV=GOOS=linux GOARCH=amd64 CGO_ENABLED=0
DOCKER_BUILD=$(shell pwd)/.docker_build

BIN=$(DOCKER_BUILD)/go-getting-started

$(BIN): clean
	mkdir -p $(DOCKER_BUILD)
	$(GO_ENV) go build -o $(BIN) .

clean:
	rm -rf $(DOCKER_BUILD)