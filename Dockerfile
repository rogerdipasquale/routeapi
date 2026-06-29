FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN apk --no-cache add make
RUN make docs-install \ 
    && make docs
RUN make build

FROM alpine:3.19

RUN apk --no-cache add ca-certificates

WORKDIR /opt/app

COPY --from=builder /app/routeapi .
COPY --from=builder docs/ ./
COPY web/ ./web/

EXPOSE 8080

CMD ["./routeapi"]
