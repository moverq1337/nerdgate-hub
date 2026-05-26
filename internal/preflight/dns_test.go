package preflight

import (
	"net/netip"
	"testing"
)

func TestContainsIPMatchesMappedAddress(t *testing.T) {
	expected := netip.MustParseAddr("203.0.113.10")
	addrs := []netip.Addr{
		netip.MustParseAddr("::ffff:203.0.113.10"),
	}

	if !ContainsIP(addrs, expected) {
		t.Fatal("expected mapped IPv4 address to match")
	}
}

func TestFormatIPs(t *testing.T) {
	got := FormatIPs([]netip.Addr{
		netip.MustParseAddr("203.0.113.10"),
		netip.MustParseAddr("2001:db8::1"),
	})

	want := "203.0.113.10, 2001:db8::1"
	if got != want {
		t.Fatalf("expected %q, got %q", want, got)
	}
}
