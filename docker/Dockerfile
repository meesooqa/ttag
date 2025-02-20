# builder
FROM golang:1.23-alpine AS builder

ENV CGO_ENABLED=0

WORKDIR /build

COPY go.mod go.sum ./
COPY app/ ./app/

WORKDIR /build/app

RUN go mod download
RUN go build -o /build/bin/ttag -ldflags "-s -w"

# run
FROM alpine:latest

COPY --from=builder /build/bin/ttag /srv/ttag
COPY ../templates /srv/templates
COPY ../static /srv/static

WORKDIR /srv
EXPOSE 8080
CMD ["./ttag"]
