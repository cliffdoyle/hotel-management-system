# ---- Stage 1: Build ----
# Use the official Go image as the base for our build stage.
# The 'alpine' tag gives us a smaller base image.
FROM golang:1.21-alpine AS builder

# Set the working directory inside the container.
WORKDIR /app

# Copy the Go module files and download dependencies first.
# This step is cached, so Docker won't re-download dependencies on every build
# unless go.mod or go.sum has changed.
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application source code.
COPY . .

# Build the application.
# -o /app/server: specifies the output file path.
# -ldflags="-s -w": strips debugging information, making the binary smaller.
# CGO_ENABLED=0: disables CGO, ensuring a static binary.
RUN CGO_ENABLED=0 go build -o /app/server -ldflags="-s -w" ./cmd/api


# ---- Stage 2: Run ----
# Use a minimal, non-root base image for the final stage.
# 'distroless' is a Google project that contains only our app and its runtime dependencies.
# It does NOT contain a package manager, shell, or any other programs. Very secure.
FROM gcr.io/distroless/static-debian11 AS final

# Copy the compiled binary from the builder stage.
COPY --from=builder /app/server /server

# Copy the migrations folder.
# We need this if we want Fly.io to run our migrations on deploy.
COPY ./migrations /migrations

# Expose the port the application will listen on.
# This must match the PORT environment variable.
EXPOSE 8080

# The command to run when the container starts.
ENTRYPOINT ["/server"]