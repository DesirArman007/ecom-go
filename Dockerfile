# Stage 1: use the Go image only to compile the application.
# This stage will not be included in the final runtime image.
FROM golang:1.25.4 AS builder

# /src is the working directory inside the temporary build container.
WORKDIR /src

# Copy dependency files first so Docker can reuse this layer when only
# application source code changes.
COPY go.mod go.sum  ./
RUN go mod download

# Copy the rest of the project, including cmd/ and internal/.
COPY . .

# Your executable Go package is in cmd/, so build ./cmd explicitly.
# CGO_ENABLED=0 creates a self-contained binary suitable for Alpine.
# GOOS=linux makes the target operating system explicit.
RUN CGO_ENABLED=0 GOOS=linux go build -o /out/server ./cmd

# Stage 2: create a small image containing only the compiled application.
FROM alpine:3.20

# /app is the working directory in the final runtime container.
WORKDIR /app

# Copy only the compiled binary from the builder stage.
COPY --from=builder /out/server ./server

# The Go server listens on port 8080 (see cmd/main.go).
EXPOSE 8080

# Start the compiled binary when the container starts.
CMD ["./server"]
