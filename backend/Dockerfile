# Build stage
FROM golang:1.27-alpine AS build
WORKDIR /app
COPY go.mod go.sum* ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /bin/api ./cmd/api

# Runtime stage
FROM alpine:3.21
COPY --from=build /bin/api /api
EXPOSE 8000
ENTRYPOINT ["/api"]
