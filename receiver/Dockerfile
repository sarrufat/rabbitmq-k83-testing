# Build the manager binary
FROM golang:1.13 as builder

WORKDIR /workspace
# Copy the go source
COPY worker.go worker.go

# Copy the Go Modules manifests
# COPY go.mod go.mod
# COPY go.sum go.sum


# Build
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GO111MODULE=on go build -a -o worker worker.go

FROM ubuntu:latest

ARG GIT_COMMIT
LABEL GitCommit=$GIT_COMMIT

WORKDIR /
COPY --from=builder /workspace/worker .

# Create operator system user & group
RUN set -eux; \
	groupadd --gid 1000 --system rabbitmq-cluster; \
	useradd --uid 1000 --system --gid rabbitmq-cluster rabbitmq-cluster

USER 1000:1000

ENTRYPOINT ["/worker"]
