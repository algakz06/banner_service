# Use the official Golang image to create a build artifact.
# This is based on Debian and includes the Go toolset.
FROM golang:1.22 as builder

# Set the Current Working Directory inside the container
WORKDIR /app

# install psql
RUN apt-get update
RUN apt-get -y install postgresql-client

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source code into the container
COPY . .

# make wait-for-postgres.sh executable
RUN chmod +x wait-for-postgres.sh

# Build the Go app
RUN go build ./cmd/main.go

# Command to run the executable
CMD ["./main"]
