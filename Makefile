build:
	rm -rf ./bin/*
	go build -o bin ./...
	mv ./bin/cmd ./bin/vaulty

run:
	go run ./cmd

