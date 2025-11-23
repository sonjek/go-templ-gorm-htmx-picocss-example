FROM golang:1.25 AS builder
WORKDIR /app
RUN apt update && apt install -y --no-install-recommends ca-certificates tzdata make curl upx
COPY --from=oven/bun:1.3 /usr/local/bin/bun /usr/local/bin/bun
COPY . .
RUN make build

# Pack binary file
RUN upx bin/app

FROM scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo/
COPY --from=builder /app/bin/app /

EXPOSE 8089

ENTRYPOINT ["/app"]
