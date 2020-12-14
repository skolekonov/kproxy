VERSION ?= $(shell git describe)
IMG ?= kproxy:${VERSION}

all: kproxy

# Build the controller binary
kproxy:
	go build -o kproxy

kproxy-linux:
	GOOS=linux GOARCH=amd64 go build -o kproxy

# Build the docker image
docker-build: kproxy-linux
	docker build . -f Dockerfile -t ${IMG}

clean:
	rm -f kproxy
	docker rmi -f ${IMG}
