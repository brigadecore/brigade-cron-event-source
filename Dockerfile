FROM --platform=$BUILDPLATFORM brigadecore/go-tools:v0.6.0 as builder

ARG VERSION
ARG COMMIT
ARG TARGETOS
ARG TARGETARCH
ENV CGO_ENABLED=0

WORKDIR /src
COPY *.go .
COPY go.mod go.mod
COPY go.sum go.sum

RUN GOOS=$TARGETOS GOARCH=$TARGETARCH go build \
  -o bin/cron-event-source \
  -ldflags "-w -X github.com/brigadecore/brigade-foundations/version.version=$VERSION -X github.com/brigadecore/brigade-foundations/version.commit=$COMMIT" \
  .

FROM gcr.io/distroless/static:nonroot as final
COPY --from=builder /src/bin/ /brigade-cron-event-source/bin/
ENTRYPOINT ["/brigade-cron-event-source/bin/cron-event-source"]