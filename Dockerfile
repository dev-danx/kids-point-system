FROM golang:1.20-bullseye AS builder

RUN go install golang.org/dl/go1.20@latest \
    && go1.20 download

COPY . / /src/

WORKDIR /src/

RUN go mod download

RUN GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /go/bin/server /src/main.go
FROM gcr.io/distroless/base-debian11

COPY --from=builder /go/bin/* /go/bin/
COPY --from=builder /src/template/* /go/bin/template/
COPY --from=builder /src/data.json /go/bin/data.json

ENTRYPOINT ["/go/bin/server"]