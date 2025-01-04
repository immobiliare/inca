FROM golang:alpine3.21 as builder
ENV GOOS=linux
ENV GOARCH=amd64
WORKDIR /workspace
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN go build

FROM alpine:3.21
WORKDIR /tmp
COPY --from=builder /workspace/inca /usr/sbin/
RUN mkdir -p /tmp/server/webroot
ENTRYPOINT ["/usr/sbin/inca"]
LABEL org.opencontainers.image.source=https://github.com/immobiliare/inca
