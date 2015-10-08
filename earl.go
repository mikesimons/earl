package earl

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"
)

// URL represents state for a parsed URL
type URL struct {
	Input string

	Scheme   string
	Opaque   string
	Query    string
	Fragment string

	// Elements of Opaque
	Authority string
	Path      string

	// Elements of Authority
	UserInfo string
	Host     string
	Port     string
}

// RFC3986: https://www.ietf.org/rfc/rfc3986.txt
// URI scheme registry: http://www.iana.org/assignments/uri-schemes/uri-schemes.xhtml
// TODO: Normalize method; See RFC3986 section 6.2.2 for normalization ref

func namedMatches(matches []string, r *regexp.Regexp) map[string]string {
	result := make(map[string]string)
	for i, name := range r.SubexpNames() {
		if name == "" {
			continue
		}
		if i >= len(matches) {
			result[name] = ""
		} else {
			result[name] = matches[i]
		}
	}
	return result
}

func splitAuthorityFromPath(opaque string) (string, string) {
	r := regexp.MustCompile("((//)?(?P<authority>[^/]+))?(?P<path>/.*)?")
	matches := namedMatches(r.FindStringSubmatch(opaque), r)
	return matches["authority"], matches["path"]
}

func splitUserinfoHostPortFromAuthority(authority string) (string, string, string) {
	userinfo := ""
	if delimPos := strings.LastIndex(authority, "@"); delimPos != -1 {
		userinfo = authority[0:delimPos]
		authority = authority[delimPos+1:]
	}

	parts := []string{
		"(", "(\\[(?P<host6>[^\\]]+)\\])", "|", "(?P<host>[^:]+)", ")?", // host6 | host
		"(:(?P<port>[0-9]+))?",
	}

	r := regexp.MustCompile(strings.Join(parts, ""))
	matches := namedMatches(r.FindStringSubmatch(authority), r)
	if matches["host"] == "" {
		matches["host"] = matches["host6"]
	}

	return userinfo, matches["host"], matches["port"]
}

// Split splits an URL in to its major components (scheme, opaque, query, fragment)
func Split(url string) (string, string, string, string) {
	parts := []string{
		"^((?P<scheme>[^:\\.]+):)?", // scheme is required by RFC3986 (S3) but we are intentionally allowing it to be omitted for convenience
		"(?P<opaque>(//)?[^?#]+)",   // hier-part
		"(\\?(?P<query>[^#]+))?",    // query
		"(#(?P<fragment>.*))?",      // fragment
	}

	r := regexp.MustCompile(strings.Join(parts, ""))
	matches := namedMatches(r.FindStringSubmatch(url), r)

	return matches["scheme"], matches["opaque"], matches["query"], matches["fragment"]
}

// Parse splits an URL in to as many parts as it can
func Parse(url string) *URL {
	result := &URL{}
	result.Input = url
	result.Scheme, result.Opaque, result.Query, result.Fragment = Split(url)
	result.Authority, result.Path = splitAuthorityFromPath(result.Opaque)
	result.UserInfo, result.Host, result.Port = splitUserinfoHostPortFromAuthority(result.Authority)
	return result
}

// ToNetURL converts an earl.URL in to a net/url.URL
func (u *URL) ToNetURL() *url.URL {
	// FIXME users of net/url may expect most of these to be decoded
	host := ""
	if u.Host != "" {
		host = u.Host

		if u.Port != "" {
			host = fmt.Sprintf("%s:%s", host, u.Port)
		}
	}

	ret := &url.URL{
		Scheme: u.Scheme,
		//User: TODO
		Host:     host,
		Path:     u.Path,
		RawPath:  u.Path,
		RawQuery: u.Query,
		Fragment: u.Fragment,
	}

	if u.Authority == "" {
		ret.Opaque = u.Opaque
	}

	return ret
}

// Normalize is intended to produce an expanded and valid URL representation
// It is presently incomplete (also read not even started).
func (u *URL) Normalize() *URL {
	result := &URL{}
	result.Scheme = strings.ToLower(u.Scheme)
	result.Host = strings.ToLower(u.Host)
	return result
}
