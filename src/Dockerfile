FROM golang:1.16-alpine

# Set the Current Working Directory inside the container
WORKDIR /versioncontrol-service

# Copy go mod file
COPY go.mod ./

# Download all dependencies.
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

EXPOSE 8080

# Build the Go app
RUN go build -o versioncontrol ./src/

# Command to run the executable
ENTRYPOINT [ "/versioncontrol-service/versioncontrol" ]