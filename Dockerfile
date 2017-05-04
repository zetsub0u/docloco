FROM golang
WORKDIR /go/src/github.com/zetsub0u/docloco
ADD . .
EXPOSE 8080
CMD go run main.go