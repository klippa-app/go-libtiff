package tiffcp_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestTiffcp(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "tiffcp Suite")
}
