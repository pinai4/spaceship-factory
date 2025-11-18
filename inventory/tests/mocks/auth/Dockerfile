FROM golang:1.24.2-alpine AS builder

RUN apk add --no-cache git

WORKDIR /app/inventory

COPY platform/go.mod platform/go.sum /app/platform/
COPY shared/go.mod shared/go.sum /app/shared/
COPY inventory/go.mod inventory/go.sum ./

RUN go mod download

COPY platform /app/platform
COPY shared /app/shared
COPY inventory ./

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o mock-auth ./tests/mocks/auth/main.go

FROM alpine:3.21.3

RUN addgroup -S appgroup && adduser -S appuser -G appgroup

WORKDIR /app

USER appuser

COPY --from=builder /app/inventory/mock-auth .

EXPOSE 50051

ENTRYPOINT ["./mock-auth"]