# GeoIP Go Server

A high-performance HTTP server built with Go and `fasthttp` that provides IP geolocation and ASN information using MaxMind GeoLite2 databases.

## Features

- **Fast**: Built on `valyala/fasthttp` for high throughput.
- **Geolocation**: Resolves IP addresses to Country, City, Latitude, Longitude, and Timezone.
- **Network Info**: Provides ASN and ISP details (requires ASN database).
- **Dockerized**: Includes a Dockerfile that automatically downloads the latest databases during the build.

## API Endpoints

### IP Lookup

- **URL**: `/lookup/<ip_address>`
- **Method**: `GET`
- **Example**: `GET /lookup/8.8.8.8`

**Response (JSON):**

```json
{
  "ip": "8.8.8.8",
  "code": "US",
  "country": "United States",
  "city": "Glenmont",
  "lat": 40.5369,
  "lon": -82.1286,
  "tz": "America/New_York",
  "asn": 15169,
  "isp": "Google LLC"
}
```

## Running Locally

1. **Prerequisites**:

   - Go 1.21+
   - `GeoLite2-City.mmdb` (Required)
   - `GeoLite2-ASN.mmdb` (Optional, for ISP info)

   Place the `.mmdb` files in the project root.

2. **Run**:
   ```bash
   go run main.go
   ```
   Server starts on `http://localhost:8080`.

## Running with Docker

The Docker build process automatically downloads the latest MaxMind databases. You need a MaxMind License Key (free).

1. **Build the image**:

   ```bash
    docker build --build-arg MAXMIND_LICENSE_KEY=YOUR_LICENSE_KEY -t geoip-server .
   ```

2. **Run the container**:
   ```bash
    docker run -p 8080:8080 geoip-server
   ```
