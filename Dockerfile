# Build the manager binary
FROM golang:1 AS builder

ARG CHIA_IMAGE_TAG=latest
ARG EXPORTER_IMAGE_TAG=latest
ARG HEALTHCHECK_IMAGE_TAG=latest

ENV CHIA_IMAGE_TAG=${CHIA_IMAGE_TAG}
ENV EXPORTER_IMAGE_TAG=${EXPORTER_IMAGE_TAG}
ENV HEALTHCHECK_IMAGE_TAG=${HEALTHCHECK_IMAGE_TAG}

WORKDIR /workspace
# Copy the Go Modules manifests
COPY go.mod go.mod
COPY go.sum go.sum
# cache deps before building and copying source so that we don't need to re-download as much
# and so that source changes don't invalidate our downloaded layer
RUN go mod download

# Copy the go source
COPY Makefile Makefile
COPY cmd/main.go cmd/main.go
COPY api/ api/
COPY internal/ internal/
COPY hack/ hack/

RUN make build

# Use distroless as minimal base image to package the manager binary
# Refer to https://github.com/GoogleContainerTools/distroless for more details
FROM gcr.io/distroless/static:nonroot
WORKDIR /
COPY --from=builder /workspace/bin/manager .
USER 65532:65532

ENTRYPOINT ["/manager"]
