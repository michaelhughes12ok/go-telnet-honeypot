FROM golang:1.27 AS builder
WORKDIR /app
COPY . .
RUN go build -o honeypot .


FROM debian:bookworm-slim
WORKDIR /app
COPY --from=builder /app/honeypot .
EXPOSE 23
CMD ["./honeypot"]