package pal2rgb_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestPal2rgb(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "pal2rgb Suite")
}
