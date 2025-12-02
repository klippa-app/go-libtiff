package fax2ps_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestFax2ps(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "fax2ps Suite")
}
