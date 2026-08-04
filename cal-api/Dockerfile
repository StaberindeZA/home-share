# Stage 1: Build the binary
FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /bin/main cmd/server/main.go

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /bin/init cmd/db/init.go

FROM alpine:latest AS initializer
COPY --from=builder /bin/init .
CMD ["./init"]

# Target 2: Main Service
FROM alpine:latest AS app
#RUN adduser -D appuser
#USER appuser
WORKDIR /app
COPY --from=builder /bin/main .
EXPOSE 8080
CMD ["./main"]
