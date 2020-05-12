FROM golang:1.14.1
RUN go get -u github.com/mattn/go-sqlite3
RUN go get -u github.com/satori/go.uuid
RUN go get -u golang.org/x/crypto/bcrypt
WORKDIR /go/src/forum
COPY . .
RUN go build -o main .
CMD ["/go/src/forum/main"]