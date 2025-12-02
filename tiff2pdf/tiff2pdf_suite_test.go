package tiff2pdf_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestTiff2pdf(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "tiff2pdf Suite")
}
