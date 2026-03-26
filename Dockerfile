FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o routeapi ./cmd/server

FROM alpine:3.19

RUN apk --no-cache add ca-certificates

WORKDIR /opt/app

COPY --from=builder /app/routeapi .
COPY web/ ./web/

EXPOSE 8080

CMD ["./routeapi"]
