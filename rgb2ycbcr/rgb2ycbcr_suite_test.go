package rgb2ycbcr_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestRgb2ycbcr(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "rgb2ycbcr Suite")
}
