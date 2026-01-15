# Use the official Golang 1.22 as the base image
FROM golang:1.22-alpine

# Set the Current Working Directory inside the container
WORKDIR /usr/src/app

RUN apt update -yq
RUN apt install -y postgresql-client

# Download Go modules
COPY go.mod go.sum ./
RUN go mod download && go mod verify

# Copy the source code into the container
COPY . .

# Build the binary
RUN go build -o /usr/local/bin/app ./cmd

# Expose the port on which the application will run
EXPOSE 8080

# Command to run the executable
CMD ["/usr/local/bin/app"]
