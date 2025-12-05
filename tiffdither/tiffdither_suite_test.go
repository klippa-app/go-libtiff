package tiffdither_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestTiffdither(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "tiffdither Suite")
}
