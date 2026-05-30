FROM node:20-alpine AS frontend-builder
WORKDIR /app/frontend
COPY frontend/package.json frontend/package-lock.json* ./
RUN npm install
COPY frontend/ .
RUN npm run build

FROM golang:alpine AS backend-builder
RUN apk add --no-cache gcc musl-dev
WORKDIR /app/backend
COPY backend/go.mod backend/go.sum ./
RUN go mod download
COPY backend/ .
COPY --from=frontend-builder /app/frontend/dist ./cmd/dist
RUN CGO_ENABLED=1 CGO_CFLAGS="-D_LARGEFILE64_SOURCE" go build -o /app/k8s-ui-admin ./cmd/...

FROM alpine:3.19
RUN apk add --no-cache ca-certificates tzdata
WORKDIR /app
COPY --from=backend-builder /app/k8s-ui-admin .
EXPOSE 8080
VOLUME ["/app/data"]
ENV PORT=8080
ENV DB_PATH=/app/data/k8s_ui_admin.db
ENV GIN_MODE=release
CMD ["./k8s-ui-admin"]
