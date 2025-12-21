# Multi-stage build for optimized Lambda image

# Build stage
FROM golang:1.23-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git make

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-s -w" \
    -o bootstrap \
    ./cmd/lambda

# Runtime stage
FROM public.ecr.aws/lambda/provided:al2

# Copy the binary from builder
COPY --from=builder /app/bootstrap ${LAMBDA_RUNTIME_DIR}/bootstrap

# Set the CMD to your handler
CMD [ "bootstrap" ]

