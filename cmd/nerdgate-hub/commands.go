package main

import (
	"context"
	"fmt"
	"net"
	"net/netip"
	"os"
	"time"

	"github.com/nerdgatehub/nerdgate-hub/internal/preflight"
	"github.com/nerdgatehub/nerdgate-hub/internal/store"
)

func runCommand(name string, args []string) {
	switch name {
	case "check-domain":
		runCheckDomain(args)
	case "reset-password":
		runResetPassword(args)
	default:
		fmt.Fprintf(os.Stderr, "unknown command %q\n", name)
		os.Exit(2)
	}
}

func runCheckDomain(args []string) {
	if len(args) < 1 || len(args) > 2 {
		fmt.Fprintln(os.Stderr, "usage: nerdgate-hub check-domain DOMAIN [EXPECTED_IP]")
		os.Exit(2)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 8*time.Second)
	defer cancel()

	expected, err := expectedIP(ctx, args)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to determine expected IP: %v\n", err)
		os.Exit(1)
	}

	result, err := preflight.CheckDNS(ctx, net.DefaultResolver, args[0], expected)
	if err != nil {
		fmt.Fprintf(os.Stderr, "DNS check failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("domain: %s\n", result.Domain)
	fmt.Printf("expected: %s\n", result.Expected)
	fmt.Printf("resolved: %s\n", preflight.FormatIPs(result.Resolved))

	if !result.Match {
		fmt.Fprintln(os.Stderr, "result: mismatch")
		os.Exit(1)
	}

	fmt.Println("result: ok")
}

func expectedIP(ctx context.Context, args []string) (netip.Addr, error) {
	if len(args) == 2 {
		return netip.ParseAddr(args[1])
	}
	return preflight.DetectPublicIP(ctx)
}

func runResetPassword(args []string) {
	if len(args) != 3 {
		fmt.Fprintln(os.Stderr, "usage: nerdgate-hub reset-password DATA_DIR USERNAME PASSWORD")
		os.Exit(2)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	db, err := store.Open(args[0])
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open store: %v\n", err)
		os.Exit(1)
	}
	defer db.Close()

	if err := db.ResetPassword(ctx, args[1], args[2]); err != nil {
		fmt.Fprintf(os.Stderr, "failed to reset password: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("password reset for %s\n", args[1])
}
