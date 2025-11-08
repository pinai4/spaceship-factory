# ============================
# Stage 1: Build stage
# ============================

# Use the official Go image based on Alpine — lightweight, fast, and secure
FROM golang:1.24.2-alpine AS builder

# Install git — required for downloading dependencies from private and public repositories
RUN apk add --no-cache git

# Set the working directory inside the container
WORKDIR /app

# Copy go.work and its checksum — this allows Go to properly link the modules
COPY go.work ./

# Copy go.mod and go.sum of all dependent modules — needed for caching go mod download
COPY platform/go.mod platform/go.sum ./platform/
COPY shared/go.mod shared/go.sum ./shared/
COPY payment/go.mod payment/go.sum ./payment/
COPY order/go.mod order/go.sum ./order/
COPY inventory/go.mod inventory/go.sum ./inventory/
COPY assembly/go.mod assembly/go.sum ./assembly/
COPY notification/go.mod notification/go.sum ./notification/

# Download all dependencies listed in the module files
RUN go mod download

# Copy the source code of all modules — including dependencies, as they are needed during compilation
COPY platform ./platform
COPY shared ./shared
COPY order ./order
COPY payment ./payment
COPY inventory ./inventory
COPY assembly ./assembly
COPY notification ./notification

## Download grpc-health-probe — a utility for checking the health of a gRPC service
#ADD https://github.com/grpc-ecosystem/grpc-health-probe/releases/download/v0.4.37/grpc_health_probe-linux-amd64 ./grpc-health-probe
#
## Make the downloaded file executable
#RUN chmod +x grpc-health-probe

# Build the binary file of the UFO service for the Linux architecture without CGO
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o app-inventory ./inventory/cmd/main.go


# ============================
# Stage 2: Final image
# ============================

# Use a clean Alpine image as the final container — lightweight and secure
FROM alpine:3.21.3

# Create a system user without root privileges — best practice for security
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

# Set the working directory inside the container
WORKDIR /app

# Run the application under a non-root user
USER appuser

# Copy the compiled binary from the builder stage
COPY --from=builder /app/app-inventory .

## Copy the grpc-health-probe binary to /bin
#COPY --from=builder /app/grpc-health-probe /bin/grpc-health-probe

# Copy the .env file for application configuration
COPY deploy/compose/inventory/.env ./deploy/compose/inventory/.env

# Expose the gRPC service port of the UFO
EXPOSE 50051

# Set the startup command — run our binary
ENTRYPOINT ["./app-inventory"]
