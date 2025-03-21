FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o /artsus_server main.go

FROM alpine:latest
COPY --from=builder /artsus_server /artsus_server
VOLUME [ "/data" ]
EXPOSE 8080
CMD ["/artsus_server", "-db-path", "/data/artsus.db", "-host", "0.0.0.0"]
