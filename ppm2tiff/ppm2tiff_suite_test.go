package ppm2tiff_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestPpm2tiff(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "ppm2tiff Suite")
}
