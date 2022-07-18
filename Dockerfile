FROM golang:1.17-buster AS build

ENV GOPATH=/

WORKDIR /src/
COPY ./ /src/

RUN go mod download 
RUN go build -o server ./cmd/main.go

CMD ["./sever"]
