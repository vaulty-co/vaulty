FROM golang:1.14-alpine as builder
RUN mkdir /build 
ADD go.sum /build
ADD go.mod /build
WORKDIR /build 
RUN go mod download
ADD . /build/
RUN go build -o bin ./... && mv ./bin/cmd ./bin/vaulty
FROM alpine
RUN adduser -S -D -H -h /vaulty appuser
RUN mkdir /.vaulty 
USER appuser
COPY --from=builder /build/bin/vaulty /vaulty/
WORKDIR /vaulty
CMD ["./vaulty", "proxy", "-r", "/.vaulty/routes.json", "--ca", "/.vaulty"]
EXPOSE 8080
