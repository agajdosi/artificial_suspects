FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY ./backend ./backend
RUN go build -o /artsus_server backend/main.go

FROM alpine:latest
COPY --from=builder /artsus_server /artsus_server
EXPOSE 8080
CMD ["/artsus_server", "-db-path", "/data/artsus.db", "-host", "0.0.0.0"]
