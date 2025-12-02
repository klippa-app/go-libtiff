package fax2tiff_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestFax2tiff(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "fax2tiff Suite")
}
