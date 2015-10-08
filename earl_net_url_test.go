package earl_test

import (
	//"fmt"
	. "github.com/mikesimons/earl"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

// https://github.com/golang/go/blob/master/src/net/url/url_test.go

var _ = Describe("Earl net/url tests", func() {
	It("should parse with no path", func() {
		url := Parse("http://www.google.com")
		Expect(url.Scheme).Should(Equal("http"))
		Expect(url.Host).Should(Equal("www.google.com"))
	})

	It("should parse with path", func() {
		url := Parse("http://www.google.com/")
		Expect(url.Scheme).Should(Equal("http"))
		Expect(url.Host).Should(Equal("www.google.com"))
		Expect(url.Path).Should(Equal("/"))
	})

	PIt("should parse path with hex escaping", func() {
		url := Parse("http://www.google.com/file%20one%26two")
		Expect(url.Scheme).Should(Equal("http"))
		Expect(url.Host).Should(Equal("www.google.com"))
		Expect(url.Path).Should(Equal("/file%20one%26two"))

		nUrl := url.Normalize()
		Expect(nUrl.Path).Should(Equal("/file one&two"))
	})

	It("should parse user", func() {
		url := Parse("ftp://webmaster@www.google.com/")
		Expect(url.Scheme).Should(Equal("ftp"))
		Expect(url.UserInfo).Should(Equal("webmaster"))
		Expect(url.Host).Should(Equal("www.google.com"))
		Expect(url.Path).Should(Equal("/"))
	})

	PIt("should parse user with pct-encoding in username", func() {
		url := Parse("ftp://john%20doe@www.google.com/")
		Expect(url.Scheme).Should(Equal("ftp"))
		Expect(url.UserInfo).Should(Equal("john%20doe"))
		Expect(url.Host).Should(Equal("www.google.com"))
		Expect(url.Path).Should(Equal("/"))

		nUrl := url.Normalize()
		Expect(nUrl.UserInfo).Should(Equal("john doe"))
	})

	It("should parse query", func() {
		url := Parse("http://www.google.com/?q=go+language")
		Expect(url.Path).Should(Equal("/"))
		Expect(url.Query).Should(Equal("q=go+language"))
	})

	It("should not decode query with pct-encoding", func() {
		url := Parse("http://www.google.com/?q=go%20language")
		Expect(url.Path).Should(Equal("/"))
		Expect(url.Query).Should(Equal("q=go%20language"))
	})

	PIt("should decode path with pct-encoding", func() {
		url := Parse("http://www.google.com/a%20b?q=c+d")
		Expect(url.Path).Should(Equal("/a%20b"))
		Expect(url.Query).Should(Equal("q=c+d"))

		nUrl := url.Normalize()
		Expect(nUrl.Path).Should(Equal("/a b"))
	})

	It("should correctly parse paths without leading slash", func() {
		url := Parse("http:www.google.com/?q=go+language")
		Expect(url.Scheme).Should(Equal("http"))
		Expect(url.Opaque).Should(Equal("www.google.com/"))
		Expect(url.Query).Should(Equal("q=go+language"))
	})

	It("should correctly parse mailto with path", func() {
		url := Parse("mailto:/webmaster@golang.org")
		Expect(url.Scheme).Should(Equal("mailto"))
		Expect(url.Path).Should(Equal("/webmaster@golang.org"))
	})

	It("should correctly parse mailto", func() {
		url := Parse("mailto:webmaster@golang.org")
		Expect(url.Scheme).Should(Equal("mailto"))
		Expect(url.Opaque).Should(Equal("webmaster@golang.org"))
	})

	It("should not produce invalid scheme if there is an unescaped :// in query", func() {
		url := Parse("/foo?query=http://bad")
		Expect(url.Scheme).Should(Equal(""))
		Expect(url.Path).Should(Equal("/foo"))
		Expect(url.Query).Should(Equal("query=http://bad"))
	})

	It("should handle urls starting //", func() {
		url := Parse("//foo")
		Expect(url.Host).Should(Equal("foo"))
	})

	It("should handle urls starting // with userinfo, path & query", func() {
		url := Parse("//user@foo/path?a=b")
		Expect(url.Host).Should(Equal("foo"))
		Expect(url.UserInfo).Should(Equal("user"))
		Expect(url.Query).Should(Equal("a=b"))
		Expect(url.Path).Should(Equal("/path"))
	})

	It("should handle urls starting ///", func() {
		url := Parse("///threeslashes")
		Expect(url.Path).Should(Equal("///threeslashes"))
	})

	It("should handle user / pass", func() {
		url := Parse("http://user:password@google.com")
		Expect(url.Scheme).Should(Equal("http"))
		Expect(url.UserInfo).Should(Equal("user:password"))
		Expect(url.Host).Should(Equal("google.com"))
	})

	It("should handle unescaped @ in username", func() {
		url := Parse("http://j@ne:password@google.com")
		Expect(url.UserInfo).Should(Equal("j@ne:password"))
		Expect(url.Host).Should(Equal("google.com"))
	})

	It("should handle unescaped @ in password", func() {
		url := Parse("http://jane:p@ssword@google.com")
		Expect(url.UserInfo).Should(Equal("jane:p@ssword"))
		Expect(url.Host).Should(Equal("google.com"))
	})

	It("should handle @ all over the place", func() {
		url := Parse("http://j@ne:p@ssword@google.com/p@th?q=@go")
		Expect(url.Scheme).Should(Equal("http"))
		Expect(url.UserInfo).Should(Equal("j@ne:p@ssword"))
		Expect(url.Host).Should(Equal("google.com"))
		Expect(url.Path).Should(Equal("/p@th"))
		Expect(url.Query).Should(Equal("q=@go"))
	})

	It("should handle fragment", func() {
		url := Parse("http://www.google.com/?q=go+language#foo")
		Expect(url.Query).Should(Equal("q=go+language"))
		Expect(url.Fragment).Should(Equal("foo"))
	})
/*
{
	"http://www.google.com/?q=go+language#foo%26bar",
	&URL{
		Scheme:   "http",
		Host:     "www.google.com",
		Path:     "/",
		RawQuery: "q=go+language",
		Fragment: "foo&bar",
	},
	"http://www.google.com/?q=go+language#foo&bar",
},
{
	"file:///home/adg/rabbits",
	&URL{
		Scheme: "file",
		Host:   "",
		Path:   "/home/adg/rabbits",
	},
	"file:///home/adg/rabbits",
},
// "Windows" paths are no exception to the rule.
// See golang.org/issue/6027, especially comment #9.
{
	"file:///C:/FooBar/Baz.txt",
	&URL{
		Scheme: "file",
		Host:   "",
		Path:   "/C:/FooBar/Baz.txt",
	},
	"file:///C:/FooBar/Baz.txt",
},
// case-insensitive scheme
{
	"MaIlTo:webmaster@golang.org",
	&URL{
		Scheme: "mailto",
		Opaque: "webmaster@golang.org",
	},
	"mailto:webmaster@golang.org",
},
// Relative path
{
	"a/b/c",
	&URL{
		Path: "a/b/c",
	},
	"a/b/c",
},
// escaped '?' in username and password
{
	"http://%3Fam:pa%3Fsword@google.com",
	&URL{
		Scheme: "http",
		User:   UserPassword("?am", "pa?sword"),
		Host:   "google.com",
	},
	"",
},
// host subcomponent; IPv4 address in RFC 3986
{
	"http://192.168.0.1/",
	&URL{
		Scheme: "http",
		Host:   "192.168.0.1",
		Path:   "/",
	},
	"",
},
// host and port subcomponents; IPv4 address in RFC 3986
{
	"http://192.168.0.1:8080/",
	&URL{
		Scheme: "http",
		Host:   "192.168.0.1:8080",
		Path:   "/",
	},
	"",
},
// host subcomponent; IPv6 address in RFC 3986
{
	"http://[fe80::1]/",
	&URL{
		Scheme: "http",
		Host:   "[fe80::1]",
		Path:   "/",
	},
	"",
},
// host and port subcomponents; IPv6 address in RFC 3986
{
	"http://[fe80::1]:8080/",
	&URL{
		Scheme: "http",
		Host:   "[fe80::1]:8080",
		Path:   "/",
	},
	"",
},
// host subcomponent; IPv6 address with zone identifier in RFC 6847
{
	"http://[fe80::1%25en0]/", // alphanum zone identifier
	&URL{
		Scheme: "http",
		Host:   "[fe80::1%en0]",
		Path:   "/",
	},
	"",
},
// host and port subcomponents; IPv6 address with zone identifier in RFC 6847
{
	"http://[fe80::1%25en0]:8080/", // alphanum zone identifier
	&URL{
		Scheme: "http",
		Host:   "[fe80::1%en0]:8080",
		Path:   "/",
	},
	"",
},
// host subcomponent; IPv6 address with zone identifier in RFC 6847
{
	"http://[fe80::1%25%65%6e%301-._~]/", // percent-encoded+unreserved zone identifier
	&URL{
		Scheme: "http",
		Host:   "[fe80::1%en01-._~]",
		Path:   "/",
	},
	"http://[fe80::1%25en01-._~]/",
},
// host and port subcomponents; IPv6 address with zone identifier in RFC 6847
{
	"http://[fe80::1%25%65%6e%301-._~]:8080/", // percent-encoded+unreserved zone identifier
	&URL{
		Scheme: "http",
		Host:   "[fe80::1%en01-._~]:8080",
		Path:   "/",
	},
	"http://[fe80::1%25en01-._~]:8080/",
},
// alternate escapings of path survive round trip
{
	"http://rest.rsc.io/foo%2fbar/baz%2Fquux?alt=media",
	&URL{
		Scheme:   "http",
		Host:     "rest.rsc.io",
		Path:     "/foo/bar/baz/quux",
		RawPath:  "/foo%2fbar/baz%2Fquux",
		RawQuery: "alt=media",
	},
	"",
},
// issue 12036
{
	"mysql://a,b,c/bar",
	&URL{
		Scheme: "mysql",
		Host:   "a,b,c",
		Path:   "/bar",
	},
	"",
},
// worst case host, still round trips
{
	"scheme://!$&'()*+,;=hello!:port/path",
	&URL{
		Scheme: "scheme",
		Host:   "!$&'()*+,;=hello!:port",
		Path:   "/path",
	},
	"",
},
// worst case path, still round trips
{
	"http://host/!$&'()*+,;=:@[hello]",
	&URL{
		Scheme:  "http",
		Host:    "host",
		Path:    "/!$&'()*+,;=:@[hello]",
		RawPath: "/!$&'()*+,;=:@[hello]",
	},
	"",
},
// golang.org/issue/5684
{
	"http://example.com/oid/[order_id]",
	&URL{
		Scheme:  "http",
		Host:    "example.com",
		Path:    "/oid/[order_id]",
		RawPath: "/oid/[order_id]",
	},
	"",
},
*/
})
