# syntax=docker/dockerfile:1.7

FROM golang:1.23-alpine AS base

WORKDIR /src

ENV GOPROXY=https://proxy.golang.org,direct \
    GOSUMDB=sum.golang.org

RUN --mount=type=cache,target=/var/cache/apk \
    --mount=type=cache,target=/var/lib/apk \
    apk add --no-cache ca-certificates tzdata git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

FROM base AS test
RUN go test ./...

FROM base AS builder
ARG TARGETOS=linux
ARG TARGETARCH=amd64
RUN mkdir -p /out && \
    CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} \
    go build -ldflags="-s -w" -o /out/api ./cmd/main.go

FROM gcr.io/distroless/static-debian12:nonroot

WORKDIR /app

COPY --from=builder /out/api ./api

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["./api"]
