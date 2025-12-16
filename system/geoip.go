package system

// Uses GeoLite2 Country database by MaxMind
// https://www.maxmind.com

import (
	_ "embed"
	"net"

	"github.com/oschwald/maxminddb-golang"
)

type GeoIP struct {
	db *maxminddb.Reader
}

var geoip *GeoIP

func init() {
	var err error
	geoip, err = OpenGeoIP()
	if err != nil {
		geoip = nil
	}
}

func GetCountryByIP(ipStr string) (string, error) {
	if geoip == nil {
		return "", nil
	}
	var result struct {
		Country struct {
			IsoCode string            `maxminddb:"iso_code"`
			Names   map[string]string `maxminddb:"names"`
		} `maxminddb:"country"`
	}
	ip := net.ParseIP(ipStr)
	err := geoip.db.Lookup(ip, &result)
	if err != nil {
		return "", err
	}

	name := result.Country.Names["en"]
	if name == "" {
		name = result.Country.IsoCode
	}

	return name, nil
}

//go:embed geoip.mmdb
var geoip_mmdb []byte

func OpenGeoIP() (*GeoIP, error) {
	db, err := maxminddb.FromBytes(geoip_mmdb)
	if err != nil {
		return nil, err
	}
	return &GeoIP{db: db}, nil
}

func (g *GeoIP) Close() error {
	return g.db.Close()
}
