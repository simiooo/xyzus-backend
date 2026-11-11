FROM golang:1.26-alpine AS builder

WORKDIR /app

ENV GOPROXY=https://goproxy.cn,direct
ENV CGO_ENABLED=0

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -ldflags="-s -w" -o /xyz .

FROM alpine:3.20

RUN apk add --no-cache ca-certificates tzdata

COPY --from=builder /xyz /usr/local/bin/xyz

EXPOSE 23020

CMD ["xyz"]