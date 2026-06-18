FROM golang:1.21-alpine AS builder

WORKDIR /usr/src/app

# Install git if needed for dependencies
RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /usr/local/bin/server cmd/server/main.go

FROM alpine:3.19 AS runner

WORKDIR /app

COPY --from=builder /usr/local/bin/server /usr/local/bin/server

EXPOSE 8080

CMD ["/usr/local/bin/server"]
