FROM golang:1.20 as builder

ENV GOOS linux
ENV CGO_ENABLED 0

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o app

FROM scratch
COPY --from=builder app .
ENTRYPOINT ["./app"]