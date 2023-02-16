FROM golang:bullseye AS builder

RUN mkdir -p /tmp/metering/

COPY . /tmp/metering/

WORKDIR /tmp/metering/

RUN go build metering.go

FROM kong:3.0-ubuntu

USER root

RUN mkdir -p /opt/amberflo

COPY --from=builder  /tmp/metering/metering /opt/amberflo
COPY kong.conf /etc/kong/

USER kong
