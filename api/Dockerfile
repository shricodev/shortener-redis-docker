# Build the application using the golang image
# Stage: builder
FROM golang:1.22.0-alpine3.18 as builder

# Set the working directory to /build
WORKDIR /build

# Copy the application source code into the container
COPY . /build/

# Build the application
RUN go build -o main .

# Use the alpine image to run the application
FROM alpine:3.18

# Create a non-root user for running the application
RUN adduser -S -D -H -h /app appuser

# Switch to the non-root user
USER appuser

# Set the working directory to /app
WORKDIR /app

# Copy the application source code into the container
COPY . /app

# Copy the built application from the builder stage
COPY --from=builder /build/main /app/

# Expose port  3000 for the application
EXPOSE  3000

# Command to run the application
CMD ["./main"]
