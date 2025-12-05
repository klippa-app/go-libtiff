package tiffset_test

import (
	"bytes"
	"context"

	"github.com/klippa-app/go-libtiff/libtiff"
	"github.com/klippa-app/go-libtiff/tiffset"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/sys"
)

var _ = Describe("tiffset", func() {
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
		err := tiffset.Run(ctx, []string{})
		// For some reason tiffset works different from the other tools.
		// It outputs the help on stderr and returns exit code 1.
		Expect(err).To(MatchError(sys.NewExitError(1)))
		Expect(stdout.String()).To(BeEmpty())
		Expect(stderr.String()).To(ContainSubstring("usage: tiffset"))
	})
})
