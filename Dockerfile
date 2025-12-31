# Build stage
FROM golang:1.23.1 AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy go mod files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application
COPY . .

# Build the Go binary for Linux with Lambda compatibility
RUN go build -tags lambda.norpc -o main main.go

# Runtime stage using the AWS Lambda base image for Go
FROM public.ecr.aws/lambda/provided:al2023

# Copy the binary from the builder stage
COPY --from=builder /app/main ./main

ENTRYPOINT [ "./main" ]