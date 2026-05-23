FROM golang:1.22-alpine AS builder

WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./
RUN go mod download
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/expense-tracker ./cmd/api/main.go

FROM alpine:latest
WORKDIR /app

# Install CA certificates for AWS S3 and other HTTPS requests
RUN apk --no-cache add ca-certificates

# Copy the binary from the builder stage
COPY --from=builder /app/expense-tracker .

EXPOSE 8080
CMD ["./expense-tracker"]
