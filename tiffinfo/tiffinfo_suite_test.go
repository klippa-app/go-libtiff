package tiffinfo_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestTiffinfo(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "tiffinfo Suite")
}
