# Use an official Go runtime as the base image
FROM golang:1.20

# Set the working directory in the container
WORKDIR /src

# Copy the local code to the container
COPY ./src /src
COPY ./graphs /graphs

# Create output_graphs directory
RUN mkdir /output_graphs

# Build the Go application
RUN go build -o pregel .

# Run the application
CMD ["./pregel"]