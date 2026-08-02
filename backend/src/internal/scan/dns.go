package scan

import (
	"context"
	"net"
	"time"
)

// dnsRetryAttempts is the number of times a DNS lookup is retried before
// an operator is marked as red. Under heavy concurrency, local resolvers
// (e.g. systemd-resolved) can return a transient SERVFAIL/timeout for a
// name that actually resolves fine; retrying avoids misclassifying those
// operators as unreachable.
const dnsRetryAttempts = 3

const dnsRetryBackoff = 300 * time.Millisecond

// dnsResolver queries Google's public DNS (8.8.8.8) directly instead of the
// host's configured resolver. The local stub resolver (e.g.
// systemd-resolved) buckles under the concurrent lookup load of a full
// scan and returns transient SERVFAIL/timeout for names that resolve fine
// (observed misclassifying known-good operators, e.g. TIM Italy, as red).
var dnsResolver = &net.Resolver{
	PreferGo: true,
	Dial: func(ctx context.Context, network, _ string) (net.Conn, error) {
		return (&net.Dialer{}).DialContext(ctx, network, "8.8.8.8:53")
	},
}

// lookupEPDGAddresses resolves fqdn, retrying up to dnsRetryAttempts times
// on failure. It returns an empty slice if every attempt fails.
func lookupEPDGAddresses(fqdn string, timeout time.Duration) []net.IP {
	for attempt := 1; attempt <= dnsRetryAttempts; attempt++ {
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		ips, err := dnsResolver.LookupIP(ctx, "ip", fqdn)
		cancel()

		if err == nil && len(ips) > 0 {
			return ips
		}

		if attempt < dnsRetryAttempts {
			time.Sleep(dnsRetryBackoff)
		}
	}

	return nil
}
