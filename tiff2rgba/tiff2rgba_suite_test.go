package tiff2rgba_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestTiff2rgba(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "tiff2rgba Suite")
}
