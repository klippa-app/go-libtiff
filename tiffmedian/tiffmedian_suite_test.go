package tiffmedian_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestTiffmedian(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "tiffmedian Suite")
}
