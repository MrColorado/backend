FROM golang:1.20 as builder

ENV GOOS linux
ENV CGO_ENABLED 0

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

# COPY configuration ./
# COPY converter ./
# COPY dataWrapper ./
# COPY grpcWrapper ./
# COPY scraper ./
# COPY models ./
# COPY utils ./
# COPY main.go ./

COPY . .

RUN go build -o app

FROM scratch
COPY --from=builder app .
ENTRYPOINT ["./app"]

# FROM golang:1.20-alpine as builder

# EXPOSE 55051

# WORKDIR /app

# COPY go.mod go.sum ./
# RUN go mod dowload

# COPY . .
# RUN CGO_ENABLED=0 GOOS=linux go build -o myapp

# FROM scratch
# COPY --from=builder /app/myapp /

# ENTRYPOINT ["/myapp"]