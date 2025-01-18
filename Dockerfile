FROM golang:1.21.0 AS builder
WORKDIR /app
COPY go.mod ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o ./cloudflare-dyndns cmd/cloudflare-dyndns/cloudflare-dyndns.go 

FROM scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /app/cloudflare-dyndns /cloudflare-dyndns
EXPOSE 8000
ENTRYPOINT ["/cloudflare-dyndns"]

