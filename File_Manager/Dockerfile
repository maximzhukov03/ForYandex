FROM golang:1.23-alpine AS builder
WORKDIR /app

COPY go.mod ./
RUn go mod download

COPY . .

RUN go build -o server ./cmd/main.go

FROM alpine:latest
WORKDIR /app

COPY --from=builder /app/server .
COPY site ./site

RUN mkdir -p ./storage

EXPOSE 8080
CMD ["./server"]