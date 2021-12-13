FROM --platform=$BUILDPLATFORM brigadecore/go-tools:v0.5.0 as builder

ARG VERSION
ARG COMMIT
ARG TARGETOS
ARG TARGETARCH
ENV CGO_ENABLED=0

WORKDIR /src
COPY main.go main.go
COPY config.go config.go
COPY go.mod go.mod
COPY go.sum go.sum

RUN GOOS=$TARGETOS GOARCH=$TARGETARCH go build \
  -o bin/cron-gateway \
  -ldflags "-w -X github.com/brigadecore/brigade-foundations/version.version=$VERSION -X github.com/brigadecore/brigade-foundations/version.commit=$COMMIT" \
  .

FROM scratch
COPY --from=builder /src/bin/ /brigade-cron-gateway/bin/
ENTRYPOINT ["/brigade-cron-gateway/bin/cron-gateway"]