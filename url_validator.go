package utils

import (
	"fmt"
	"net"
	"net/url"
	"strings"
)

var allowedSchemes = map[string]bool{
	"http":  true,
	"https": true,
}

func ValidateURL(rawURL string) error {
	if rawURL == "" {
		return nil
	}

	parsed, err := url.Parse(rawURL)
	if err != nil {
		return fmt.Errorf("invalid URL: %w", err)
	}

	if !allowedSchemes[parsed.Scheme] {
		return fmt.Errorf("unsupported scheme: %s", parsed.Scheme)
	}

	host := parsed.Hostname()
	if host == "" {
		return fmt.Errorf("empty host in URL")
	}

	if isPrivateOrLocalhost(host) {
		return fmt.Errorf("access to private/localhost addresses is not allowed: %s", host)
	}

	return nil
}

func isPrivateOrLocalhost(host string) bool {
	if host == "localhost" || strings.HasSuffix(host, ".localhost") {
		return true
	}

	ip := net.ParseIP(host)
	if ip == nil {
		return false
	}

	if ip.IsLoopback() || ip.IsPrivate() || ip.IsLinkLocalUnicast() || ip.IsUnspecified() {
		return true
	}

	return false
}
