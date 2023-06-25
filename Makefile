VERSION=1.0
BUILD=`date +%FT%T%z`
LDFLAGS=-ldflags "-w -s -X main.Version=${VERSION} -X main.Build=${BUILD} -X main.Entry=f1"
 
NAME   := ciceros/vaulty
IMG    := ${NAME}:${VERSION}
LATEST := ${NAME}:latest

build:
	rm -rf ./bin/*
	go build -o bin ${LDFLAGS} ./...
	mv ./bin/cmd ./bin/vaulty

run:
	go run ./cmd

image:
	docker build -t ${IMG} .
	docker tag ${IMG} ${LATEST}

push:
	docker push ${NAME}

login:
	docker log -u ${DOCKER_USER} -p ${DOCKER_PASS}

