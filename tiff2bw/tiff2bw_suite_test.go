package tiff2bw_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestTiff2bw(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "tiff2bw Suite")
}
