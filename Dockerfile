# Build stage

FROM golang:1.23 AS builder
WORKDIR /go/src/app

COPY go.sum go.mod config.example.yaml /go/src/app/
RUN go mod download && \
    go mod tidy
COPY . .

RUN go build ./cmd/blacklog-exporter

# Final stage

FROM gcr.io/distroless/base
LABEL maintainer=github.com/debman
WORKDIR /app

COPY --from=builder /go/src/app/blacklog-exporter .
COPY --from=builder /go/src/app/config.example.yaml .

ENTRYPOINT ["/app/blacklog-exporter"]
CMD [ "-c","./config.yaml" ]
