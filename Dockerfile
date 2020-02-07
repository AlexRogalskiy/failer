FROM golang:1.13.7-buster as golang

WORKDIR /failer

COPY go.mod go.mod
RUN go mod download
COPY main.go main.go

RUN CGO_ENABLED=0 GOOS=linux go build -o failer -mod=readonly ./main.go

FROM debian:buster-20191224-slim

RUN apt-get update \
    && apt-get install -y --no-install-recommends curl \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /failer

COPY --from=golang /failer/failer .

ENTRYPOINT ["/failer/failer"]
