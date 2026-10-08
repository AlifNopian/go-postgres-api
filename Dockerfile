# Stage 1: Build (Gunakan golang:alpine tanpa pin versi lama)
FROM golang:alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# Stage 2: Minimal Runtime
FROM alpine:latest

WORKDIR /app
COPY --from=builder /app/main .

EXPOSE 8081
CMD ["./main"]