FROM golang:1.22.5 as builder
WORKDIR /workspace
ENV GOPROXY=https://goproxy.cn,direct
COPY . ./
RUN go mod download
RUN CGO_ENABLED=0 go build -a -o consumer .
RUN CGO_ENABLED=0 go build -a -o consume-cpu/consume-cpu ./consume-cpu

FROM debian:bookworm-20240722

RUN set -ex \
    && apt update \
    && apt install stress -y \
    && apt install procps -y

ARG TARGETARCH
WORKDIR /
COPY --from=builder /workspace/consumer /consumer
COPY --from=builder /workspace/consume-cpu/consume-cpu /consume-cpu/consume-cpu

EXPOSE 9100
ENTRYPOINT ["/consumer"]
