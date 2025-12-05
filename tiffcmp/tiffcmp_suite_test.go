package tiffcmp_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestTiffcmp(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "tiffcmp Suite")
}
