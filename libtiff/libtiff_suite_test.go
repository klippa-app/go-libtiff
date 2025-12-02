package libtiff_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestLibtiff(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Libtiff Suite")
}
