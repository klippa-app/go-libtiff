package tiffset_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestTiffset(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "tiffset Suite")
}
