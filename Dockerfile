# Use the official Go image from DockerHub as the base image
FROM golang:1.21.1 as builder

# Set the current working directory inside the container
WORKDIR /app

# Copy go mod and sum files to the workspace
COPY go.mod go.sum ./

# Download all dependencies. 
RUN go mod download

# Copy the source code from the current directory to the Working Directory inside the container
COPY . .

# Build the application. We'll use CGO_ENABLED=0 to ensure the binary is statically linked.
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o main .

### Start a new stage from scratch ###
FROM alpine:latest  

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the binary from the builder stage
COPY --from=builder /app/main .

# Command to run the executable
CMD ["./main"] 
