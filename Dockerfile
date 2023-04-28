# Stage 1: Build
FROM golang:1.20-alpine AS build

# Set the working directory inside the container
WORKDIR /app

# Copy files from the host machine to the container
COPY . .

# Compile the main.go file and create an executable
RUN go build -o portService main.go

# Stage 2: Production
FROM alpine:3.17

# Create a non-root user and group
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

# Set the working directory inside the container
WORKDIR /app

# Copy the compiled binary from the build stage and the ports.json file
COPY --from=build /app/portService /app/portService
COPY ports.json /app/ports.json

# Set permissions and change ownership to the non-root user
RUN chown -R appuser:appgroup /app

# Run the container as the non-root user
USER appuser

# Run the application with the specified configuration file
CMD ["/app/portService", "ports.json"]
