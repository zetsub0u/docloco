FROM golang:latest
WORKDIR /go/src/github.com/zetsub0u/docloco
COPY . .
EXPOSE 8080
CMD go run main.go