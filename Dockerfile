# Build stage
ARG GOOS=
ARG GOARCH=

FROM golang:1.23 AS builder
WORKDIR /go/src/app

COPY go.sum go.mod config.example.yaml /go/src/app/
RUN go mod download && \
    go mod tidy && \
    go vet -v
COPY . .

WORKDIR /go/src/app/cmd/blacklog-exporter
RUN go build -o /blacklog-exporter

# Final stage
LABEL maintainer github.com/debman
FROM gcr.io/distroless/static
WORKDIR /app

COPY --from=builder /go/src/app/cmd/blacklog-exporter /app/
COPY --from=builder /go/src/app/config.example.yaml /config.yml

ENTRYPOINT ["/app/blacklog-exporter"]
CMD [ "-c","./config.yaml" ]
