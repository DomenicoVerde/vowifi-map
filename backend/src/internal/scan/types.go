package scan

// Epdg is a single EPDG IP address resolved for an operator, together with
// its geolocation and the outcome of the IKE_SA_INIT probe. Ip is nil when
// the operator's EPDG name didn't resolve to any address at all. Latitude
// and Longitude are nil when unknown (e.g. an IP with no GeoIP match).
// Response is "yes"/"no" depending on whether the EPDG answered the probe,
// or nil if it was never probed (e.g. a sinkholed loopback address).
type Epdg struct {
	Ip        *string  `json:"ip"`
	Latitude  *float64 `json:"latitude"`
	Longitude *float64 `json:"longitude"`
	Response  *string  `json:"response"`
}

// Mno mirrors the struct used by the original vulnerability scanner
// (backend/src/main.go), with EpdgsIp/Responses/Coordinates merged into a
// single Epdgs list.
type Mno struct {
	Mcc     string `json:"mcc"`
	Mnc     string `json:"mnc"`
	Iso     string `json:"iso"`
	Country string `json:"country"`
	Code    string `json:"code"`
	Network string `json:"network"`
	Epdgs   []Epdg `json:"epdgs"`
}
