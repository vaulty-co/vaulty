VERSION=`git describe --tags`
BUILD=`date +%FT%T%z`
LDFLAGS=-ldflags "-w -s -X main.Version=${VERSION} -X main.Build=${BUILD} -X main.Entry=f1"

build:
	rm -rf ./bin/*
	go build ${LDFLAGS} -o bin ./...
	mv ./bin/cmd ./bin/vaulty

run:
	go run ./cmd

image:
	docker build -t vaulty:${VERSION} .

push:
	docker push vaulty/vaulty:${VERSION}
