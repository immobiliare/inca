FROM golang:alpine3.14 as builder
ENV GOOS=linux
ENV GOARCH=amd64
WORKDIR /workspace
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN go build

FROM alpine:3.14
WORKDIR /tmp
COPY --from=builder /workspace/inca /usr/sbin/
CMD ["/usr/sbin/inca"]