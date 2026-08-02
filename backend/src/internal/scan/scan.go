package scan

import (
	"fmt"
	"time"

	"vowifi_scanner/internal/persistence"
)

// ScanOptions controls how an operator's EPDG is probed.
type ScanOptions struct {
	DNSTimeout time.Duration
	IKETimeout time.Duration
}

const ikeResponseYes = "yes"
const ikeResponseNo = "no"

// ScanOperator resolves the EPDG FQDN of the given operator, probes every
// resolved IP address with a single IKE_SA_INIT request, and geolocates
// each IP via the provided GeoIP database. Loopback/sinkhole addresses
// (e.g. 127.0.0.1) are never probed or geolocated. If the operator's EPDG
// name doesn't resolve to any address at all, Epdgs is left empty.
func ScanOperator(operator persistence.Operator, geo *persistence.GeoIP, options ScanOptions) Mno {
	result := Mno{
		Mcc:     operator.Mcc,
		Mnc:     operator.Mnc,
		Iso:     operator.Iso,
		Country: operator.Country,
		Code:    operator.CountryCode,
		Network: operator.Network,
		Epdgs:   []Epdg{},
	}

	fqdn := fmt.Sprintf("epdg.epc.mnc%s.mcc%s.pub.3gppnetwork.org", operator.Mnc, operator.Mcc)

	ips := lookupEPDGAddresses(fqdn, options.DNSTimeout)
	if len(ips) == 0 {
		return result
	}

	for _, ip := range ips {
		ipStr := ip.String()
		entry := Epdg{Ip: &ipStr}

		if ip.IsLoopback() {
			result.Epdgs = append(result.Epdgs, entry)
			continue
		}

		if coordinate, geoErr := geo.Lookup(ip); geoErr == nil && coordinate != nil {
			latitude := coordinate.Latitude
			longitude := coordinate.Longitude
			entry.Latitude = &latitude
			entry.Longitude = &longitude
		}

		response := ikeResponseNo
		if ProbeEPDG(ip, options.IKETimeout) {
			response = ikeResponseYes
		}
		entry.Response = &response

		result.Epdgs = append(result.Epdgs, entry)
	}

	return result
}
