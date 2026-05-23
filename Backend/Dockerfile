# Build context = root repo (docker compose / CI: context . file Backend/Dockerfile)
FROM golang:1.22-alpine AS build
WORKDIR /src
COPY Backend/go.mod Backend/go.sum ./
RUN go mod download
COPY Backend/ .
RUN CGO_ENABLED=0 go build -o /api ./cmd/api

FROM alpine:3.20
RUN apk add --no-cache ca-certificates tzdata
WORKDIR /app
COPY --from=build /api /app/api
COPY Backend/migrations /app/migrations
COPY Frontend-admin/public/admin/_partials /app/partials
COPY Frontend-admin/static /app/static
COPY Frontend-admin/templates /app/templates
ENV ADDR=:8080 \
    ADMIN_PARTIALS_DIR=/app/partials \
    STATIC_DIR=/app/static \
    ADMIN_TEMPLATES_DIR=/app/templates \
    MIGRATIONS_DIR=/app/migrations
EXPOSE 8080
CMD ["/app/api"]
