# Stage 1: build the Go binary
FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o blocky-visor-sidecar .

# Stage 2: minimal runtime image
FROM alpine:3.21
RUN apk add --no-cache ca-certificates
WORKDIR /app
COPY --from=builder /app/blocky-visor-sidecar .
EXPOSE 8550
ENTRYPOINT ["./blocky-visor-sidecar"]
CMD ["-config", "/etc/blocky-visor/config.yaml"]
