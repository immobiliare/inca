FROM cgr.dev/chainguard/go:latest AS builder
ENV CGO_ENABLED=0
RUN mkdir -p /tmp/server/webroot
WORKDIR /workspace
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN go build

FROM ghcr.io/anchore/syft:latest AS sbomgen
COPY --from=builder /workspace/inca /usr/sbin/inca
RUN ["/syft", "--output", "spdx-json=/inca.spdx.json", "/usr/sbin/inca"]

FROM cgr.dev/chainguard/static:latest
WORKDIR /tmp
COPY --from=builder /workspace/inca /usr/sbin/
COPY --from=builder /tmp/server /tmp/server
COPY --from=sbomgen /inca.spdx.json /var/lib/db/sbom/inca.spdx.json
ENTRYPOINT ["/usr/sbin/inca"]
LABEL org.opencontainers.image.title="inca"
LABEL org.opencontainers.image.description="INternal CA is an API around Certificate Authority flows to handle internal and global certificates at ease"
LABEL org.opencontainers.image.url="https://github.com/immobiliare/inca"
LABEL org.opencontainers.image.source="https://github.com/immobiliare/inca"
LABEL org.opencontainers.image.licenses="MIT"
LABEL io.containers.autoupdate=registry
