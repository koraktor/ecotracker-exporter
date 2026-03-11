ARG ARCH="amd64"
ARG OS="linux"
FROM golang:1.26-alpine AS builder

WORKDIR /src
COPY ./. ./
RUN go build -o /build/${OS}-${ARCH}/ecotracker_exporter .

ARG ARCH="amd64"
ARG OS="linux"
FROM golang:1.26-alpine

RUN apk upgrade --no-cache

LABEL authors="koraktor"

COPY --from=builder /build/${OS}-${ARCH}/ecotracker_exporter /bin/ecotracker_exporter

EXPOSE 9776
ENTRYPOINT [ "/bin/sh", "-c", "/bin/ecotracker_exporter --host=${HOST} --port=${PORT} --listen-address=${LISTEN_ADDRESS:-:9776} --log-level=${LOG_LEVEL:-warn} ${*}", "entrypoint" ]
