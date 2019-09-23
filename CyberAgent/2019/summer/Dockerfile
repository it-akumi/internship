FROM golang:1.12

WORKDIR /go/src
COPY . .

ENV GOOS=linux
ENV GOARCH=amd64
ENV GO111MODULE=on
RUN go build -o server api.go model.go

CMD ["./server"]
