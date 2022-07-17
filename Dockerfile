FROM golang:1.17-buster AS build

ENV GOPATH=/

WORKDIR /src/
COPY ./ /src/

RUN go mod download; go build -o /scanner_backend_api ./cmd/main.go

FROM alpine:latest

COPY --from=build /scanner_backend_api /scanner_backend_api
COPY ./configs/ /configs/
COPY ./wait-for-postgres.sh ./

RUN apk --no-cache add postgresql-client && chmod +x wait-for-postgres.sh

CMD ["/scanner_backend_api"]