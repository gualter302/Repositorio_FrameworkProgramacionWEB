FROM node:24-alpine AS frontend
WORKDIR /src/frontend
COPY frontend/package.json frontend/package-lock.json ./
RUN npm ci
COPY frontend/ ./
RUN npm run build

FROM golang:1.23-alpine AS backend
WORKDIR /src/backend
COPY backend/go.mod backend/go.sum ./
RUN go mod download
COPY backend/ ./
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/sigodental .

FROM alpine:3.21
RUN addgroup -S sigo && adduser -S sigo -G sigo
WORKDIR /app
COPY --from=backend /out/sigodental /app/sigodental
COPY --from=frontend /src/frontend/dist /app/static
ENV PORT=10000
ENV STATIC_DIR=/app/static
EXPOSE 10000
USER sigo
CMD ["/app/sigodental"]

