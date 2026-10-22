FROM golang:1.22-bookworm AS executor
WORKDIR /src
COPY go.mod ./
COPY *.go ./
RUN CGO_ENABLED=0 go build -o /executor . && CGO_ENABLED=0 go test -c -o /executor.test .

FROM debian:bookworm-slim AS isolate
RUN apt-get update \
    && apt-get install -y --no-install-recommends build-essential libcap-dev libsystemd-dev pkg-config git ca-certificates \
    && rm -rf /var/lib/apt/lists/*
ARG ISOLATE_VERSION=v2.0
RUN git clone --depth 1 --branch ${ISOLATE_VERSION} https://github.com/ioi/isolate /isolate \
    && make -C /isolate isolate

FROM debian:bookworm-slim AS base
RUN apt-get update \
    && apt-get install -y --no-install-recommends gcc g++ libc6-dev python3 ruby rustc default-jdk-headless libcap2 libsystemd0 \
    && rm -rf /var/lib/apt/lists/*
COPY --from=isolate /isolate/isolate /usr/local/bin/isolate
COPY --from=isolate /isolate/default.cf /usr/local/etc/isolate

# Runs the real-isolate tests: docker build --target integration, then run with --privileged.
FROM base AS integration
COPY --from=executor /executor.test /usr/local/bin/executor.test
ENV EXECUTOR_INTEGRATION=1
CMD ["executor.test", "-test.v", "-test.run", "Integration"]

FROM base
COPY --from=executor /executor /usr/local/bin/executor

# isolate creates namespaces and mounts, so the container needs --privileged (see docker-compose.yml).
EXPOSE 8001
CMD ["executor"]
