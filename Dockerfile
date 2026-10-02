# syntax=docker/dockerfile:1

# ---- frontend build ----
FROM node:20-alpine AS frontend-build
WORKDIR /app/frontend
COPY frontend/package.json frontend/package-lock.json ./
RUN npm ci
COPY frontend/ ./
RUN npm run build

# ---- backend build ----
FROM golang:1.23-alpine AS backend-build
WORKDIR /app/backend
COPY backend/go.mod ./
RUN go mod download
COPY backend/ ./
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /server ./cmd/server

# ---- final ----
FROM alpine:3.20
RUN addgroup -S app && adduser -S -G app app
WORKDIR /app
COPY --from=backend-build /server /app/server
COPY --from=frontend-build /app/frontend/dist /app/static

ENV PORT=8080
ENV STATIC_DIR=/app/static
EXPOSE 8080

USER app:app
ENTRYPOINT ["/app/server"]
