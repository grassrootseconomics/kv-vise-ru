FROM golang:1.24-bookworm AS build

ENV CGO_ENABLED=0

ARG BUILDPLATFORM
ARG TARGETPLATFORM
ARG BUILD=dev

RUN echo "Building on $BUILDPLATFORM, building for $TARGETPLATFORM"
WORKDIR /build

COPY . .
RUN go mod download
RUN go build -o kv-vise-ru -ldflags="-X main.build=${BUILD} -s -w" cmd/main.go

FROM debian:bookworm-slim

ENV DEBIAN_FRONTEND=noninteractive

WORKDIR /service

COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /build .

EXPOSE 5059

CMD ["./kv-vise-ru"]