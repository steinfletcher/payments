FROM golang:1.12 as builder
COPY . /payments
WORKDIR /payments
ENV GO111MODULE=on
RUN CGO_ENABLED=0 GOOS=linux go build -o payments cmd/payments/main.go

FROM alpine:3.9
WORKDIR /root/
RUN apk add --no-cache tzdata
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /payments .
RUN ls /root/
CMD ["./payments"]
