# 多架构构建：builder 固定在构建机平台，用 Go 交叉编译到目标平台
# （TARGETOS/TARGETARCH 由 buildx 注入），避免 QEMU 模拟 go build。
FROM --platform=$BUILDPLATFORM golang:1.26-alpine AS builder

ARG TARGETOS
ARG TARGETARCH

WORKDIR /app

ENV GOPROXY=https://goproxy.cn,direct
ENV CGO_ENABLED=0

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -ldflags="-s -w" -o /xyz .

FROM alpine:3.20

RUN apk add --no-cache ca-certificates tzdata

COPY --from=builder /xyz /usr/local/bin/xyz

EXPOSE 23020

CMD ["xyz"]