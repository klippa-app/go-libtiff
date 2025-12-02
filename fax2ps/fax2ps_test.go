package fax2ps_test

import (
	"bytes"
	"context"

	"github.com/klippa-app/go-libtiff/fax2ps"
	"github.com/klippa-app/go-libtiff/libtiff"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/tetratelabs/wazero"
)

var _ = Describe("fax2ps", func() {
	var ctx context.Context
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	BeforeEach(func() {
		ctx = context.Background()
		stdout = bytes.Buffer{}
		stderr = bytes.Buffer{}
		config := &libtiff.Config{
			FSConfig: wazero.NewFSConfig().
				WithDirMount("../testdata", "/testdata").
				WithDirMount("/tmp", "/tmp"),
			Stdout: &stdout,
			Stderr: &stderr,
		}
		ctx = libtiff.ConfigInContext(ctx, config)
	})
	It("shows the help text", func() {
		err := fax2ps.Run(ctx, []string{"-help"})
		Expect(err).To(BeNil())
		Expect(stdout.String()).To(ContainSubstring("usage: fax2ps"))
		Expect(stderr.String()).To(BeEmpty())
	})
})
