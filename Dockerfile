FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod ./
COPY go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application
# Using 'app' as a default name - replace with your actual application name if different
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/app

# Use a smaller image for the final application
FROM alpine:3.19

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the binary from the builder stage
COPY --from=builder /app/main .

# Copy necessary config files if needed
COPY --from=builder /app/configs ./configs

# Expose the port your application runs on (update as needed)
EXPOSE 8080

# Command to run the application
CMD ["./main"]
