FROM golang:1.26-alpine AS builder

ARG ARCH="amd64"
ARG OS="linux"

WORKDIR /src
COPY ./. ./
RUN go build -o /build/${OS}-${ARCH}/ecotracker_exporter .

FROM golang:1.26-alpine

ARG ARCH="amd64"
ARG OS="linux"

RUN apk upgrade --no-cache

LABEL authors="koraktor"

COPY --from=builder /build/${OS}-${ARCH}/ecotracker_exporter /bin/ecotracker_exporter

EXPOSE 9776
ENTRYPOINT [ "/bin/sh", "-c", "/bin/ecotracker_exporter --host=${HOST} --port=${PORT} --listen-address=${LISTEN_ADDRESS:-:9776} --log-level=${LOG_LEVEL:-warn} ${*}", "entrypoint" ]
