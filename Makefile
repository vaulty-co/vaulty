build:
	rm -rf ./bin/*
	go build -o bin ./...
	mv ./bin/cmd ./bin/vaulty

run:
	go run ./cmd

image:
	docker build -t vaulty .
push:
	docker push vaulty/vaulty:latest
