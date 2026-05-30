FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN sed -i 's/go 1.26.1/go 1.23/' go.mod && go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /app/bin/api-gateway ./api-gateway/main.go

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/bin/api-gateway .
EXPOSE 8080
CMD ["./api-gateway"]
