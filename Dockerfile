# Use the official Alpine Linux image with Golang as the base for building
FROM --platform=$BUILDPLATFORM golang:alpine AS build

# Set the working directory inside the container
WORKDIR /app

# Copy the package files to the working directory
COPY . .

# Download the package dependencies
RUN go mod download

# Build the Golang package
RUN CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH go build -o main ./cmd/main.go

# Use Alpine Linux as the base for the final image
FROM --platform=$TARGETPLATFORM alpine:latest

# Set the working directory inside the container
WORKDIR /app

# Copy the built binary from the build stage
COPY --from=build /app/main .

# Expose port 3400
EXPOSE 3400

# Specify the command to run when the container starts
CMD ["./main"]