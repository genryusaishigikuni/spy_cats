# Stage 1: Build the Go application
FROM golang:1.24.1-alpine as builder

WORKDIR /app
COPY . .

# Install dependencies and build the Go app
RUN go mod tidy
RUN go build -o spy_cats cmd/main.go

# Stage 2: Final image that will run both Go app and PostgreSQL
FROM alpine:latest

# Install PostgreSQL client utilities to get `pg_isready`
RUN apk update && apk add --no-cache postgresql-client bash

# Create a user for the application
RUN adduser -D spycat
USER spycat

# Set the working directory
WORKDIR /app

# Copy the Go binary from the builder stage
COPY --from=builder /app/spy_cats .

# Expose the port the app will run on
EXPOSE 8080

# Run the app
CMD ["sh", "-c", "until pg_isready -h db -p 5432; do echo waiting for db; sleep 2; done; ./spy_cats"]
