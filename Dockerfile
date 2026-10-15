FROM golang:1.26.9 AS builder
WORKDIR /app
COPY go.mod ./

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o honeypot . 

FROM gcr.io/distroless/static:nonroot
COPY --from=builder /app/honeypot /honeypot
EXPOSE 2323
ENTRYPOINT ["/honeypot"]
