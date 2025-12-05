package tiff2bw_test

import (
	"bytes"
	"context"

	"github.com/klippa-app/go-libtiff/libtiff"
	"github.com/klippa-app/go-libtiff/tiff2bw"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/tetratelabs/wazero"
)

var _ = Describe("tiff2bw", func() {
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
		err := tiff2bw.Run(ctx, []string{"-help"})
		Expect(err).To(BeNil())
		Expect(stdout.String()).To(ContainSubstring("usage: tiff2bw"))
		Expect(stderr.String()).To(BeEmpty())
	})
})
