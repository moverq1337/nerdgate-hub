package preflight

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/netip"
	"sort"
	"strings"
)

type DNSResult struct {
	Domain   string
	Expected netip.Addr
	Resolved []netip.Addr
	Match    bool
}

func CheckDNS(ctx context.Context, resolver *net.Resolver, domain string, expected netip.Addr) (DNSResult, error) {
	domain = strings.TrimSpace(strings.TrimSuffix(domain, "."))
	if domain == "" {
		return DNSResult{}, errors.New("domain is required")
	}
	if !expected.IsValid() {
		return DNSResult{}, errors.New("expected IP is invalid")
	}

	addrs, err := resolver.LookupIPAddr(ctx, domain)
	if err != nil {
		return DNSResult{}, err
	}

	resolved := make([]netip.Addr, 0, len(addrs))
	for _, addr := range addrs {
		parsed, ok := netip.AddrFromSlice(addr.IP)
		if !ok {
			continue
		}
		resolved = append(resolved, parsed.Unmap())
	}

	sort.Slice(resolved, func(i, j int) bool {
		return resolved[i].String() < resolved[j].String()
	})

	return DNSResult{
		Domain:   domain,
		Expected: expected.Unmap(),
		Resolved: resolved,
		Match:    ContainsIP(resolved, expected),
	}, nil
}

func ContainsIP(addrs []netip.Addr, expected netip.Addr) bool {
	expected = expected.Unmap()
	for _, addr := range addrs {
		if addr.Unmap() == expected {
			return true
		}
	}
	return false
}

func FormatIPs(addrs []netip.Addr) string {
	if len(addrs) == 0 {
		return "(none)"
	}

	values := make([]string, 0, len(addrs))
	for _, addr := range addrs {
		values = append(values, addr.String())
	}
	return strings.Join(values, ", ")
}

func FormatMismatch(domain string, expected netip.Addr, resolved []netip.Addr) string {
	return fmt.Sprintf("%s must resolve to %s, currently resolves to %s", domain, expected, FormatIPs(resolved))
}
