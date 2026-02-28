package libtiff_test

import (
	"context"

	"github.com/klippa-app/go-libtiff/libtiff"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("TIFFGetTagListCount", func() {
	ctx := context.Background()

	It("returns the number of tags in the current directory", func() {
		tiffFile, err := instance.TIFFOpenFileFromPath(ctx, "/testdata/lena512color.jpeg.tiff", nil)
		Expect(err).To(BeNil())
		defer tiffFile.Close(ctx)

		count, err := tiffFile.TIFFGetTagListCount(ctx)
		Expect(err).To(BeNil())
		Expect(count).To(BeNumerically(">", 0))
	})

	It("returns zero for a minimal TIFF with no custom tags", func() {
		// TIFFGetTagListCount only counts custom/extended tags, not standard TIFF tags
		// like ImageWidth, BitsPerSample, etc. which are stored in dedicated struct fields.
		readTiff, cleanup := writeMinimalTiff(ctx, func(_ context.Context, _ *libtiff.File) {})
		defer cleanup()

		count, err := readTiff.TIFFGetTagListCount(ctx)
		Expect(err).To(BeNil())
		Expect(count).To(Equal(0))
	})

	It("returns a positive count for a TIFF with custom tags", func() {
		readTiff, cleanup := writeMinimalTiff(ctx, func(ctx context.Context, f *libtiff.File) {
			Expect(f.TIFFSetFieldString(ctx, libtiff.TIFFTAG_SOFTWARE, "test-software")).To(Succeed())
		})
		defer cleanup()

		count, err := readTiff.TIFFGetTagListCount(ctx)
		Expect(err).To(BeNil())
		Expect(count).To(BeNumerically(">", 0))
	})
})

var _ = Describe("TIFFGetTagListEntry", func() {
	ctx := context.Background()

	It("returns valid tag numbers for each entry", func() {
		tiffFile, err := instance.TIFFOpenFileFromPath(ctx, "/testdata/lena512color.jpeg.tiff", nil)
		Expect(err).To(BeNil())
		defer tiffFile.Close(ctx)

		count, err := tiffFile.TIFFGetTagListCount(ctx)
		Expect(err).To(BeNil())

		for i := 0; i < count; i++ {
			tag, err := tiffFile.TIFFGetTagListEntry(ctx, i)
			Expect(err).To(BeNil())
			Expect(tag).To(BeNumerically(">", 0))
		}
	})

	It("returns custom tags that can be looked up", func() {
		// TIFFGetTagListEntry only returns custom tags, not standard tags like IMAGEWIDTH.
		readTiff, cleanup := writeMinimalTiff(ctx, func(ctx context.Context, f *libtiff.File) {
			Expect(f.TIFFSetFieldString(ctx, libtiff.TIFFTAG_SOFTWARE, "test-software")).To(Succeed())
		})
		defer cleanup()

		count, err := readTiff.TIFFGetTagListCount(ctx)
		Expect(err).To(BeNil())
		Expect(count).To(BeNumerically(">", 0))

		// Each entry should have a valid tag number.
		for i := 0; i < count; i++ {
			tag, err := readTiff.TIFFGetTagListEntry(ctx, i)
			Expect(err).To(BeNil())
			Expect(tag).To(BeNumerically(">", 0))
		}
	})
})

var _ = Describe("TIFFFieldWithTag", func() {
	ctx := context.Background()

	It("returns a field descriptor for a known tag", func() {
		tiffFile, err := instance.TIFFOpenFileFromPath(ctx, "/testdata/lena512color.jpeg.tiff", nil)
		Expect(err).To(BeNil())
		defer tiffFile.Close(ctx)

		field, err := tiffFile.TIFFFieldWithTag(ctx, uint32(libtiff.TIFFTAG_IMAGEWIDTH))
		Expect(err).To(BeNil())
		Expect(field).ToNot(BeNil())

		name, err := field.Name(ctx)
		Expect(err).To(BeNil())
		Expect(name).To(Equal("ImageWidth"))
	})

	It("returns nil for an unknown tag", func() {
		tiffFile, err := instance.TIFFOpenFileFromPath(ctx, "/testdata/lena512color.jpeg.tiff", nil)
		Expect(err).To(BeNil())
		defer tiffFile.Close(ctx)

		field, err := tiffFile.TIFFFieldWithTag(ctx, 99999)
		Expect(err).To(BeNil())
		Expect(field).To(BeNil())
	})
})

var _ = Describe("TIFFFieldWithName", func() {
	ctx := context.Background()

	It("returns a field descriptor for a known tag name", func() {
		tiffFile, err := instance.TIFFOpenFileFromPath(ctx, "/testdata/lena512color.jpeg.tiff", nil)
		Expect(err).To(BeNil())
		defer tiffFile.Close(ctx)

		field, err := tiffFile.TIFFFieldWithName(ctx, "ImageWidth")
		Expect(err).To(BeNil())
		Expect(field).ToNot(BeNil())

		tag, err := field.Tag(ctx)
		Expect(err).To(BeNil())
		Expect(tag).To(Equal(uint32(libtiff.TIFFTAG_IMAGEWIDTH)))
	})
})

var _ = Describe("TIFFField", func() {
	ctx := context.Background()

	It("returns field metadata for IMAGEWIDTH", func() {
		tiffFile, err := instance.TIFFOpenFileFromPath(ctx, "/testdata/lena512color.jpeg.tiff", nil)
		Expect(err).To(BeNil())
		defer tiffFile.Close(ctx)

		field, err := tiffFile.TIFFFieldWithTag(ctx, uint32(libtiff.TIFFTAG_IMAGEWIDTH))
		Expect(err).To(BeNil())
		Expect(field).ToNot(BeNil())

		tag, err := field.Tag(ctx)
		Expect(err).To(BeNil())
		Expect(tag).To(Equal(uint32(256)))

		name, err := field.Name(ctx)
		Expect(err).To(BeNil())
		Expect(name).To(Equal("ImageWidth"))

		dataType, err := field.DataType(ctx)
		Expect(err).To(BeNil())
		// ImageWidth is TIFF_LONG or TIFF_SHORT
		Expect(dataType).To(Or(Equal(libtiff.TIFF_LONG), Equal(libtiff.TIFF_SHORT)))

		isAnon, err := field.IsAnonymous(ctx)
		Expect(err).To(BeNil())
		Expect(isAnon).To(BeFalse())
	})

	It("returns read/write count for BitsPerSample", func() {
		tiffFile, err := instance.TIFFOpenFileFromPath(ctx, "/testdata/lena512color.jpeg.tiff", nil)
		Expect(err).To(BeNil())
		defer tiffFile.Close(ctx)

		field, err := tiffFile.TIFFFieldWithTag(ctx, uint32(libtiff.TIFFTAG_BITSPERSAMPLE))
		Expect(err).To(BeNil())
		Expect(field).ToNot(BeNil())

		readCount, err := field.ReadCount(ctx)
		Expect(err).To(BeNil())
		// BitsPerSample has TIFF_VARIABLE read count (-1) in libtiff
		Expect(readCount).To(Or(Equal(-1), BeNumerically(">", 0)))

		writeCount, err := field.WriteCount(ctx)
		Expect(err).To(BeNil())
		Expect(writeCount).To(Or(Equal(-1), BeNumerically(">", 0)))
	})

	It("returns SetGetSize for a known tag", func() {
		tiffFile, err := instance.TIFFOpenFileFromPath(ctx, "/testdata/lena512color.jpeg.tiff", nil)
		Expect(err).To(BeNil())
		defer tiffFile.Close(ctx)

		field, err := tiffFile.TIFFFieldWithTag(ctx, uint32(libtiff.TIFFTAG_IMAGEWIDTH))
		Expect(err).To(BeNil())
		Expect(field).ToNot(BeNil())

		size, err := field.SetGetSize(ctx)
		Expect(err).To(BeNil())
		Expect(size).To(BeNumerically(">", 0))
	})
})

var _ = Describe("Tag enumeration", func() {
	ctx := context.Background()

	It("can enumerate all tags and look up their metadata", func() {
		tiffFile, err := instance.TIFFOpenFileFromPath(ctx, "/testdata/lena512color.jpeg.tiff", nil)
		Expect(err).To(BeNil())
		defer tiffFile.Close(ctx)

		count, err := tiffFile.TIFFGetTagListCount(ctx)
		Expect(err).To(BeNil())

		for i := 0; i < count; i++ {
			tagNum, err := tiffFile.TIFFGetTagListEntry(ctx, i)
			Expect(err).To(BeNil())

			field, err := tiffFile.TIFFFieldWithTag(ctx, tagNum)
			Expect(err).To(BeNil())
			// Some tags may not have field descriptors
			if field == nil {
				continue
			}

			name, err := field.Name(ctx)
			Expect(err).To(BeNil())
			Expect(name).ToNot(BeEmpty())
		}
	})
})
