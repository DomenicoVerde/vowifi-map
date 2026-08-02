package scan

import (
	"net"
	"time"

	"vowifi_scanner/internal/ike"
)

// ProbeEPDG sends a single IKE_SA_INIT request to the given EPDG IP address
// on port 500/UDP and reports whether any IKE response was received within
// the given timeout. It does not attempt to complete the IKE handshake or
// test for vulnerabilities: any well-formed IKE reply (SA, Notify, Cookie,
// ...) is treated as proof that the EPDG is reachable and alive.
func ProbeEPDG(ip net.IP, timeout time.Duration) bool {
	network := "udp"
	host := ip.String()
	if ip.To4() == nil {
		network = "udp6"
	}

	remoteAddr, err := net.ResolveUDPAddr(network, net.JoinHostPort(host, "500"))
	if err != nil {
		return false
	}

	connection, err := net.DialUDP(network, nil, remoteAddr)
	if err != nil {
		return false
	}
	defer connection.Close()

	requestData, err := ike.BuildIKESAInitRequest()
	if err != nil {
		return false
	}

	if _, err = connection.Write(requestData); err != nil {
		return false
	}

	if err = connection.SetReadDeadline(time.Now().Add(timeout)); err != nil {
		return false
	}

	buffer := make([]byte, 65535)
	n, err := connection.Read(buffer)
	if err != nil || n < 28 {
		return false
	}

	response := new(ike.IKEMessage)
	if err = response.Decode(buffer[:n]); err != nil {
		return false
	}

	return true
}
