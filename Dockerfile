FROM --platform=$TARGETPLATFORM go:1.21-bookworm as builder

# Set the working directory
WORKDIR /go/src/app

# Copy the Go source code into the image
COPY . .
WORKDIR /go/src/app/src

# We must fetch dependencies
RUN go get .

# Build the Go app for the target architecture (e.g., amd64)
RUN GOARCH=$TARGETARCH go build -o wait-for-mysql main.go

# Use a minimal Alpine image for the runtime
FROM --platform=$TARGETPLATFORM debian:bookworm
LABEL org.opencontainers.image.authors="AndersonPEM https://github.com/andersonpem"

# Copy the binary from the builder stage to /usr/local/bin
COPY --from=builder /go/src/app/src/wait-for-mysql /usr/local/bin/netknock

# Make the binary executable
RUN chmod +x /usr/local/bin/wait-for-mysql

# Run the binary
ENTRYPOINT ["/usr/local/bin/wait-for-mysql"]