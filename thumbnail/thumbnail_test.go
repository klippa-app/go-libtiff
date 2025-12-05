package thumbnail_test

import (
	"bytes"
	"context"

	"github.com/klippa-app/go-libtiff/libtiff"
	"github.com/klippa-app/go-libtiff/thumbnail"
	"github.com/tetratelabs/wazero/sys"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/tetratelabs/wazero"
)

var _ = Describe("thumbnail", func() {
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
		err := thumbnail.Run(ctx, []string{"-help"})
		// For some reason thumbnail works different from the other tools.
		// It outputs the help on stderr and returns exit code 1.
		Expect(err).To(MatchError(sys.NewExitError(1)))
		Expect(stdout.String()).To(BeEmpty())
		Expect(stderr.String()).To(ContainSubstring("usage: thumbnail"))
	})
})
