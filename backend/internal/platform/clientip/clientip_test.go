package clientip

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestResolver_NoTrustedProxiesIgnoresSpoofedHeaders(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		resolver *Resolver
		remote   string
		headers  map[string]string
		want     string
	}{
		{
			name:     "nil resolver strips port and ignores X-Forwarded-For",
			resolver: nil,
			remote:   "198.51.100.77:54321",
			headers: map[string]string{
				"X-Forwarded-For": "1.2.3.4",
				"X-Real-IP":       "5.6.7.8",
			},
			want: "198.51.100.77",
		},
		{
			name:     "zero-value resolver ignores headers",
			resolver: &Resolver{},
			remote:   "203.0.113.10:1234",
			headers: map[string]string{
				"X-Forwarded-For": "10.0.0.1",
			},
			want: "203.0.113.10",
		},
		{
			name:     "empty trusted list ignores X-Forwarded-For",
			resolver: mustResolver(t, nil, ""),
			remote:   "203.0.113.11:9999",
			headers: map[string]string{
				"X-Forwarded-For": "198.51.100.1",
			},
			want: "203.0.113.11",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			req.RemoteAddr = tc.remote
			for k, v := range tc.headers {
				req.Header.Set(k, v)
			}
			if got := tc.resolver.Resolve(req); got != tc.want {
				t.Fatalf("Resolve() = %q, want %q", got, tc.want)
			}
		})
	}
}

func TestResolver_TrustedProxyCIDR(t *testing.T) {
	t.Parallel()

	t.Run("reads single IP from X-Forwarded-For", func(t *testing.T) {
		t.Parallel()
		res := mustResolver(t, []string{"10.0.0.0/8"}, "")
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.RemoteAddr = "10.1.2.3:4567"
		req.Header.Set("X-Forwarded-For", "198.51.100.42")

		if got := res.Resolve(req); got != "198.51.100.42" {
			t.Fatalf("Resolve() = %q, want 198.51.100.42", got)
		}
	})

	t.Run("ignores left-most spoofed value across multiple proxies", func(t *testing.T) {
		t.Parallel()
		res := mustResolver(t, []string{"10.0.0.0/8", "172.16.0.0/12"}, "X-Forwarded-For")
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.RemoteAddr = "10.0.0.5:8080"
		// Attacker sent "9.9.9.9"; edge proxy 172.16.0.10 saw 198.51.100.99.
		req.Header.Set("X-Forwarded-For", "9.9.9.9, 198.51.100.99, 172.16.0.10")

		if got := res.Resolve(req); got != "198.51.100.99" {
			t.Fatalf("Resolve() = %q, want 198.51.100.99", got)
		}
	})

	t.Run("supports multiple X-Forwarded-For header lines", func(t *testing.T) {
		t.Parallel()
		res := mustResolver(t, []string{"127.0.0.1"}, "")
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.RemoteAddr = "127.0.0.1:5000"
		req.Header.Add("X-Forwarded-For", "9.9.9.9")
		req.Header.Add("X-Forwarded-For", "203.0.113.55")

		if got := res.Resolve(req); got != "203.0.113.55" {
			t.Fatalf("Resolve() = %q, want 203.0.113.55", got)
		}
	})

	t.Run("supports custom single-IP header like X-Real-IP", func(t *testing.T) {
		t.Parallel()
		res := mustResolver(t, []string{"10.20.0.0/16"}, "x-real-ip")
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.RemoteAddr = "10.20.1.2:443"
		req.Header.Set("X-Real-IP", "198.51.100.88")
		req.Header.Set("X-Forwarded-For", "1.1.1.1")

		if got := res.Resolve(req); got != "198.51.100.88" {
			t.Fatalf("Resolve() = %q, want 198.51.100.88", got)
		}
	})

	t.Run("falls back to peer IP when trusted proxy sends no header", func(t *testing.T) {
		t.Parallel()
		res := mustResolver(t, []string{"10.0.0.0/8"}, "")
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.RemoteAddr = "10.0.0.8:1234"

		if got := res.Resolve(req); got != "10.0.0.8" {
			t.Fatalf("Resolve() = %q, want 10.0.0.8", got)
		}
	})
}

func TestResolver_UntrustedProxyFallsBackToRemoteAddr(t *testing.T) {
	t.Parallel()

	res := mustResolver(t, []string{"10.0.0.0/8"}, "X-Forwarded-For")
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.RemoteAddr = "198.51.100.25:33333"
	req.Header.Set("X-Forwarded-For", "127.0.0.1")
	req.Header.Set("X-Real-IP", "127.0.0.1")

	if got := res.Resolve(req); got != "198.51.100.25" {
		t.Fatalf("Resolve() = %q, want 198.51.100.25", got)
	}
}

func TestResolver_MalformedHeadersAndIPv6(t *testing.T) {
	t.Parallel()

	t.Run("malformed X-Forwarded-For falls back to peer IP", func(t *testing.T) {
		t.Parallel()
		res := mustResolver(t, []string{"10.0.0.0/8"}, "")
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.RemoteAddr = "10.0.0.3:7000"
		req.Header.Set("X-Forwarded-For", "198.51.100.1, not-an-ip")

		if got := res.Resolve(req); got != "10.0.0.3" {
			t.Fatalf("Resolve() = %q, want 10.0.0.3", got)
		}
	})

	t.Run("malformed X-Real-IP falls back to peer IP", func(t *testing.T) {
		t.Parallel()
		res := mustResolver(t, []string{"10.0.0.0/8"}, "X-Real-IP")
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.RemoteAddr = "10.0.0.3:7000"
		req.Header.Set("X-Real-IP", "1.2.3.4, 5.6.7.8")

		if got := res.Resolve(req); got != "10.0.0.3" {
			t.Fatalf("Resolve() = %q, want 10.0.0.3", got)
		}
	})

	t.Run("IPv6 trusted proxy and IPv6 client", func(t *testing.T) {
		t.Parallel()
		res := mustResolver(t, []string{"fd00::/8"}, "")
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.RemoteAddr = "[fd00::10]:44321"
		req.Header.Set("X-Forwarded-For", "[2001:db8::42]:9000")

		if got := res.Resolve(req); got != "2001:db8::42" {
			t.Fatalf("Resolve() = %q, want 2001:db8::42", got)
		}
	})

	t.Run("IPv4-mapped IPv6 peer matches IPv4 CIDR", func(t *testing.T) {
		t.Parallel()
		res := mustResolver(t, []string{"127.0.0.0/8"}, "")
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.RemoteAddr = "[::ffff:127.0.0.1]:54321"
		req.Header.Set("X-Forwarded-For", "203.0.113.7")

		if got := res.Resolve(req); got != "203.0.113.7" {
			t.Fatalf("Resolve() = %q, want 203.0.113.7", got)
		}
	})

	t.Run("malformed RemoteAddr returns empty string", func(t *testing.T) {
		t.Parallel()
		res := mustResolver(t, []string{"10.0.0.0/8"}, "")
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.RemoteAddr = "not-a-valid-addr"

		if got := res.Resolve(req); got != "" {
			t.Fatalf("Resolve() = %q, want empty", got)
		}
	})

	t.Run("invalid CIDR returns error on construction", func(t *testing.T) {
		t.Parallel()
		if _, err := NewResolver([]string{"10.0.0.0/99"}, ""); err == nil {
			t.Fatal("expected error for invalid CIDR, got nil")
		}
	})
}

func mustResolver(t *testing.T, cidrs []string, header string) *Resolver {
	t.Helper()
	res, err := NewResolver(cidrs, header)
	if err != nil {
		t.Fatalf("NewResolver(%v, %q) unexpected error: %v", cidrs, header, err)
	}
	return res
}
