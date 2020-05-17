FROM golang:1.14-alpine as builder
RUN mkdir /build 
ADD . /build/
WORKDIR /build 
RUN go get -d -v ./...
RUN go build -o bin ./... && mv ./bin/cmd ./bin/vaulty
FROM alpine
RUN adduser -S -D -H -h /vaulty appuser
USER appuser
COPY --from=builder /build/bin/vaulty /vaulty/
WORKDIR /vaulty
CMD ["./vaulty", "proxy"]
EXPOSE 8080
