# Use the official Golang image as the base
FROM golang:latest

# Set the working directory inside the container
WORKDIR /app

# Copy the package files to the working directory
COPY . .

# Download the package dependencies
RUN go mod download

# Build the Golang package
RUN go build -o main ./cmd/main.go

# Expose port 3400
EXPOSE 3400

# Specify the command to run when the container starts
CMD ["./main"]