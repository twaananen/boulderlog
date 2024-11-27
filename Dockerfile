# Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Install Node.js, npm, and other dependencies
RUN apk add --no-cache nodejs npm git

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Install templ
RUN go install github.com/a-h/templ/cmd/templ@latest

# Copy the source code
COPY . .

# Install Tailwind CSS
RUN npm install -g tailwindcss

# Generate Tailwind CSS
RUN tailwindcss -i ./static/css/input.css -o ./static/css/tailwind.css --minify

# Generate templ files
RUN templ generate

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Final stage
FROM alpine:latest

WORKDIR /root/

# Copy the binary from the builder stage
COPY --from=builder /app/main .
COPY --from=builder /app/static ./static

# Expose port 8080
EXPOSE 8080

# Command to run the executable
CMD ["./main"]