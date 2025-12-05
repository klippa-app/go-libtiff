package tiffsplit_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestTiffsplit(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "tiffsplit Suite")
}
