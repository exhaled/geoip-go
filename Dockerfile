# Build stage
FROM golang:alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o geoip-server .

# Downloader stage
FROM alpine:latest AS downloader

RUN apk add --no-cache curl tar

ARG MAXMIND_LICENSE_KEY

WORKDIR /downloads

# Download and extract GeoLite2-City
RUN if [ -z "$MAXMIND_LICENSE_KEY" ]; then \
      echo "Error: MAXMIND_LICENSE_KEY build argument is required"; \
      exit 1; \
    fi && \
    curl -L "https://download.maxmind.com/app/geoip_download?edition_id=GeoLite2-City&license_key=${MAXMIND_LICENSE_KEY}&suffix=tar.gz" -o City.tar.gz && \
    tar -xzf City.tar.gz && \
    find . -name "GeoLite2-City.mmdb" -exec mv {} /downloads/ \;

# Download and extract GeoLite2-ASN
RUN curl -L "https://download.maxmind.com/app/geoip_download?edition_id=GeoLite2-ASN&license_key=${MAXMIND_LICENSE_KEY}&suffix=tar.gz" -o ASN.tar.gz && \
    tar -xzf ASN.tar.gz && \
    find . -name "GeoLite2-ASN.mmdb" -exec mv {} /downloads/ \;

# Final stage
FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/geoip-server .
COPY --from=downloader /downloads/GeoLite2-City.mmdb .
COPY --from=downloader /downloads/GeoLite2-ASN.mmdb .

EXPOSE 8080

CMD ["./geoip-server"]

