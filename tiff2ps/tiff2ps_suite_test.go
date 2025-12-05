package tiff2ps_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestTiff2ps(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "tiff2ps Suite")
}
