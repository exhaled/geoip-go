package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"strings"

	"github.com/oschwald/geoip2-golang"
	"github.com/valyala/fasthttp"
)

func main() {
	// Open the GeoLite2-City database
	cityDB, err := geoip2.Open("GeoLite2-City.mmdb")
	if err != nil {
		log.Fatal(err)
	}
	defer cityDB.Close()

	// Try to open the GeoLite2-ASN database (optional but needed for ISP/ASN info)
	asnDB, err := geoip2.Open("GeoLite2-ASN.mmdb")
	if err != nil {
		log.Printf("Warning: GeoLite2-ASN.mmdb not found. ISP and ASN info will be missing.")
	} else {
		defer asnDB.Close()
	}

	writeJSONError := func(ctx *fasthttp.RequestCtx, message string, statusCode int) {
		ctx.SetStatusCode(statusCode)
		ctx.Response.Header.SetContentType("application/json")
		if err := json.NewEncoder(ctx).Encode(map[string]string{"error": message}); err != nil {
			log.Printf("Error sending JSON error: %v", err)
		}
	}

	requestHandler := func(ctx *fasthttp.RequestCtx) {
		path := string(ctx.Path())

		switch {
		case path == "/":
			fmt.Fprintf(ctx, "Hello World!")
		case strings.HasPrefix(path, "/lookup/"):
			ipStr := strings.TrimPrefix(path, "/lookup/")
			if ipStr == "" {
				writeJSONError(ctx, "Missing IP address in path", fasthttp.StatusBadRequest)
				return
			}

			ip := net.ParseIP(ipStr)
			if ip == nil {
				writeJSONError(ctx, "Invalid IP address", fasthttp.StatusBadRequest)
				return
			}

			cityRecord, err := cityDB.City(ip)
			if err != nil {
				writeJSONError(ctx, fmt.Sprintf("Error looking up IP in City DB: %v", err), fasthttp.StatusInternalServerError)
				return
			}

			response := map[string]interface{}{
				"ip":      ipStr,
				"code":    cityRecord.Country.IsoCode,
				"country": cityRecord.Country.Names["en"],
				"city":    cityRecord.City.Names["en"],
				"lat":     cityRecord.Location.Latitude,
				"lon":     cityRecord.Location.Longitude,
				"tz":      cityRecord.Location.TimeZone,
			}

			if asnDB != nil {
				asnRecord, err := asnDB.ASN(ip)
				if err == nil {
					response["asn"] = asnRecord.AutonomousSystemNumber
					response["isp"] = asnRecord.AutonomousSystemOrganization
				}
			}

			ctx.Response.Header.SetContentType("application/json")
			if err := json.NewEncoder(ctx).Encode(response); err != nil {
				writeJSONError(ctx, "Error encoding response", fasthttp.StatusInternalServerError)
			}
		default:
			writeJSONError(ctx, "Not Found", fasthttp.StatusNotFound)
		}
	}

	fmt.Println("Starting server on http://localhost:8080")
	if err := fasthttp.ListenAndServe(":8080", requestHandler); err != nil {
		log.Fatalf("Error in ListenAndServe: %s", err)
	}
}
