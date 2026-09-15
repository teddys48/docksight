# Stage 1: Build Svelte Frontend with Bun
FROM oven/bun:1 AS frontend-builder
WORKDIR /app/frontend

COPY frontend/package.json frontend/bun.lock* ./
RUN bun install --frozen-lockfile || bun install

COPY frontend/ ./
RUN bun run build

# Stage 2: Build Go Backend with Embedded Frontend
FROM golang:alpine AS backend-builder
WORKDIR /app/backend

COPY backend/ ./
# Copy built static frontend files into backend/dist for go:embed
COPY --from=frontend-builder /app/frontend/dist ./dist

RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o docksight main.go

# Stage 3: Minimal Final Runtime
FROM alpine:3.20

RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app

# Copy binary from backend builder
COPY --from=backend-builder /app/backend/docksight /app/docksight

# Create data directory for SQLite database
RUN mkdir -p /data
ENV DB_PATH=/data/monitoring.db
ENV PORT=8080

EXPOSE 8080

ENTRYPOINT ["/app/docksight"]
