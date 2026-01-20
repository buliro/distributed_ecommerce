FROM golang:1.23-alpine AS base

WORKDIR /src

ENV GOPROXY=https://proxy.golang.org,direct \
    GOSUMDB=sum.golang.org

RUN apk add --no-cache ca-certificates tzdata git

COPY go.mod go.sum ./
RUN /bin/sh -c 'set -e; for attempt in 1 2 3; do \
        go mod download && exit 0; \
        echo "go mod download failed (attempt ${attempt}), retrying..." >&2; \
        sleep $((attempt * 5)); \
    done; exit 1'

COPY . .

FROM base AS test
RUN go test ./...

FROM base AS builder
ARG TARGETOS=linux
ARG TARGETARCH=amd64
RUN mkdir -p /out && \
    CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} \
    go build -ldflags="-s -w" -o /out/api ./cmd/main.go

FROM alpine:3.20 AS runtime

RUN apk add --no-cache ca-certificates

WORKDIR /app

COPY --from=builder /out/api ./api

EXPOSE 8080

ENTRYPOINT ["./api"]
