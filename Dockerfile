FROM golang:1.24-alpine AS builder

WORKDIR /app

# Swag CLI kurulumu
RUN go install github.com/swaggo/swag/cmd/swag@latest

COPY go.mod go.sum ./
ENV GOTOOLCHAIN=auto
RUN go mod download

COPY . .

# Swagger dokümantasyonu oluştur
RUN swag init -g cmd/server/main.go -o docs

RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/server/main.go

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/main .

EXPOSE 8080

CMD ["./main"]