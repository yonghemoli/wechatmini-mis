# syntax=docker/dockerfile:1.7

ARG ALPINE_MIRROR=mirrors.aliyun.com
ARG GITHUB_REWRITE_BASE=github.com
ARG NODE_IMAGE=node:22-alpine
ARG GO_IMAGE=golang:1.24-alpine
ARG ALPINE_IMAGE=alpine:3.20

# ===== 阶段2: 后端构建 =====
FROM ${GO_IMAGE} AS builder
ARG ALPINE_MIRROR
ARG GITHUB_REWRITE_BASE
RUN sed -i "s/dl-cdn.alpinelinux.org/${ALPINE_MIRROR}/g" /etc/apk/repositories \
    && apk add --no-cache git ca-certificates \
    && update-ca-certificates
WORKDIR /app
ENV GOPROXY=https://goproxy.cn,https://goproxy.io,direct
ENV GOPRIVATE=github.com/xiuxianjs/*
COPY go.mod go.sum ./
RUN --mount=type=secret,id=github_token,required=false \
    if [ -s /run/secrets/github_token ]; then \
        TOKEN="$(cat /run/secrets/github_token)" && \
        GIT_CONFIG_COUNT=2 \
        GIT_CONFIG_KEY_0="url.https://${TOKEN}@${GITHUB_REWRITE_BASE}/.insteadOf" \
        GIT_CONFIG_VALUE_0=https://github.com/ \
        GIT_CONFIG_KEY_1=http.version \
        GIT_CONFIG_VALUE_1=HTTP/1.1 \
        go mod download; \
    else \
        GIT_CONFIG_COUNT=1 \
        GIT_CONFIG_KEY_0=http.version \
        GIT_CONFIG_VALUE_0=HTTP/1.1 \
        go mod download; \
    fi
COPY src ./src
COPY main.go ./
ARG TARGETOS
ARG TARGETARCH
ARG VERSION=0.0.1
RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} \
    go build -ldflags "-X main.Version=${VERSION} -X main.BuildTime=$(date +%s) -s -w" \
    -o yonghemolimis .

# ===== 阶段3: 最小运行镜像 =====
FROM ${ALPINE_IMAGE}
ARG ALPINE_MIRROR
WORKDIR /app
RUN sed -i "s/dl-cdn.alpinelinux.org/${ALPINE_MIRROR}/g" /etc/apk/repositories \
    && apk add --no-cache ca-certificates tzdata
COPY --from=builder /app/yonghemolimis .
CMD ["./yonghemolimis"]
