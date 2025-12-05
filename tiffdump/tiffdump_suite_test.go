package tiffdump_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestTiffdump(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "tiffdump Suite")
}
