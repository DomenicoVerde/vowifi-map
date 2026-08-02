package persistence

import (
	"fmt"
	"net"

	"github.com/oschwald/geoip2-golang"
)

// Coordinate is the geolocation of a single EPDG IP address, as resolved
// from the GeoLite2-City database.
type Coordinate struct {
	Ip        string  `json:"ip"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	City      string  `json:"city,omitempty"`
	Country   string  `json:"country,omitempty"`
}

// GeoIP wraps a GeoLite2-City database reader.
type GeoIP struct {
	reader *geoip2.Reader
}

// OpenGeoIP opens the GeoLite2-City.mmdb database at the given path.
func OpenGeoIP(path string) (*GeoIP, error) {
	reader, err := geoip2.Open(path)
	if err != nil {
		return nil, fmt.Errorf("OpenGeoIP(): failed to open %q: %w", path, err)
	}
	return &GeoIP{reader: reader}, nil
}

// Close releases the underlying database file.
func (g *GeoIP) Close() error {
	return g.reader.Close()
}

// Lookup returns the geolocation of the given IP address, or nil if the
// address has no location in the database (e.g. private/reserved ranges).
func (g *GeoIP) Lookup(ip net.IP) (*Coordinate, error) {
	record, err := g.reader.City(ip)
	if err != nil {
		return nil, fmt.Errorf("Lookup(): city lookup failed for %s: %w", ip, err)
	}

	if record.Location.Latitude == 0 && record.Location.Longitude == 0 {
		return nil, nil
	}

	coordinate := &Coordinate{
		Ip:        ip.String(),
		Latitude:  record.Location.Latitude,
		Longitude: record.Location.Longitude,
		City:      record.City.Names["en"],
		Country:   record.Country.Names["en"],
	}

	return coordinate, nil
}
