# Use the official Golang image as the base image
FROM golang:latest

# Set the working directory inside the container
WORKDIR /app

# Copy files from the host machine to the container
COPY . .

# Compile the main.go file and create an executable
RUN go build -o portService main.go

# Run the test inside the container
CMD ["/app/portService", "ports.json"]
