# Build stage
ARG GOOS=
ARG GOARCH=

FROM golang:1.23 AS builder
WORKDIR /go/src/app

COPY go.sum go.mod config.example.yaml /go/src/app/
RUN go mod download && \
    go mod tidy
COPY . .

RUN go build ./cmd/blacklog-exporter

# Final stage
FROM gcr.io/distroless/static
LABEL maintainer github.com/debman
WORKDIR /app

COPY --from=builder /go/src/app/blacklog-exporter /app/
COPY --from=builder /go/src/app/config.example.yaml /config.yml

ENTRYPOINT ["/app/blacklog-exporter"]
CMD [ "-c","./config.yaml" ]
