package mkg3states_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestMkg3states(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "mkg3states Suite")
}
