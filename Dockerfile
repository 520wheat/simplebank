# 阶段 1：编译
FROM golang:1.26-alpine AS builder

RUN apk add --no-cache git

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o /server ./cmd/server

# 阶段 2：运行
FROM alpine:3.21
RUN apk add --no-cache curl

WORKDIR /app
COPY --from=builder /server .
COPY app.env.example ./app.env
COPY db/migration ./db/migration

EXPOSE 8080 9090
HEALTHCHECK --interval=30s --timeout=3s CMD curl -f http://localhost:8080/health || exit 1

ENTRYPOINT ["/app/server"]