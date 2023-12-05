FROM golang:1.21.1 as builder

WORKDIR /app

COPY go.* ./
RUN go mod download

COPY . ./

RUN CGO_ENABLED=0 GOOS=linux go build -v -o main ./

FROM scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /app/main /app/main
WORKDIR /app
EXPOSE 8933
EXPOSE 161
CMD ["./main"]