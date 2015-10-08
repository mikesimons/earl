package earl_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestEarl(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Earl Suite")
}
