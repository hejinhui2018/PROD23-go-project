FROM golang:1.23.12

WORKDIR /src
COPY . .
RUN go test ./... && go vet ./... && go build ./...
