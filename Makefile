GO_IMAGE := golang:1.23
DOCKER_GO := docker run --rm -u $(shell id -u):$(shell id -g) \
	-v $(CURDIR):/src -w /src \
	-e GOCACHE=/tmp/gocache -e GOOS=windows -e GOARCH=amd64 -e CGO_ENABLED=0 \
	$(GO_IMAGE)

.PHONY: build vet clean

build:
	$(DOCKER_GO) go build -trimpath -ldflags "-s -w" -o bin/vdi_poke.exe .

vet:
	$(DOCKER_GO) go vet ./...

clean:
	rm -rf bin
