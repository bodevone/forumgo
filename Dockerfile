FROM golang:1.14.1
RUN go get -u github.com/satori/go.uuid
RUN go get -u github.com/mattn/go-sqlite3
WORKDIR /go/src/forum
COPY . .
RUN go build -o main .