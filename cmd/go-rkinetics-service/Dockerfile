FROM golang:1.24 AS builder

COPY ../../ /src
WORKDIR /src

# RUN GOPROXY=https://goproxy.cn make build
RUN make build

FROM debian:stable-slim

RUN apt-get update && apt-get install -y --no-install-recommends \
    ca-certificates  \
    netbase \
    && rm -rf /var/lib/apt/lists/ \
    && apt-get autoremove -y && apt-get autoclean -y

COPY --from=builder /src/bin /app
COPY ../../configs /app/configs

WORKDIR /app

EXPOSE 8000
EXPOSE 9000

CMD ["./server", "-conf", "/app/configs/config.yaml"]
