# Dockerfile
# Stage 1: Build the frontend
FROM node:22.14.0-slim AS frontend-builder
WORKDIR /app/frontend
COPY frontend/package.json frontend/package-lock.json* ./
# Use ci instead of install for potentially faster and more reliable builds in CI environments
RUN npm ci
COPY frontend/ ./
RUN npm run build

# Stage 2: Build the backend
FROM golang:1.23.8-alpine3.20 AS backend-builder
WORKDIR /app/backend
COPY backend/go.mod backend/go.sum ./
RUN go mod download
COPY backend/ ./
# Copy built frontend assets from the previous stage
# Ensure the target directory exists
RUN mkdir -p ./frontend/dist
COPY --from=frontend-builder /app/frontend/dist ./frontend/dist
# Build the Go application statically linked
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /app/server ./cmd/server

# Stage 3: Create the final image
FROM alpine:3.20.6
WORKDIR /app

# Create a non-root user and group
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

# Copy the built Go binary from the backend-builder stage
COPY --from=backend-builder /app/server .
# Copy the built frontend assets from the backend-builder stage
# The Go app expects them in ./frontend/dist relative to its execution path
COPY --from=backend-builder /app/backend/frontend/dist ./frontend/dist

# Change ownership to the non-root user
RUN chown -R appuser:appgroup /app

# Switch to the non-root user
USER appuser

# Expose the port the backend listens on
EXPOSE 8080

# Command to run the application
CMD ["/app/server"]
