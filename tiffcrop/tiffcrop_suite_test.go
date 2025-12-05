package tiffcrop_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestTiffcrop(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "tiffcrop Suite")
}
