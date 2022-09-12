FROM docker.io/golang:1.18.6-alpine AS builder

WORKDIR /build

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
      -ldflags='-w -s -extldflags "-static"' -a \
      -o /build/x-delay ./cmd/x-delay

FROM gcr.io/distroless/static AS app

COPY --from=builder /build/x-delay /go/bin/x-delay

EXPOSE 10025

ENTRYPOINT ["/go/bin/x-delay"]



