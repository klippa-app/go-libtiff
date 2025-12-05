package raw2tiff_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestRaw2tiff(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "raw2tiff Suite")
}
