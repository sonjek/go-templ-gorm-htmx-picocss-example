FROM golang:1.25-alpine3.22 AS builder
WORKDIR /app
RUN apk --update add --no-cache ca-certificates tzdata make && update-ca-certificates
COPY . .
RUN export GOPATH="/root/go" && make build

# Pach binary file
RUN apk add upx
RUN upx bin/app

FROM scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo/
COPY --from=builder /app/bin/app /

EXPOSE 8089

ENTRYPOINT ["/app"]
