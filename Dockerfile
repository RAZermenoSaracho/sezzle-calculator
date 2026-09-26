# Build the frontend.
FROM node:24-alpine AS frontend
WORKDIR /src/frontend
COPY frontend/package.json frontend/package-lock.json ./
RUN npm ci
COPY frontend/ ./
RUN npm run build

# Build a static Go binary.
FROM golang:1.27.1-alpine AS backend
WORKDIR /src/backend
COPY backend/ ./
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/server ./cmd/server

# Single runtime image: the Go server serves /api and the built frontend.
FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=backend /out/server /app/server
COPY --from=frontend /src/frontend/dist /app/static
ENV PORT=8080 STATIC_DIR=/app/static
EXPOSE 8080
USER nonroot:nonroot
ENTRYPOINT ["/app/server"]
