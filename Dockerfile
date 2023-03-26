FROM golang:1.20-bullseye AS builder

RUN go install golang.org/dl/go1.20@latest \
    && go1.20 download

COPY . / /src/

RUN go mod download

RUN GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /app/bin/server /src/main.go
FROM alpine:latest

COPY --from=builder /app/bin/* /app/bin/

CMD /app/bin/server
