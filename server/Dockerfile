# syntax=docker/dockerfile:experimental
FROM golang:bookworm as builder_base

ENV GOPRIVATE="github.com/MrColorado"
ARG DIRECTORY

WORKDIR /go/src/github.com/MrColorado/backend/server

COPY ${DIRECTORY}/go.mod .
COPY ${DIRECTORY}/go.sum .
COPY internal/go.mod /go/src/github.com/MrColorado/backend/internal/go.mod
COPY logger/go.mod /go/src/github.com/MrColorado/backend/logger/go.mod
RUN go mod download

# Build image
FROM builder_base as builder

COPY ${DIRECTORY}/. .
COPY internal/ /go/src/github.com/MrColorado/backend/internal/
COPY logger/ /go/src/github.com/MrColorado/backend/logger/
RUN make install

# Deployment image
FROM golang:bookworm

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /go/bin/server /service

ENTRYPOINT ["/service"]