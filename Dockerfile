# Build frontend stage
FROM node:20-alpine AS frontend-builder

WORKDIR /app

# Copy frontend files
COPY frontend/package*.json ./
RUN npm install

COPY frontend/ ./
RUN npm run build

# Build backend stage
FROM golang:1.24-alpine AS backend-builder

WORKDIR /app

# Install dependencies
RUN apk add --no-cache git

# Copy go mod files
COPY backend/go.mod backend/go.sum ./
COPY backend/go.mod backend/go.sum ./backend/
RUN go mod download

# Copy source code
COPY backend/ ./backend/

# Build the application
WORKDIR /app/backend
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o doodle-backend main.go

# Final stage
FROM alpine:latest

WORKDIR /app

# Install ca-certificates for HTTPS
RUN apk --no-cache add ca-certificates

# Copy the binary from backend-builder
COPY --from=backend-builder /app/backend/doodle-backend /app/doodle-backend

# Copy frontend build from frontend-builder
COPY --from=frontend-builder /app/dist /app/frontend

# Expose port
EXPOSE 8080

# Run the application
CMD ["/app/doodle-backend"]
