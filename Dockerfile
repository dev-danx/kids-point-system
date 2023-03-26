FROM golang:1.20-bullseye AS builder

RUN go install golang.org/dl/go1.20@latest \
    && go1.20 download

COPY . / /src/

WORKDIR /src/

RUN go mod download

RUN GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /go/bin/server /src/main.go
FROM alpine:latest

COPY --from=builder /go/bin/* /go/bin/

CMD /go/bin/server
