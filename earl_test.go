package earl_test

import (
	"fmt"
	. "github.com/mikesimons/earl"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Earl", func() {
	Describe("Split", func() {
		It("should split url in to separate components", func() {
			r1, r2, r3, r4 := Split("scheme://opaque?query#fragment")
			Expect(r1).Should(Equal("scheme"))
			Expect(r2).Should(Equal("//opaque"))
			Expect(r3).Should(Equal("query"))
			Expect(r4).Should(Equal("fragment"))
		})

		It("should allow omission of scheme component", func() {
			r1, r2, r3, r4 := Split("opaque?query#fragment")
			Expect(r1).Should(Equal(""))
			Expect(r2).Should(Equal("opaque"))
			Expect(r3).Should(Equal("query"))
			Expect(r4).Should(Equal("fragment"))
		})

		It("should allow omission of query component", func() {
			r1, r2, r3, r4 := Split("scheme://opaque#fragment")
			Expect(r1).Should(Equal("scheme"))
			Expect(r2).Should(Equal("//opaque"))
			Expect(r3).Should(Equal(""))
			Expect(r4).Should(Equal("fragment"))
		})

		It("should allow omission of fragment component", func() {
			r1, r2, r3, r4 := Split("scheme://opaque?query")
			Expect(r1).Should(Equal("scheme"))
			Expect(r2).Should(Equal("//opaque"))
			Expect(r3).Should(Equal("query"))
			Expect(r4).Should(Equal(""))
		})
	})

	Describe("Parse", func() {
		It("should populate all major components of URL", func() {
			url := Parse("http://user:pass@google.com:80/path?query=query#fragment")
			Expect(url.Scheme).Should(Equal("http"))
			Expect(url.Opaque).Should(Equal("//user:pass@google.com:80/path"))
			Expect(url.Query).Should(Equal("query=query"))
			Expect(url.Fragment).Should(Equal("fragment"))
		})

		It("should separate opaque in to authority & path", func() {
			url := Parse("http://user:pass@google.com:80/path?query=query#fragment")
			Expect(url.Authority).Should(Equal("user:pass@google.com:80"))
			Expect(url.Path).Should(Equal("/path"))
		})

		It("should separate authority in to userinfo, host and port", func() {
			url := Parse("http://user:pass@google.com:80/path?query=query#fragment")
			Expect(url.Host).Should(Equal("google.com"))
			Expect(url.Port).Should(Equal("80"))
			Expect(url.UserInfo).Should(Equal("user:pass"))
		})

		It("should handle empty path", func() {
			url := Parse("http://google.com")
			Expect(url.Host).Should(Equal("google.com"))
			Expect(url.Path).Should(Equal(""))
		})

		It("should handle mailto: url", func() {
			url := Parse("mailto:mike@mike.mike")
			Expect(url.Scheme).Should(Equal("mailto"))
			Expect(url.Opaque).Should(Equal("mike@mike.mike"))
		})

		It("should handle IPv6 url", func() {
			url := Parse("http://[2001:db8:1f70::999:de8:7648:6e8]:9090?test=test")
			Expect(url.Scheme).Should(Equal("http"))
			Expect(url.Opaque).Should(Equal("//[2001:db8:1f70::999:de8:7648:6e8]:9090"))
			Expect(url.Host).Should(Equal("2001:db8:1f70::999:de8:7648:6e8"))
			Expect(url.Port).Should(Equal("9090"))
			Expect(url.Query).Should(Equal("test=test"))
		})

		It("should handle naked host:port", func() {
			url := Parse("google.com:8080")
			fmt.Printf("%#v\n", url)
			Expect(url.Host).Should(Equal("google.com"))
			Expect(url.Port).Should(Equal("8080"))
		})

		It("should handle naked host:port with localhost", func() {
			url := Parse("localhost:8080")
			Expect(url.Host).Should(Equal("google.com"))
			Expect(url.Port).Should(Equal("8080"))
		})

	})
})
