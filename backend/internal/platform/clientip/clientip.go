// Package clientip resuelve la IP real del cliente únicamente desde cabeceras
// fijadas por proxies o balanceadores de confianza (BR-169).
package clientip

import (
	"fmt"
	"net/http"
	"net/netip"
	"strings"
)

// DefaultHeader es la cabecera estándar consultada cuando el peer inmediato es
// un proxy de confianza y no se especificó otra cabecera.
const DefaultHeader = "X-Forwarded-For"

// Resolver determina la IP del cliente verificando primero si la conexión TCP
// inmediata proviene de un CIDR o IP de proxy confiable.
type Resolver struct {
	trustedProxies []netip.Prefix
	header         string
}

// NewResolver construye un Resolver a partir de una lista de CIDRs o IPs
// individuales confiables y el nombre de la cabecera de reenvío.
// Si trustedCIDRs está vacío, siempre devuelve la IP directa de r.RemoteAddr.
func NewResolver(trustedCIDRs []string, header string) (*Resolver, error) {
	prefixes := make([]netip.Prefix, 0, len(trustedCIDRs))
	for _, raw := range trustedCIDRs {
		entry := strings.TrimSpace(raw)
		if entry == "" {
			continue
		}
		prefix, err := parseTrustedPrefix(entry)
		if err != nil {
			return nil, fmt.Errorf("proxy confiable inválido %q: %w", entry, err)
		}
		prefixes = append(prefixes, prefix)
	}

	h := strings.TrimSpace(header)
	if h == "" {
		h = DefaultHeader
	} else {
		h = http.CanonicalHeaderKey(h)
	}

	return &Resolver{
		trustedProxies: prefixes,
		header:         h,
	}, nil
}

// Resolve retorna la IP del cliente para la petición HTTP dada.
// Nunca confía en cabeceras como X-Forwarded-For o X-Real-IP salvo que el peer
// inmediato (r.RemoteAddr) pertenezca a un proxy confiable configurado.
func (res *Resolver) Resolve(r *http.Request) string {
	if r == nil {
		return ""
	}
	peerIP, ok := parseAddr(r.RemoteAddr)
	if !ok {
		return ""
	}
	if res == nil || len(res.trustedProxies) == 0 || !res.isTrusted(peerIP) {
		return peerIP.String()
	}

	header := res.header
	if header == "" {
		header = DefaultHeader
	}
	values := r.Header.Values(header)
	if len(values) == 0 {
		return peerIP.String()
	}

	if strings.EqualFold(header, DefaultHeader) {
		return res.resolveForwardedFor(values, peerIP)
	}
	raw := strings.TrimSpace(values[len(values)-1])
	if addr, ok := parseAddr(raw); ok {
		return addr.String()
	}
	return peerIP.String()
}

// resolveForwardedFor recorre X-Forwarded-For de derecha a izquierda para
// descartar proxies confiables de la cadena y evitar falsificación en la izquierda.
func (res *Resolver) resolveForwardedFor(values []string, peerIP netip.Addr) string {
	var tokens []string
	for _, v := range values {
		tokens = append(tokens, strings.Split(v, ",")...)
	}

	var lastTrusted netip.Addr
	for i := len(tokens) - 1; i >= 0; i-- {
		addr, ok := parseAddr(tokens[i])
		if !ok {
			return peerIP.String()
		}
		if !res.isTrusted(addr) {
			return addr.String()
		}
		lastTrusted = addr
	}
	if lastTrusted.IsValid() {
		return lastTrusted.String()
	}
	return peerIP.String()
}

func (res *Resolver) isTrusted(addr netip.Addr) bool {
	for _, p := range res.trustedProxies {
		if p.Contains(addr) {
			return true
		}
	}
	return false
}

func parseTrustedPrefix(entry string) (netip.Prefix, error) {
	if addr, err := netip.ParseAddr(entry); err == nil {
		addr = addr.Unmap()
		return netip.PrefixFrom(addr, addr.BitLen()), nil
	}
	prefix, err := netip.ParsePrefix(entry)
	if err != nil {
		return netip.Prefix{}, err
	}
	return prefix.Masked(), nil
}

func parseAddr(raw string) (netip.Addr, bool) {
	s := strings.TrimSpace(raw)
	if s == "" {
		return netip.Addr{}, false
	}
	if ap, err := netip.ParseAddrPort(s); err == nil {
		return ap.Addr().Unmap(), true
	}
	if strings.HasPrefix(s, "[") && strings.HasSuffix(s, "]") {
		s = s[1 : len(s)-1]
	}
	addr, err := netip.ParseAddr(s)
	if err != nil {
		return netip.Addr{}, false
	}
	return addr.Unmap(), true
}
