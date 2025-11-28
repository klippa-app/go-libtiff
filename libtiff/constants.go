package libtiff

type TIFFTAG uint32

// https://gitlab.com/libtiff/libtiff/-/blob/master/libtiff/tiff.h

var (
	TIFFTAG_SUBFILETYPE     = TIFFTAG(254)   /* subfile data descriptor */
	FILETYPE_REDUCEDIMAGE   = TIFFTAG(0x1)   /* reduced resolution version */
	FILETYPE_PAGE           = TIFFTAG(0x2)   /* one page of many */
	FILETYPE_MASK           = TIFFTAG(0x4)   /* transparency mask */
	TIFFTAG_OSUBFILETYPE    = TIFFTAG(255)   /* +kind of data in subfile */
	OFILETYPE_IMAGE         = TIFFTAG(1)     /* full resolution image data */
	OFILETYPE_REDUCEDIMAGE  = TIFFTAG(2)     /* reduced size image data */
	OFILETYPE_PAGE          = TIFFTAG(3)     /* one page of many */
	TIFFTAG_IMAGEWIDTH      = TIFFTAG(256)   /* image width in pixels */
	TIFFTAG_IMAGELENGTH     = TIFFTAG(257)   /* image height in pixels */
	TIFFTAG_BITSPERSAMPLE   = TIFFTAG(258)   /* bits per channel (sample) */
	TIFFTAG_COMPRESSION     = TIFFTAG(259)   /* data compression technique */
	COMPRESSION_NONE        = TIFFTAG(1)     /* dump mode */
	COMPRESSION_CCITTRLE    = TIFFTAG(2)     /* CCITT modified Huffman RLE */
	COMPRESSION_CCITTFAX3   = TIFFTAG(3)     /* CCITT Group 3 fax encoding */
	COMPRESSION_CCITT_T4    = TIFFTAG(3)     /* CCITT T.4 (TIFF 6 name) */
	COMPRESSION_CCITTFAX4   = TIFFTAG(4)     /* CCITT Group 4 fax encoding */
	COMPRESSION_CCITT_T6    = TIFFTAG(4)     /* CCITT T.6 (TIFF 6 name) */
	COMPRESSION_LZW         = TIFFTAG(5)     /* Lempel-Ziv  & Welch */
	COMPRESSION_OJPEG       = TIFFTAG(6)     /* !6.0 JPEG */
	COMPRESSION_JPEG        = TIFFTAG(7)     /* %JPEG DCT compression */
	COMPRESSION_T85         = TIFFTAG(9)     /* !TIFF/FX T.85 JBIG compression */
	COMPRESSION_T43         = TIFFTAG(10)    /* !TIFF/FX T.43 colour by layered JBIG compression */
	COMPRESSION_NEXT        = TIFFTAG(32766) /* NeXT 2-bit RLE */
	COMPRESSION_CCITTRLEW   = TIFFTAG(32771) /* #1 w/ word alignment */
	COMPRESSION_PACKBITS    = TIFFTAG(32773) /* Macintosh RLE */
	COMPRESSION_THUNDERSCAN = TIFFTAG(32809) /* ThunderScan RLE */
	/* codes 32895-32898 are reserved for ANSI IT8 TIFF/IT <dkelly@apago.com) */
	COMPRESSION_IT8CTPAD = TIFFTAG(32895) /* IT8 CT w/padding */
	COMPRESSION_IT8LW    = TIFFTAG(32896) /* IT8 Linework RLE */
	COMPRESSION_IT8MP    = TIFFTAG(32897) /* IT8 Monochrome picture */
	COMPRESSION_IT8BL    = TIFFTAG(32898) /* IT8 Binary line art */
	/* compression codes 32908-32911 are reserved for Pixar */
	COMPRESSION_PIXARFILM     = TIFFTAG(32908) /* Pixar companded 10bit LZW */
	COMPRESSION_PIXARLOG      = TIFFTAG(32909) /* Pixar companded 11bit ZIP */
	COMPRESSION_DEFLATE       = TIFFTAG(32946) /* Deflate compression, legacy tag */
	COMPRESSION_ADOBE_DEFLATE = TIFFTAG(8)     /* Deflate compression, as recognized by Adobe */
	/* compression code 32947 is reserved for Oceana Matrix <dev@oceana.com> */
	COMPRESSION_DCS      = TIFFTAG(32947) /* Kodak DCS encoding */
	COMPRESSION_JBIG     = TIFFTAG(34661) /* ISO JBIG */
	COMPRESSION_SGILOG   = TIFFTAG(34676) /* SGI Log Luminance RLE */
	COMPRESSION_SGILOG24 = TIFFTAG(34677) /* SGI Log 24-bit packed */
	COMPRESSION_JP2000   = TIFFTAG(34712) /* Leadtools JPEG2000 */
	COMPRESSION_LERC     = TIFFTAG(34887) /* ESRI Lerc codec: https://github.com/Esri/lerc */
	/* compression codes 34887-34889 are reserved for ESRI */
	COMPRESSION_LZMA               = TIFFTAG(34925) /* LZMA2 */
	COMPRESSION_ZSTD               = TIFFTAG(50000) /* ZSTD: WARNING not registered in Adobe-maintained registry */
	COMPRESSION_WEBP               = TIFFTAG(50001) /* WEBP: WARNING not registered in Adobe-maintained registry */
	COMPRESSION_JXL                = TIFFTAG(50002) /* JPEGXL: WARNING not registered in Adobe-maintained registry */
	COMPRESSION_JXL_DNG_1_7        = TIFFTAG(52546) /* JPEGXL from DNG 1.7 specification */
	TIFFTAG_PHOTOMETRIC            = TIFFTAG(262)   /* photometric interpretation */
	PHOTOMETRIC_MINISWHITE         = TIFFTAG(0)     /* min value is white */
	PHOTOMETRIC_MINISBLACK         = TIFFTAG(1)     /* min value is black */
	PHOTOMETRIC_RGB                = TIFFTAG(2)     /* RGB color model */
	PHOTOMETRIC_PALETTE            = TIFFTAG(3)     /* color map indexed */
	PHOTOMETRIC_MASK               = TIFFTAG(4)     /* $holdout mask */
	PHOTOMETRIC_SEPARATED          = TIFFTAG(5)     /* !color separations */
	PHOTOMETRIC_YCBCR              = TIFFTAG(6)     /* !CCIR 601 */
	PHOTOMETRIC_CIELAB             = TIFFTAG(8)     /* !1976 CIE L*a*b* */
	PHOTOMETRIC_ICCLAB             = TIFFTAG(9)     /* ICC L*a*b* [Adobe TIFF Technote 4] */
	PHOTOMETRIC_ITULAB             = TIFFTAG(10)    /* ITU L*a*b* */
	PHOTOMETRIC_CFA                = TIFFTAG(32803) /* color filter array */
	PHOTOMETRIC_LOGL               = TIFFTAG(32844) /* CIE Log2(L) */
	PHOTOMETRIC_LOGLUV             = TIFFTAG(32845) /* CIE Log2(L) (u',v') */
	TIFFTAG_THRESHHOLDING          = TIFFTAG(263)   /* +thresholding used on data */
	THRESHHOLD_BILEVEL             = TIFFTAG(1)     /* b&w art scan */
	THRESHHOLD_HALFTONE            = TIFFTAG(2)     /* or dithered scan */
	THRESHHOLD_ERRORDIFFUSE        = TIFFTAG(3)     /* usually floyd-steinberg */
	TIFFTAG_CELLWIDTH              = TIFFTAG(264)   /* +dithering matrix width */
	TIFFTAG_CELLLENGTH             = TIFFTAG(265)   /* +dithering matrix height */
	TIFFTAG_FILLORDER              = TIFFTAG(266)   /* data order within a byte */
	FILLORDER_MSB2LSB              = TIFFTAG(1)     /* most significant -> least */
	FILLORDER_LSB2MSB              = TIFFTAG(2)     /* least significant -> most */
	TIFFTAG_DOCUMENTNAME           = TIFFTAG(269)   /* name of doc. image is from */
	TIFFTAG_IMAGEDESCRIPTION       = TIFFTAG(270)   /* info about image */
	TIFFTAG_MAKE                   = TIFFTAG(271)   /* scanner manufacturer name */
	TIFFTAG_MODEL                  = TIFFTAG(272)   /* scanner model name/number */
	TIFFTAG_STRIPOFFSETS           = TIFFTAG(273)   /* offsets to data strips */
	TIFFTAG_ORIENTATION            = TIFFTAG(274)   /* +image orientation */
	ORIENTATION_TOPLEFT            = TIFFTAG(1)     /* row 0 top, col 0 lhs */
	ORIENTATION_TOPRIGHT           = TIFFTAG(2)     /* row 0 top, col 0 rhs */
	ORIENTATION_BOTRIGHT           = TIFFTAG(3)     /* row 0 bottom, col 0 rhs */
	ORIENTATION_BOTLEFT            = TIFFTAG(4)     /* row 0 bottom, col 0 lhs */
	ORIENTATION_LEFTTOP            = TIFFTAG(5)     /* row 0 lhs, col 0 top */
	ORIENTATION_RIGHTTOP           = TIFFTAG(6)     /* row 0 rhs, col 0 top */
	ORIENTATION_RIGHTBOT           = TIFFTAG(7)     /* row 0 rhs, col 0 bottom */
	ORIENTATION_LEFTBOT            = TIFFTAG(8)     /* row 0 lhs, col 0 bottom */
	TIFFTAG_SAMPLESPERPIXEL        = TIFFTAG(277)   /* samples per pixel */
	TIFFTAG_ROWSPERSTRIP           = TIFFTAG(278)   /* rows per strip of data */
	TIFFTAG_STRIPBYTECOUNTS        = TIFFTAG(279)   /* bytes counts for strips */
	TIFFTAG_MINSAMPLEVALUE         = TIFFTAG(280)   /* +minimum sample value */
	TIFFTAG_MAXSAMPLEVALUE         = TIFFTAG(281)   /* +maximum sample value */
	TIFFTAG_XRESOLUTION            = TIFFTAG(282)   /* pixels/resolution in x */
	TIFFTAG_YRESOLUTION            = TIFFTAG(283)   /* pixels/resolution in y */
	TIFFTAG_PLANARCONFIG           = TIFFTAG(284)   /* storage organization */
	PLANARCONFIG_CONTIG            = TIFFTAG(1)     /* single image plane */
	PLANARCONFIG_SEPARATE          = TIFFTAG(2)     /* separate planes of data */
	TIFFTAG_PAGENAME               = TIFFTAG(285)   /* page name image is from */
	TIFFTAG_XPOSITION              = TIFFTAG(286)   /* x page offset of image lhs */
	TIFFTAG_YPOSITION              = TIFFTAG(287)   /* y page offset of image lhs */
	TIFFTAG_FREEOFFSETS            = TIFFTAG(288)   /* +byte offset to free block */
	TIFFTAG_FREEBYTECOUNTS         = TIFFTAG(289)   /* +sizes of free blocks */
	TIFFTAG_GRAYRESPONSEUNIT       = TIFFTAG(290)   /* $gray scale curve accuracy */
	GRAYRESPONSEUNIT_10S           = TIFFTAG(1)     /* tenths of a unit */
	GRAYRESPONSEUNIT_100S          = TIFFTAG(2)     /* hundredths of a unit */
	GRAYRESPONSEUNIT_1000S         = TIFFTAG(3)     /* thousandths of a unit */
	GRAYRESPONSEUNIT_10000S        = TIFFTAG(4)     /* ten-thousandths of a unit */
	GRAYRESPONSEUNIT_100000S       = TIFFTAG(5)     /* hundred-thousandths */
	TIFFTAG_GRAYRESPONSECURVE      = TIFFTAG(291)   /* $gray scale response curve */
	TIFFTAG_GROUP3OPTIONS          = TIFFTAG(292)   /* 32 flag bits */
	TIFFTAG_T4OPTIONS              = TIFFTAG(292)   /* TIFF 6.0 proper name alias */
	GROUP3OPT_2DENCODING           = TIFFTAG(0x1)   /* 2-dimensional coding */
	GROUP3OPT_UNCOMPRESSED         = TIFFTAG(0x2)   /* data not compressed */
	GROUP3OPT_FILLBITS             = TIFFTAG(0x4)   /* fill to byte boundary */
	TIFFTAG_GROUP4OPTIONS          = TIFFTAG(293)   /* 32 flag bits */
	TIFFTAG_T6OPTIONS              = TIFFTAG(293)   /* TIFF 6.0 proper name */
	GROUP4OPT_UNCOMPRESSED         = TIFFTAG(0x2)   /* data not compressed */
	TIFFTAG_RESOLUTIONUNIT         = TIFFTAG(296)   /* units of resolutions */
	RESUNIT_NONE                   = TIFFTAG(1)     /* no meaningful units */
	RESUNIT_INCH                   = TIFFTAG(2)     /* english */
	RESUNIT_CENTIMETER             = TIFFTAG(3)     /* metric */
	TIFFTAG_PAGENUMBER             = TIFFTAG(297)   /* page numbers of multi-page */
	TIFFTAG_COLORRESPONSEUNIT      = TIFFTAG(300)   /* $color curve accuracy */
	COLORRESPONSEUNIT_10S          = TIFFTAG(1)     /* tenths of a unit */
	COLORRESPONSEUNIT_100S         = TIFFTAG(2)     /* hundredths of a unit */
	COLORRESPONSEUNIT_1000S        = TIFFTAG(3)     /* thousandths of a unit */
	COLORRESPONSEUNIT_10000S       = TIFFTAG(4)     /* ten-thousandths of a unit */
	COLORRESPONSEUNIT_100000S      = TIFFTAG(5)     /* hundred-thousandths */
	TIFFTAG_TRANSFERFUNCTION       = TIFFTAG(301)   /* !colorimetry info */
	TIFFTAG_SOFTWARE               = TIFFTAG(305)   /* name & release */
	TIFFTAG_DATETIME               = TIFFTAG(306)   /* creation date and time */
	TIFFTAG_ARTIST                 = TIFFTAG(315)   /* creator of image */
	TIFFTAG_HOSTCOMPUTER           = TIFFTAG(316)   /* machine where created */
	TIFFTAG_PREDICTOR              = TIFFTAG(317)   /* prediction scheme w/ LZW */
	PREDICTOR_NONE                 = TIFFTAG(1)     /* no prediction scheme used */
	PREDICTOR_HORIZONTAL           = TIFFTAG(2)     /* horizontal differencing */
	PREDICTOR_FLOATINGPOINT        = TIFFTAG(3)     /* floating point predictor */
	TIFFTAG_WHITEPOINT             = TIFFTAG(318)   /* image white point */
	TIFFTAG_PRIMARYCHROMATICITIES  = TIFFTAG(319)   /* !primary chromaticities */
	TIFFTAG_COLORMAP               = TIFFTAG(320)   /* RGB map for palette image */
	TIFFTAG_HALFTONEHINTS          = TIFFTAG(321)   /* !highlight+shadow info */
	TIFFTAG_TILEWIDTH              = TIFFTAG(322)   /* !tile width in pixels */
	TIFFTAG_TILELENGTH             = TIFFTAG(323)   /* !tile height in pixels */
	TIFFTAG_TILEOFFSETS            = TIFFTAG(324)   /* !offsets to data tiles */
	TIFFTAG_TILEBYTECOUNTS         = TIFFTAG(325)   /* !byte counts for tiles */
	TIFFTAG_BADFAXLINES            = TIFFTAG(326)   /* lines w/ wrong pixel count */
	TIFFTAG_CLEANFAXDATA           = TIFFTAG(327)   /* regenerated line info */
	CLEANFAXDATA_CLEAN             = TIFFTAG(0)     /* no errors detected */
	CLEANFAXDATA_REGENERATED       = TIFFTAG(1)     /* receiver regenerated lines */
	CLEANFAXDATA_UNCLEAN           = TIFFTAG(2)     /* uncorrected errors exist */
	TIFFTAG_CONSECUTIVEBADFAXLINES = TIFFTAG(328)   /* max consecutive bad lines */
	TIFFTAG_SUBIFD                 = TIFFTAG(330)   /* subimage descriptors */
	TIFFTAG_INKSET                 = TIFFTAG(332)   /* !inks in separated image */
	INKSET_CMYK                    = TIFFTAG(1)     /* !cyan-magenta-yellow-black color */
	INKSET_MULTIINK                = TIFFTAG(2)     /* !multi-ink or hi-fi color */
	TIFFTAG_INKNAMES               = TIFFTAG(333)   /* !ascii names of inks */
	TIFFTAG_NUMBEROFINKS           = TIFFTAG(334)   /* !number of inks */
	TIFFTAG_DOTRANGE               = TIFFTAG(336)   /* !0% and 100% dot codes */
	TIFFTAG_TARGETPRINTER          = TIFFTAG(337)   /* !separation target */
	TIFFTAG_EXTRASAMPLES           = TIFFTAG(338)   /* !info about extra samples */
	EXTRASAMPLE_UNSPECIFIED        = TIFFTAG(0)     /* !unspecified data */
	EXTRASAMPLE_ASSOCALPHA         = TIFFTAG(1)     /* !associated alpha data */
	EXTRASAMPLE_UNASSALPHA         = TIFFTAG(2)     /* !unassociated alpha data */
	TIFFTAG_SAMPLEFORMAT           = TIFFTAG(339)   /* !data sample format */
	SAMPLEFORMAT_UINT              = TIFFTAG(1)     /* !unsigned integer data */
	SAMPLEFORMAT_INT               = TIFFTAG(2)     /* !signed integer data */
	SAMPLEFORMAT_IEEEFP            = TIFFTAG(3)     /* !IEEE floating point data */
	SAMPLEFORMAT_VOID              = TIFFTAG(4)     /* !untyped data */
	SAMPLEFORMAT_COMPLEXINT        = TIFFTAG(5)     /* !complex signed int */
	SAMPLEFORMAT_COMPLEXIEEEFP     = TIFFTAG(6)     /* !complex ieee floating */
	TIFFTAG_SMINSAMPLEVALUE        = TIFFTAG(340)   /* !variable MinSampleValue */
	TIFFTAG_SMAXSAMPLEVALUE        = TIFFTAG(341)   /* !variable MaxSampleValue */
	TIFFTAG_CLIPPATH               = TIFFTAG(343)   /* %ClipPath [Adobe TIFF technote 2] */
	TIFFTAG_XCLIPPATHUNITS         = TIFFTAG(344)   /* %XClipPathUnits [Adobe TIFF technote 2] */
	TIFFTAG_YCLIPPATHUNITS         = TIFFTAG(345)   /* %YClipPathUnits [Adobe TIFF technote 2] */
	TIFFTAG_INDEXED                = TIFFTAG(346)   /* %Indexed [Adobe TIFF Technote 3] */
	TIFFTAG_JPEGTABLES             = TIFFTAG(347)   /* %JPEG table stream */
	TIFFTAG_OPIPROXY               = TIFFTAG(351)   /* %OPI Proxy [Adobe TIFF technote] */
	/* Tags 400-435 are from the TIFF/FX spec */
	TIFFTAG_GLOBALPARAMETERSIFD = TIFFTAG(400)    /* ! */
	TIFFTAG_PROFILETYPE         = TIFFTAG(401)    /* ! */
	PROFILETYPE_UNSPECIFIED     = TIFFTAG(0)      /* ! */
	PROFILETYPE_G3_FAX          = TIFFTAG(1)      /* ! */
	TIFFTAG_FAXPROFILE          = TIFFTAG(402)    /* ! */
	FAXPROFILE_S                = TIFFTAG(1)      /* !TIFF/FX FAX profile S */
	FAXPROFILE_F                = TIFFTAG(2)      /* !TIFF/FX FAX profile F */
	FAXPROFILE_J                = TIFFTAG(3)      /* !TIFF/FX FAX profile J */
	FAXPROFILE_C                = TIFFTAG(4)      /* !TIFF/FX FAX profile C */
	FAXPROFILE_L                = TIFFTAG(5)      /* !TIFF/FX FAX profile L */
	FAXPROFILE_M                = TIFFTAG(6)      /* !TIFF/FX FAX profile LM */
	TIFFTAG_CODINGMETHODS       = TIFFTAG(403)    /* !TIFF/FX coding methods */
	CODINGMETHODS_T4_1D         = TIFFTAG(1 << 1) /* !T.4 1D */
	CODINGMETHODS_T4_2D         = TIFFTAG(1 << 2) /* !T.4 2D */
	CODINGMETHODS_T6            = TIFFTAG(1 << 3) /* !T.6 */
	CODINGMETHODS_T85           = TIFFTAG(1 << 4) /* !T.85 JBIG */
	CODINGMETHODS_T42           = TIFFTAG(1 << 5) /* !T.42 JPEG */
	CODINGMETHODS_T43           = TIFFTAG(1 << 6) /* !T.43 colour by layered JBIG */
	TIFFTAG_VERSIONYEAR         = TIFFTAG(404)    /* !TIFF/FX version year */
	TIFFTAG_MODENUMBER          = TIFFTAG(405)    /* !TIFF/FX mode number */
	TIFFTAG_DECODE              = TIFFTAG(433)    /* !TIFF/FX decode */
	TIFFTAG_IMAGEBASECOLOR      = TIFFTAG(434)    /* !TIFF/FX image base colour */
	TIFFTAG_T82OPTIONS          = TIFFTAG(435)    /* !TIFF/FX T.82 options */
	/*
	 * Tags 512-521 are obsoleted by Technical Note #2 which specifies a
	 * revised JPEG-in-TIFF scheme.
	 */
	TIFFTAG_JPEGPROC               = TIFFTAG(512)   /* !JPEG processing algorithm */
	JPEGPROC_BASELINE              = TIFFTAG(1)     /* !baseline sequential */
	JPEGPROC_LOSSLESS              = TIFFTAG(14)    /* !Huffman coded lossless */
	TIFFTAG_JPEGIFOFFSET           = TIFFTAG(513)   /* !pointer to SOI marker */
	TIFFTAG_JPEGIFBYTECOUNT        = TIFFTAG(514)   /* !JFIF stream length */
	TIFFTAG_JPEGRESTARTINTERVAL    = TIFFTAG(515)   /* !restart interval length */
	TIFFTAG_JPEGLOSSLESSPREDICTORS = TIFFTAG(517)   /* !lossless proc predictor */
	TIFFTAG_JPEGPOINTTRANSFORM     = TIFFTAG(518)   /* !lossless point transform */
	TIFFTAG_JPEGQTABLES            = TIFFTAG(519)   /* !Q matrix offsets */
	TIFFTAG_JPEGDCTABLES           = TIFFTAG(520)   /* !DCT table offsets */
	TIFFTAG_JPEGACTABLES           = TIFFTAG(521)   /* !AC coefficient offsets */
	TIFFTAG_YCBCRCOEFFICIENTS      = TIFFTAG(529)   /* !RGB -> YCbCr transform */
	TIFFTAG_YCBCRSUBSAMPLING       = TIFFTAG(530)   /* !YCbCr subsampling factors */
	TIFFTAG_YCBCRPOSITIONING       = TIFFTAG(531)   /* !subsample positioning */
	YCBCRPOSITION_CENTERED         = TIFFTAG(1)     /* !as in PostScript Level 2 */
	YCBCRPOSITION_COSITED          = TIFFTAG(2)     /* !as in CCIR 601-1 */
	TIFFTAG_REFERENCEBLACKWHITE    = TIFFTAG(532)   /* !colorimetry info */
	TIFFTAG_STRIPROWCOUNTS         = TIFFTAG(559)   /* !TIFF/FX strip row counts */
	TIFFTAG_XMLPACKET              = TIFFTAG(700)   /* %XML packet [Adobe XMP Specification, January 2004 */
	TIFFTAG_OPIIMAGEID             = TIFFTAG(32781) /* %OPI ImageID [Adobe TIFF technote] */
	/* For eiStream Annotation Specification, Version 1.00.06 see
	 * http://web.archive.org/web/20050309141348/http://www.kofile.com/support%20pro/faqs/annospec.htm */
	TIFFTAG_TIFFANNOTATIONDATA = TIFFTAG(32932)
	/* tags 32952-32956 are private tags registered to Island Graphics */
	TIFFTAG_REFPTS            = TIFFTAG(32953) /* image reference points */
	TIFFTAG_REGIONTACKPOINT   = TIFFTAG(32954) /* region-xform tack point */
	TIFFTAG_REGIONWARPCORNERS = TIFFTAG(32955) /* warp quadrilateral */
	TIFFTAG_REGIONAFFINE      = TIFFTAG(32956) /* affine transformation mat */
	/* tags 32995-32999 are private tags registered to SGI */
	TIFFTAG_MATTEING   = TIFFTAG(32995) /* $use ExtraSamples */
	TIFFTAG_DATATYPE   = TIFFTAG(32996) /* $use SampleFormat */
	TIFFTAG_IMAGEDEPTH = TIFFTAG(32997) /* z depth of image */
	TIFFTAG_TILEDEPTH  = TIFFTAG(32998) /* z depth/data tile */
	/* tags 33300-33309 are private tags registered to Pixar */
	/*
	 * TIFFTAG_PIXAR_IMAGEFULLWIDTH and TIFFTAG_PIXAR_IMAGEFULLLENGTH
	 * are set when an image has been cropped out of a larger image.
	 * They reflect the size of the original uncropped image.
	 * The TIFFTAG_XPOSITION and TIFFTAG_YPOSITION can be used
	 * to determine the position of the smaller image in the larger one.
	 */
	TIFFTAG_PIXAR_IMAGEFULLWIDTH  = TIFFTAG(33300) /* full image size in x */
	TIFFTAG_PIXAR_IMAGEFULLLENGTH = TIFFTAG(33301) /* full image size in y */
	/* Tags 33302-33306 are used to identify special image modes and data
	 * used by Pixar's texture formats.
	 */
	TIFFTAG_PIXAR_TEXTUREFORMAT        = TIFFTAG(33302) /* texture map format */
	TIFFTAG_PIXAR_WRAPMODES            = TIFFTAG(33303) /* s & t wrap modes */
	TIFFTAG_PIXAR_FOVCOT               = TIFFTAG(33304) /* cotan(fov) for env. maps */
	TIFFTAG_PIXAR_MATRIX_WORLDTOSCREEN = TIFFTAG(33305)
	TIFFTAG_PIXAR_MATRIX_WORLDTOCAMERA = TIFFTAG(33306)
	/* tag 33405 is a private tag registered to Eastman Kodak */
	TIFFTAG_WRITERSERIALNUMBER  = TIFFTAG(33405) /* device serial number */
	TIFFTAG_CFAREPEATPATTERNDIM = TIFFTAG(33421) /* (alias for TIFFTAG_EP_CFAREPEATPATTERNDIM)*/
	TIFFTAG_CFAPATTERN          = TIFFTAG(33422) /* (alias for TIFFTAG_EP_CFAPATTERN) */
	TIFFTAG_BATTERYLEVEL        = TIFFTAG(33423) /* (alias for TIFFTAG_EP_BATTERYLEVEL) */
	/* tag 33432 is listed in the 6.0 spec w/ unknown ownership */
	TIFFTAG_COPYRIGHT = TIFFTAG(33432) /* copyright string */
	/* Tags 33445-33452 are used for Molecular Dynamics GEL fileformat,
	 * see http://research.stowers-institute.org/mcm/efg/ScientificSoftware/Utility/TiffTags/GEL-FileFormat.pdf
	 * (2023: the above web site is unavailable but tags are explained briefly at
	 * https://www.awaresystems.be/imaging/tiff/tifftags/docs/gel.html
	 */
	TIFFTAG_MD_FILETAG    = TIFFTAG(33445) /* Specifies the pixel data format encoding in the GEL file format. */
	TIFFTAG_MD_SCALEPIXEL = TIFFTAG(33446) /* scale factor */
	TIFFTAG_MD_COLORTABLE = TIFFTAG(33447) /* conversion from 16bit to 8bit */
	TIFFTAG_MD_LABNAME    = TIFFTAG(33448) /* name of the lab that scanned this file. */
	TIFFTAG_MD_SAMPLEINFO = TIFFTAG(33449) /* information about the scanned GEL sample */
	TIFFTAG_MD_PREPDATE   = TIFFTAG(33450) /* information about the date the sample was prepared YY/MM/DD */
	TIFFTAG_MD_PREPTIME   = TIFFTAG(33451) /* information about the time the sample was prepared HH:MM*/
	TIFFTAG_MD_FILEUNITS  = TIFFTAG(33452) /* Units for data in this file, as used in the GEL file format. */
	/* IPTC TAG from RichTIFF specifications */
	TIFFTAG_RICHTIFFIPTC               = TIFFTAG(33723)
	TIFFTAG_INGR_PACKET_DATA_TAG       = TIFFTAG(33918) /* Intergraph Application specific storage. */
	TIFFTAG_INGR_FLAG_REGISTERS        = TIFFTAG(33919) /* Intergraph Application specific flags. */
	TIFFTAG_IRASB_TRANSORMATION_MATRIX = TIFFTAG(33920) /* Originally part of Intergraph's GeoTIFF tags, but likely understood by IrasB only. */
	TIFFTAG_MODELTIEPOINTTAG           = TIFFTAG(33922) /* GeoTIFF */
	/* 34016-34029 are reserved for ANSI IT8 TIFF/IT <dkelly@apago.com) */
	TIFFTAG_IT8SITE                     = TIFFTAG(34016) /* site name */
	TIFFTAG_IT8COLORSEQUENCE            = TIFFTAG(34017) /* color seq. [RGB,CMYK,etc] */
	TIFFTAG_IT8HEADER                   = TIFFTAG(34018) /* DDES Header */
	TIFFTAG_IT8RASTERPADDING            = TIFFTAG(34019) /* raster scanline padding */
	TIFFTAG_IT8BITSPERRUNLENGTH         = TIFFTAG(34020) /* # of bits in short run */
	TIFFTAG_IT8BITSPEREXTENDEDRUNLENGTH = TIFFTAG(34021) /* # of bits in long run */
	TIFFTAG_IT8COLORTABLE               = TIFFTAG(34022) /* LW colortable */
	TIFFTAG_IT8IMAGECOLORINDICATOR      = TIFFTAG(34023) /* BP/BL image color switch */
	TIFFTAG_IT8BKGCOLORINDICATOR        = TIFFTAG(34024) /* BP/BL bg color switch */
	TIFFTAG_IT8IMAGECOLORVALUE          = TIFFTAG(34025) /* BP/BL image color value */
	TIFFTAG_IT8BKGCOLORVALUE            = TIFFTAG(34026) /* BP/BL bg color value */
	TIFFTAG_IT8PIXELINTENSITYRANGE      = TIFFTAG(34027) /* MP pixel intensity value */
	TIFFTAG_IT8TRANSPARENCYINDICATOR    = TIFFTAG(34028) /* HC transparency switch */
	TIFFTAG_IT8COLORCHARACTERIZATION    = TIFFTAG(34029) /* color character. table */
	TIFFTAG_IT8HCUSAGE                  = TIFFTAG(34030) /* HC usage indicator */
	TIFFTAG_IT8TRAPINDICATOR            = TIFFTAG(34031) /* Trapping indicator (untrapped=0, trapped=1) */
	TIFFTAG_IT8CMYKEQUIVALENT           = TIFFTAG(34032) /* CMYK color equivalents */
	/* tags 34232-34236 are private tags registered to Texas Instruments */
	TIFFTAG_FRAMECOUNT             = TIFFTAG(34232) /* Sequence Frame Count */
	TIFFTAG_MODELTRANSFORMATIONTAG = TIFFTAG(34264) /* Used in interchangeable GeoTIFF files */
	/* tag 34377 is private tag registered to Adobe for PhotoShop */
	TIFFTAG_PHOTOSHOP = TIFFTAG(34377)
	/* tags 34665, 34853 and 40965 are documented in EXIF specification */
	TIFFTAG_EXIFIFD = TIFFTAG(34665) /* Pointer to EXIF private directory */
	/* tag 34750 is a private tag registered to Adobe? */
	TIFFTAG_ICCPROFILE = TIFFTAG(34675) /* ICC profile data */
	TIFFTAG_IMAGELAYER = TIFFTAG(34732) /* !TIFF/FX image layer information */
	/* tag 34750 is a private tag registered to Pixel Magic */
	TIFFTAG_JBIGOPTIONS = TIFFTAG(34750) /* JBIG options */
	TIFFTAG_GPSIFD      = TIFFTAG(34853) /* Pointer to EXIF GPS private directory */
	/* tags 34908-34914 are private tags registered to SGI */
	TIFFTAG_FAXRECVPARAMS = TIFFTAG(34908) /* encoded Class 2 ses. params */
	TIFFTAG_FAXSUBADDRESS = TIFFTAG(34909) /* received SubAddr string */
	TIFFTAG_FAXRECVTIME   = TIFFTAG(34910) /* receive time (secs) */
	TIFFTAG_FAXDCS        = TIFFTAG(34911) /* encoded fax ses. params, Table 2/T.30 */
	/* tags 37439-37443 are registered to SGI <gregl@sgi.com> */
	TIFFTAG_STONITS = TIFFTAG(37439) /* Sample value to Nits */
	/* tag 34929 is a private tag registered to FedEx */
	TIFFTAG_FEDEX_EDR                      = TIFFTAG(34929) /* unknown use */
	TIFFTAG_IMAGESOURCEDATA                = TIFFTAG(37724) /* http://justsolve.archiveteam.org/wiki/PSD, http://www.adobe.com/devnet-apps/photoshop/fileformatashtml/ */
	TIFFTAG_INTEROPERABILITYIFD            = TIFFTAG(40965) /* Pointer to EXIF Interoperability private directory */
	TIFFTAG_GDAL_METADATA                  = TIFFTAG(42112) /* Used by the GDAL library */
	TIFFTAG_GDAL_NODATA                    = TIFFTAG(42113) /* Used by the GDAL library */
	TIFFTAG_OCE_SCANJOB_DESCRIPTION        = TIFFTAG(50215) /* Used in the Oce scanning process */
	TIFFTAG_OCE_APPLICATION_SELECTOR       = TIFFTAG(50216) /* Used in the Oce scanning process. */
	TIFFTAG_OCE_IDENTIFICATION_NUMBER      = TIFFTAG(50217)
	TIFFTAG_OCE_IMAGELOGIC_CHARACTERISTICS = TIFFTAG(50218)
	/* tags 50674 to 50677 are reserved for ESRI */
	TIFFTAG_LERC_PARAMETERS = TIFFTAG(50674) /* Stores LERC version and additional compression method */

	/* Adobe Digital Negative (DNG) format tags */
	TIFFTAG_DNGVERSION           = TIFFTAG(50706) /* &DNG version number */
	TIFFTAG_DNGBACKWARDVERSION   = TIFFTAG(50707) /* &DNG compatibility version */
	TIFFTAG_UNIQUECAMERAMODEL    = TIFFTAG(50708) /* &name for the camera model */
	TIFFTAG_LOCALIZEDCAMERAMODEL = TIFFTAG(50709) /* &localized camera model name (UTF-8) */
	TIFFTAG_CFAPLANECOLOR        = TIFFTAG(50710) /* &CFAPattern->LinearRaw space mapping */
	TIFFTAG_CFALAYOUT            = TIFFTAG(50711) /* &spatial layout of the CFA */
	TIFFTAG_LINEARIZATIONTABLE   = TIFFTAG(50712) /* &lookup table description */
	TIFFTAG_BLACKLEVELREPEATDIM  = TIFFTAG(50713) /* &repeat pattern size for the BlackLevel tag */
	TIFFTAG_BLACKLEVEL           = TIFFTAG(50714) /* &zero light encoding level */
	TIFFTAG_BLACKLEVELDELTAH     = TIFFTAG(50715) /* &zero light encoding level differences (columns) */
	TIFFTAG_BLACKLEVELDELTAV     = TIFFTAG(50716) /* &zero light encoding level differences (rows) */
	TIFFTAG_WHITELEVEL           = TIFFTAG(50717) /* &fully saturated encoding level */
	TIFFTAG_DEFAULTSCALE         = TIFFTAG(50718) /* &default scale factors */
	TIFFTAG_DEFAULTCROPORIGIN    = TIFFTAG(50719) /* &origin of the final image area */
	TIFFTAG_DEFAULTCROPSIZE      = TIFFTAG(50720) /* &size of the final image area */
	TIFFTAG_COLORMATRIX1         = TIFFTAG(50721) /* &XYZ->reference color space transformation matrix 1 */
	TIFFTAG_COLORMATRIX2         = TIFFTAG(50722) /* &XYZ->reference color space transformation matrix 2 */
	TIFFTAG_CAMERACALIBRATION1   = TIFFTAG(50723) /* &calibration matrix 1 */
	TIFFTAG_CAMERACALIBRATION2   = TIFFTAG(50724) /* &calibration matrix 2 */
	TIFFTAG_REDUCTIONMATRIX1     = TIFFTAG(50725) /* &dimensionality reduction matrix 1 */
	TIFFTAG_REDUCTIONMATRIX2     = TIFFTAG(50726) /* &dimensionality reduction matrix 2 */
	TIFFTAG_ANALOGBALANCE        = TIFFTAG(50727) /* &gain applied the stored raw values*/
	TIFFTAG_ASSHOTNEUTRAL        = TIFFTAG(50728) /* &selected white balance in linear reference space */
	TIFFTAG_ASSHOTWHITEXY        = TIFFTAG(50729) /* &selected white balance in x-y chromaticity coordinates */
	TIFFTAG_BASELINEEXPOSURE     = TIFFTAG(50730) /* &how much to move the zero point */
	TIFFTAG_BASELINENOISE        = TIFFTAG(50731) /* &relative noise level */
	TIFFTAG_BASELINESHARPNESS    = TIFFTAG(50732) /* &relative amount of sharpening */
	/* TIFFTAG_BAYERGREENSPLIT: &how closely the values of the green pixels in the blue/green rows
	 * track the values of the green pixels in the red/green rows */
	TIFFTAG_BAYERGREENSPLIT         = TIFFTAG(50733)
	TIFFTAG_LINEARRESPONSELIMIT     = TIFFTAG(50734) /* &non-linear encoding range */
	TIFFTAG_CAMERASERIALNUMBER      = TIFFTAG(50735) /* &camera's serial number */
	TIFFTAG_LENSINFO                = TIFFTAG(50736) /* info about the lens */
	TIFFTAG_CHROMABLURRADIUS        = TIFFTAG(50737) /* &chroma blur radius */
	TIFFTAG_ANTIALIASSTRENGTH       = TIFFTAG(50738) /* &relative strength of the camera's anti-alias filter */
	TIFFTAG_SHADOWSCALE             = TIFFTAG(50739) /* &used by Adobe Camera Raw */
	TIFFTAG_DNGPRIVATEDATA          = TIFFTAG(50740) /* &manufacturer's private data */
	TIFFTAG_MAKERNOTESAFETY         = TIFFTAG(50741) /* &whether the EXIF MakerNote tag is safe to preserve along with the rest of the EXIF data */
	TIFFTAG_CALIBRATIONILLUMINANT1  = TIFFTAG(50778) /* &illuminant 1 */
	TIFFTAG_CALIBRATIONILLUMINANT2  = TIFFTAG(50779) /* &illuminant 2 */
	TIFFTAG_BESTQUALITYSCALE        = TIFFTAG(50780) /* &best quality multiplier */
	TIFFTAG_RAWDATAUNIQUEID         = TIFFTAG(50781) /* &unique identifier for the raw image data */
	TIFFTAG_ORIGINALRAWFILENAME     = TIFFTAG(50827) /* &file name of the original raw file (UTF-8) */
	TIFFTAG_ORIGINALRAWFILEDATA     = TIFFTAG(50828) /* &contents of the original raw file */
	TIFFTAG_ACTIVEAREA              = TIFFTAG(50829) /* &active (non-masked) pixels of the sensor */
	TIFFTAG_MASKEDAREAS             = TIFFTAG(50830) /* &list of coordinates of fully masked pixels */
	TIFFTAG_ASSHOTICCPROFILE        = TIFFTAG(50831) /* &these two tags used to */
	TIFFTAG_ASSHOTPREPROFILEMATRIX  = TIFFTAG(50832) /* map cameras's color space  into ICC profile space */
	TIFFTAG_CURRENTICCPROFILE       = TIFFTAG(50833) /* & */
	TIFFTAG_CURRENTPREPROFILEMATRIX = TIFFTAG(50834) /* & */

	/* DNG 1.2.0.0 */
	TIFFTAG_COLORIMETRICREFERENCE       = TIFFTAG(50879) /* &colorimetric reference */
	TIFFTAG_CAMERACALIBRATIONSIGNATURE  = TIFFTAG(50931) /* &camera calibration signature (UTF-8) */
	TIFFTAG_PROFILECALIBRATIONSIGNATURE = TIFFTAG(50932) /* &profile calibration signature (UTF-8) */
	/* TIFFTAG_EXTRACAMERAPROFILES 50933 &extra camera profiles : is already defined for GeoTIFF DGIWG */
	TIFFTAG_ASSHOTPROFILENAME         = TIFFTAG(50934) /* &as shot profile name (UTF-8) */
	TIFFTAG_NOISEREDUCTIONAPPLIED     = TIFFTAG(50935) /* &amount of applied noise reduction */
	TIFFTAG_PROFILENAME               = TIFFTAG(50936) /* &camera profile name (UTF-8) */
	TIFFTAG_PROFILEHUESATMAPDIMS      = TIFFTAG(50937) /* &dimensions of HSV mapping */
	TIFFTAG_PROFILEHUESATMAPDATA1     = TIFFTAG(50938) /* &first HSV mapping table */
	TIFFTAG_PROFILEHUESATMAPDATA2     = TIFFTAG(50939) /* &second HSV mapping table */
	TIFFTAG_PROFILETONECURVE          = TIFFTAG(50940) /* &default tone curve */
	TIFFTAG_PROFILEEMBEDPOLICY        = TIFFTAG(50941) /* &profile embedding policy */
	TIFFTAG_PROFILECOPYRIGHT          = TIFFTAG(50942) /* &profile copyright information (UTF-8) */
	TIFFTAG_FORWARDMATRIX1            = TIFFTAG(50964) /* &matrix for mapping white balanced camera colors to XYZ D50 */
	TIFFTAG_FORWARDMATRIX2            = TIFFTAG(50965) /* &matrix for mapping white balanced camera colors to XYZ D50 */
	TIFFTAG_PREVIEWAPPLICATIONNAME    = TIFFTAG(50966) /* &name of application that created preview (UTF-8) */
	TIFFTAG_PREVIEWAPPLICATIONVERSION = TIFFTAG(50967) /* &version of application that created preview (UTF-8) */
	TIFFTAG_PREVIEWSETTINGSNAME       = TIFFTAG(50968) /* &name of conversion settings (UTF-8) */
	TIFFTAG_PREVIEWSETTINGSDIGEST     = TIFFTAG(50969) /* &unique id of conversion settings */
	TIFFTAG_PREVIEWCOLORSPACE         = TIFFTAG(50970) /* &preview color space */
	TIFFTAG_PREVIEWDATETIME           = TIFFTAG(50971) /* &date/time preview was rendered */
	TIFFTAG_RAWIMAGEDIGEST            = TIFFTAG(50972) /* &md5 of raw image data */
	TIFFTAG_ORIGINALRAWFILEDIGEST     = TIFFTAG(50973) /* &md5 of the data stored in the OriginalRawFileData tag */
	TIFFTAG_SUBTILEBLOCKSIZE          = TIFFTAG(50974) /* &subtile block size */
	TIFFTAG_ROWINTERLEAVEFACTOR       = TIFFTAG(50975) /* &number of interleaved fields */
	TIFFTAG_PROFILELOOKTABLEDIMS      = TIFFTAG(50981) /* &num of input samples in each dim of default "look" table */
	TIFFTAG_PROFILELOOKTABLEDATA      = TIFFTAG(50982) /* &default "look" table for use as starting point */

	/* DNG 1.3.0.0 */
	TIFFTAG_OPCODELIST1  = TIFFTAG(51008) /* &opcodes that should be applied to raw image after reading */
	TIFFTAG_OPCODELIST2  = TIFFTAG(51009) /* &opcodes that should be applied after mapping to linear reference */
	TIFFTAG_OPCODELIST3  = TIFFTAG(51022) /* &opcodes that should be applied after demosaicing */
	TIFFTAG_NOISEPROFILE = TIFFTAG(51041) /* &noise profile */

	/* DNG 1.4.0.0 */
	TIFFTAG_DEFAULTUSERCROP              = TIFFTAG(51125) /* &default user crop rectangle in relative coords */
	TIFFTAG_DEFAULTBLACKRENDER           = TIFFTAG(51110) /* &black rendering hint */
	TIFFTAG_BASELINEEXPOSUREOFFSET       = TIFFTAG(51109) /* &baseline exposure offset */
	TIFFTAG_PROFILELOOKTABLEENCODING     = TIFFTAG(51108) /* &3D LookTable indexing conversion */
	TIFFTAG_PROFILEHUESATMAPENCODING     = TIFFTAG(51107) /* &3D HueSatMap indexing conversion */
	TIFFTAG_ORIGINALDEFAULTFINALSIZE     = TIFFTAG(51089) /* &default final size of larger original file for this proxy */
	TIFFTAG_ORIGINALBESTQUALITYFINALSIZE = TIFFTAG(51090) /* &best quality final size of larger original file for this proxy */
	TIFFTAG_ORIGINALDEFAULTCROPSIZE      = TIFFTAG(51091) /* &the default crop size of larger original file for this proxy */
	TIFFTAG_NEWRAWIMAGEDIGEST            = TIFFTAG(51111) /* &modified MD5 digest of the raw image data */
	TIFFTAG_RAWTOPREVIEWGAIN             = TIFFTAG(51112) /* &The gain between the main raw FD and the preview IFD containing this tag */

	/* DNG 1.5.0.0 */
	TIFFTAG_DEPTHFORMAT      = TIFFTAG(51177) /* &encoding of the depth data in the file */
	TIFFTAG_DEPTHNEAR        = TIFFTAG(51178) /* &distance from the camera represented by value 0 in the depth map */
	TIFFTAG_DEPTHFAR         = TIFFTAG(51179) /* &distance from the camera represented by the maximum value in the depth map */
	TIFFTAG_DEPTHUNITS       = TIFFTAG(51180) /* &measurement units for DepthNear and DepthFar */
	TIFFTAG_DEPTHMEASURETYPE = TIFFTAG(51181) /* &measurement geometry for the depth map */
	TIFFTAG_ENHANCEPARAMS    = TIFFTAG(51182) /* &a string that documents how the enhanced image data was processed. */

	/* DNG 1.6.0.0 */
	TIFFTAG_PROFILEGAINTABLEMAP    = TIFFTAG(52525) /* &spatially varying gain tables that can be applied as starting point */
	TIFFTAG_SEMANTICNAME           = TIFFTAG(52526) /* &a string that identifies the semantic mask */
	TIFFTAG_SEMANTICINSTANCEID     = TIFFTAG(52528) /* &a string that identifies a specific instance in a semantic mask */
	TIFFTAG_MASKSUBAREA            = TIFFTAG(52536) /* &the crop rectangle of this IFD's mask, relative to the main image */
	TIFFTAG_RGBTABLES              = TIFFTAG(52543) /* &color transforms to apply to masked image regions */
	TIFFTAG_CALIBRATIONILLUMINANT3 = TIFFTAG(52529) /* &the illuminant used for the third set of color calibration tags */
	TIFFTAG_COLORMATRIX3           = TIFFTAG(52531) /* &matrix to convert XYZ values to reference camera native color space under CalibrationIlluminant3 */
	TIFFTAG_CAMERACALIBRATION3     = TIFFTAG(52530) /* &matrix to transform reference camera native space values to individual camera native space values under CalibrationIlluminant3 */
	TIFFTAG_REDUCTIONMATRIX3       = TIFFTAG(52538) /* &dimensionality reduction matrix for use in color conversion to XYZ under CalibrationIlluminant3 */
	TIFFTAG_PROFILEHUESATMAPDATA3  = TIFFTAG(52537) /* &the data for the third HSV table */
	TIFFTAG_FORWARDMATRIX3         = TIFFTAG(52532) /* &matrix to map white balanced camera colors to XYZ D50 */
	TIFFTAG_ILLUMINANTDATA1        = TIFFTAG(52533) /* &data for the first calibration illuminant */
	TIFFTAG_ILLUMINANTDATA2        = TIFFTAG(52534) /* &data for the second calibration illuminant */
	TIFFTAG_ILLUMINANTDATA3        = TIFFTAG(53535) /* &data for the third calibration illuminant */

	/* TIFF/EP */
	TIFFTAG_EP_CFAREPEATPATTERNDIM = TIFFTAG(33421) /* dimensions of CFA pattern */
	TIFFTAG_EP_CFAPATTERN          = TIFFTAG(33422) /* color filter array pattern */
	TIFFTAG_EP_BATTERYLEVEL        = TIFFTAG(33423) /* battery level (rational or ASCII) */
	TIFFTAG_EP_INTERLACE           = TIFFTAG(34857) /* Number of multi-field images */
	/* TIFFTAG_EP_IPTC_NAA and TIFFTAG_RICHTIFFIPTC share the same tag number (33723)
	 *   LibTIFF type is UNDEFINED or BYTE, but often times incorrectly specified as LONG,
	 *   because TIFF/EP (ISO/DIS 12234-2) specifies type LONG or ASCII. */
	TIFFTAG_EP_IPTC_NAA                 = TIFFTAG(33723) /* Alias IPTC/NAA Newspaper Association RichTIFF */
	TIFFTAG_EP_TIMEZONEOFFSET           = TIFFTAG(34858) /* Time zone offset relative to UTC */
	TIFFTAG_EP_SELFTIMERMODE            = TIFFTAG(34859) /* Number of seconds capture was delayed from button press */
	TIFFTAG_EP_FLASHENERGY              = TIFFTAG(37387) /* Flash energy, or range if there is uncertainty */
	TIFFTAG_EP_SPATIALFREQUENCYRESPONSE = TIFFTAG(37388) /* Spatial frequency response */
	TIFFTAG_EP_NOISE                    = TIFFTAG(37389) /* Camera noise measurement values */
	TIFFTAG_EP_FOCALPLANEXRESOLUTION    = TIFFTAG(37390) /* Focal plane X resolution */
	TIFFTAG_EP_FOCALPLANEYRESOLUTION    = TIFFTAG(37391) /* Focal plane Y resolution */
	TIFFTAG_EP_FOCALPLANERESOLUTIONUNIT = TIFFTAG(37392) /* Focal plane resolution unit */
	TIFFTAG_EP_IMAGENUMBER              = TIFFTAG(37393) /* Number of image when several of burst shot stored in same TIFF/EP */
	TIFFTAG_EP_SECURITYCLASSIFICATION   = TIFFTAG(37394) /* Security classification */
	TIFFTAG_EP_IMAGEHISTORY             = TIFFTAG(37395) /* Record of what has been done to the image */
	TIFFTAG_EP_EXPOSUREINDEX            = TIFFTAG(37397) /* Exposure index */
	TIFFTAG_EP_STANDARDID               = TIFFTAG(37398) /* TIFF/EP standard version, n.n.n.n */
	TIFFTAG_EP_SENSINGMETHOD            = TIFFTAG(37399) /* Type of image sensor */
	/*
	 * TIFF/EP tags equivalent to EXIF tags
	 *     Note that TIFF-EP and EXIF use nearly the same metadata tag set, but TIFF-EP stores the tags in IFD 0,
	 *     while EXIF store the tags in a separate IFD. Either location is allowed by DNG, but the EXIF location is preferred.
	 */
	TIFFTAG_EP_EXPOSURETIME           = TIFFTAG(33434) /* Exposure time */
	TIFFTAG_EP_FNUMBER                = TIFFTAG(33437) /* F number */
	TIFFTAG_EP_EXPOSUREPROGRAM        = TIFFTAG(34850) /* Exposure program */
	TIFFTAG_EP_SPECTRALSENSITIVITY    = TIFFTAG(34852) /* Spectral sensitivity */
	TIFFTAG_EP_ISOSPEEDRATINGS        = TIFFTAG(34855) /* ISO speed rating */
	TIFFTAG_EP_OECF                   = TIFFTAG(34856) /* Optoelectric conversion factor */
	TIFFTAG_EP_DATETIMEORIGINAL       = TIFFTAG(36867) /* Date and time of original data generation */
	TIFFTAG_EP_COMPRESSEDBITSPERPIXEL = TIFFTAG(37122) /* Image compression mode */
	TIFFTAG_EP_SHUTTERSPEEDVALUE      = TIFFTAG(37377) /* Shutter speed */
	TIFFTAG_EP_APERTUREVALUE          = TIFFTAG(37378) /* Aperture */
	TIFFTAG_EP_BRIGHTNESSVALUE        = TIFFTAG(37379) /* Brightness */
	TIFFTAG_EP_EXPOSUREBIASVALUE      = TIFFTAG(37380) /* Exposure bias */
	TIFFTAG_EP_MAXAPERTUREVALUE       = TIFFTAG(37381) /* Maximum lens aperture */
	TIFFTAG_EP_SUBJECTDISTANCE        = TIFFTAG(37382) /* Subject distance */
	TIFFTAG_EP_METERINGMODE           = TIFFTAG(37383) /* Metering mode */
	TIFFTAG_EP_LIGHTSOURCE            = TIFFTAG(37384) /* Light source */
	TIFFTAG_EP_FLASH                  = TIFFTAG(37385) /* Flash */
	TIFFTAG_EP_FOCALLENGTH            = TIFFTAG(37386) /* Lens focal length */
	TIFFTAG_EP_SUBJECTLOCATION        = TIFFTAG(37396) /* Subject location (area) */

	TIFFTAG_RPCCOEFFICIENT       = TIFFTAG(50844) /* Define by GDAL for geospatial georeferencing through RPC: http://geotiff.maptools.org/rpc_prop.html */
	TIFFTAG_ALIAS_LAYER_METADATA = TIFFTAG(50784) /* Alias Sketchbook Pro layer usage description. */

	/* GeoTIFF DGIWG */
	TIFFTAG_TIFF_RSID           = TIFFTAG(50908) /* https://www.awaresystems.be/imaging/tiff/tifftags/tiff_rsid.html */
	TIFFTAG_GEO_METADATA        = TIFFTAG(50909) /* https://www.awaresystems.be/imaging/tiff/tifftags/geo_metadata.html */
	TIFFTAG_EXTRACAMERAPROFILES = TIFFTAG(50933) /* http://wwwimages.adobe.com/www.adobe.com/content/dam/Adobe/en/products/photoshop/pdfs/dng_spec_1.4.0.0.pdf */

	/* tag 65535 is an undefined tag used by Eastman Kodak */
	TIFFTAG_DCSHUESHIFTVALUES = TIFFTAG(65535) /* hue shift correction data */

	/*
	 * The following are ``pseudo tags'' that can be used to control
	 * codec-specific functionality.  These tags are not written to file.
	 * Note that these values start at 0xffff+1 so that they'll never
	 * collide with Aldus-assigned tags.
	 *
	 * If you want your private pseudo tags ``registered'' (i.e. added to
	 * this file), please post a bug report via the tracking system at
	 * http://www.remotesensing.org/libtiff/bugs.html with the appropriate
	 * C definitions to add.
	 */
	TIFFTAG_FAXMODE     = TIFFTAG(65536)  /* Group 3/4 format control */
	FAXMODE_CLASSIC     = TIFFTAG(0x0000) /* default, include RTC */
	FAXMODE_NORTC       = TIFFTAG(0x0001) /* no RTC at end of data */
	FAXMODE_NOEOL       = TIFFTAG(0x0002) /* no EOL code at end of row */
	FAXMODE_BYTEALIGN   = TIFFTAG(0x0004) /* byte align row */
	FAXMODE_WORDALIGN   = TIFFTAG(0x0008) /* word align row */
	FAXMODE_CLASSF      = FAXMODE_NORTC   /* TIFF Class F */
	TIFFTAG_JPEGQUALITY = TIFFTAG(65537)  /* Compression quality level */
	/* Note: quality level is on the IJG 0-100 scale.  Default value is 75 */
	TIFFTAG_JPEGCOLORMODE  = TIFFTAG(65538)  /* Auto RGB<=>YCbCr convert? */
	JPEGCOLORMODE_RAW      = TIFFTAG(0x0000) /* no conversion (default) */
	JPEGCOLORMODE_RGB      = TIFFTAG(0x0001) /* do auto conversion */
	TIFFTAG_JPEGTABLESMODE = TIFFTAG(65539)  /* What to put in JPEGTables */
	JPEGTABLESMODE_QUANT   = TIFFTAG(0x0001) /* include quantization tbls */
	JPEGTABLESMODE_HUFF    = TIFFTAG(0x0002) /* include Huffman tbls */
	/* Note: default is JPEGTABLESMODE_QUANT | JPEGTABLESMODE_HUFF */
	TIFFTAG_FAXFILLFUNC        = TIFFTAG(65540) /* G3/G4 fill function */
	TIFFTAG_PIXARLOGDATAFMT    = TIFFTAG(65549) /* PixarLogCodec I/O data sz */
	PIXARLOGDATAFMT_8BIT       = TIFFTAG(0)     /* regular u_char samples */
	PIXARLOGDATAFMT_8BITABGR   = TIFFTAG(1)     /* ABGR-order u_chars */
	PIXARLOGDATAFMT_11BITLOG   = TIFFTAG(2)     /* 11-bit log-encoded (raw) */
	PIXARLOGDATAFMT_12BITPICIO = TIFFTAG(3)     /* as per PICIO (1.0==2048) */
	PIXARLOGDATAFMT_16BIT      = TIFFTAG(4)     /* signed short samples */
	PIXARLOGDATAFMT_FLOAT      = TIFFTAG(5)     /* IEEE float samples */
	/* 65550-65556 are allocated to Oceana Matrix <dev@oceana.com> */
	TIFFTAG_DCSIMAGERTYPE     = TIFFTAG(65550) /* imager model & filter */
	DCSIMAGERMODEL_M3         = TIFFTAG(0)     /* M3 chip (1280 x 1024) */
	DCSIMAGERMODEL_M5         = TIFFTAG(1)     /* M5 chip (1536 x 1024) */
	DCSIMAGERMODEL_M6         = TIFFTAG(2)     /* M6 chip (3072 x 2048) */
	DCSIMAGERFILTER_IR        = TIFFTAG(0)     /* infrared filter */
	DCSIMAGERFILTER_MONO      = TIFFTAG(1)     /* monochrome filter */
	DCSIMAGERFILTER_CFA       = TIFFTAG(2)     /* color filter array */
	DCSIMAGERFILTER_OTHER     = TIFFTAG(3)     /* other filter */
	TIFFTAG_DCSINTERPMODE     = TIFFTAG(65551) /* interpolation mode */
	DCSINTERPMODE_NORMAL      = TIFFTAG(0x0)   /* whole image, default */
	DCSINTERPMODE_PREVIEW     = TIFFTAG(0x1)   /* preview of image (384x256) */
	TIFFTAG_DCSBALANCEARRAY   = TIFFTAG(65552) /* color balance values */
	TIFFTAG_DCSCORRECTMATRIX  = TIFFTAG(65553) /* color correction values */
	TIFFTAG_DCSGAMMA          = TIFFTAG(65554) /* gamma value */
	TIFFTAG_DCSTOESHOULDERPTS = TIFFTAG(65555) /* toe & shoulder points */
	TIFFTAG_DCSCALIBRATIONFD  = TIFFTAG(65556) /* calibration file desc */
	/* Note: quality level is on the ZLIB 1-9 scale. Default value is -1 */
	TIFFTAG_ZIPQUALITY      = TIFFTAG(65557) /* compression quality level */
	TIFFTAG_PIXARLOGQUALITY = TIFFTAG(65558) /* PixarLog uses same scale */
	/* 65559 is allocated to Oceana Matrix <dev@oceana.com> */
	TIFFTAG_DCSCLIPRECTANGLE     = TIFFTAG(65559) /* area of image to acquire */
	TIFFTAG_SGILOGDATAFMT        = TIFFTAG(65560) /* SGILog user data format */
	SGILOGDATAFMT_FLOAT          = TIFFTAG(0)     /* IEEE float samples */
	SGILOGDATAFMT_16BIT          = TIFFTAG(1)     /* 16-bit samples */
	SGILOGDATAFMT_RAW            = TIFFTAG(2)     /* uninterpreted data */
	SGILOGDATAFMT_8BIT           = TIFFTAG(3)     /* 8-bit RGB monitor values */
	TIFFTAG_SGILOGENCODE         = TIFFTAG(65561) /* SGILog data encoding control*/
	SGILOGENCODE_NODITHER        = TIFFTAG(0)     /* do not dither encoded values*/
	SGILOGENCODE_RANDITHER       = TIFFTAG(1)     /* randomly dither encd values */
	TIFFTAG_LZMAPRESET           = TIFFTAG(65562) /* LZMA2 preset (compression level) */
	TIFFTAG_PERSAMPLE            = TIFFTAG(65563) /* interface for per sample tags */
	PERSAMPLE_MERGED             = TIFFTAG(0)     /* present as a single value */
	PERSAMPLE_MULTI              = TIFFTAG(1)     /* present as multiple values */
	TIFFTAG_ZSTD_LEVEL           = TIFFTAG(65564) /* ZSTD compression level */
	TIFFTAG_LERC_VERSION         = TIFFTAG(65565) /* LERC version */
	LERC_VERSION_2_4             = TIFFTAG(4)
	TIFFTAG_LERC_ADD_COMPRESSION = TIFFTAG(65566) /* LERC additional compression */
	LERC_ADD_COMPRESSION_NONE    = TIFFTAG(0)
	LERC_ADD_COMPRESSION_DEFLATE = TIFFTAG(1)
	LERC_ADD_COMPRESSION_ZSTD    = TIFFTAG(2)
	TIFFTAG_LERC_MAXZERROR       = TIFFTAG(65567) /* LERC maximum error */
	TIFFTAG_WEBP_LEVEL           = TIFFTAG(65568) /* WebP compression level */
	TIFFTAG_WEBP_LOSSLESS        = TIFFTAG(65569) /* WebP lossless/lossy */
	TIFFTAG_WEBP_LOSSLESS_EXACT  = TIFFTAG(65571) /* WebP lossless exact mode. Set-only mode. Default is 1. Can be set to 0 to increase compression rate, but R,G,B in areas where alpha = 0 will not be preserved */
	TIFFTAG_DEFLATE_SUBCODEC     = TIFFTAG(65570) /* ZIP codec: to get/set the sub-codec to use. Will default to libdeflate when available */
	DEFLATE_SUBCODEC_ZLIB        = TIFFTAG(0)
	DEFLATE_SUBCODEC_LIBDEFLATE  = TIFFTAG(1)

	/*
	 * EXIF tags
	 */
	EXIFTAG_EXPOSURETIME        = TIFFTAG(33434) /* Exposure time */
	EXIFTAG_FNUMBER             = TIFFTAG(33437) /* F number */
	EXIFTAG_EXPOSUREPROGRAM     = TIFFTAG(34850) /* Exposure program */
	EXIFTAG_SPECTRALSENSITIVITY = TIFFTAG(34852) /* Spectral sensitivity */
	/* After EXIF 2.2.1 ISOSpeedRatings is named PhotographicSensitivity.
	   In addition, while "Count=Any", only 1 count should be used. */
	EXIFTAG_ISOSPEEDRATINGS          = TIFFTAG(34855) /* ISO speed rating */
	EXIFTAG_PHOTOGRAPHICSENSITIVITY  = TIFFTAG(34855) /* Photographic Sensitivity (new name for tag 34855) */
	EXIFTAG_OECF                     = TIFFTAG(34856) /* Optoelectric conversion factor */
	EXIFTAG_EXIFVERSION              = TIFFTAG(36864) /* Exif version */
	EXIFTAG_DATETIMEORIGINAL         = TIFFTAG(36867) /* Date and time of original data generation */
	EXIFTAG_DATETIMEDIGITIZED        = TIFFTAG(36868) /* Date and time of digital data generation */
	EXIFTAG_COMPONENTSCONFIGURATION  = TIFFTAG(37121) /* Meaning of each component */
	EXIFTAG_COMPRESSEDBITSPERPIXEL   = TIFFTAG(37122) /* Image compression mode */
	EXIFTAG_SHUTTERSPEEDVALUE        = TIFFTAG(37377) /* Shutter speed */
	EXIFTAG_APERTUREVALUE            = TIFFTAG(37378) /* Aperture */
	EXIFTAG_BRIGHTNESSVALUE          = TIFFTAG(37379) /* Brightness */
	EXIFTAG_EXPOSUREBIASVALUE        = TIFFTAG(37380) /* Exposure bias */
	EXIFTAG_MAXAPERTUREVALUE         = TIFFTAG(37381) /* Maximum lens aperture */
	EXIFTAG_SUBJECTDISTANCE          = TIFFTAG(37382) /* Subject distance */
	EXIFTAG_METERINGMODE             = TIFFTAG(37383) /* Metering mode */
	EXIFTAG_LIGHTSOURCE              = TIFFTAG(37384) /* Light source */
	EXIFTAG_FLASH                    = TIFFTAG(37385) /* Flash */
	EXIFTAG_FOCALLENGTH              = TIFFTAG(37386) /* Lens focal length */
	EXIFTAG_SUBJECTAREA              = TIFFTAG(37396) /* Subject area */
	EXIFTAG_MAKERNOTE                = TIFFTAG(37500) /* Manufacturer notes */
	EXIFTAG_USERCOMMENT              = TIFFTAG(37510) /* User comments */
	EXIFTAG_SUBSECTIME               = TIFFTAG(37520) /* DateTime subseconds */
	EXIFTAG_SUBSECTIMEORIGINAL       = TIFFTAG(37521) /* DateTimeOriginal subseconds */
	EXIFTAG_SUBSECTIMEDIGITIZED      = TIFFTAG(37522) /* DateTimeDigitized subseconds */
	EXIFTAG_FLASHPIXVERSION          = TIFFTAG(40960) /* Supported Flashpix version */
	EXIFTAG_COLORSPACE               = TIFFTAG(40961) /* Color space information */
	EXIFTAG_PIXELXDIMENSION          = TIFFTAG(40962) /* Valid image width */
	EXIFTAG_PIXELYDIMENSION          = TIFFTAG(40963) /* Valid image height */
	EXIFTAG_RELATEDSOUNDFILE         = TIFFTAG(40964) /* Related audio file */
	EXIFTAG_FLASHENERGY              = TIFFTAG(41483) /* Flash energy */
	EXIFTAG_SPATIALFREQUENCYRESPONSE = TIFFTAG(41484) /* Spatial frequency response */
	EXIFTAG_FOCALPLANEXRESOLUTION    = TIFFTAG(41486) /* Focal plane X resolution */
	EXIFTAG_FOCALPLANEYRESOLUTION    = TIFFTAG(41487) /* Focal plane Y resolution */
	EXIFTAG_FOCALPLANERESOLUTIONUNIT = TIFFTAG(41488) /* Focal plane resolution unit */
	EXIFTAG_SUBJECTLOCATION          = TIFFTAG(41492) /* Subject location */
	EXIFTAG_EXPOSUREINDEX            = TIFFTAG(41493) /* Exposure index */
	EXIFTAG_SENSINGMETHOD            = TIFFTAG(41495) /* Sensing method */
	EXIFTAG_FILESOURCE               = TIFFTAG(41728) /* File source */
	EXIFTAG_SCENETYPE                = TIFFTAG(41729) /* Scene type */
	EXIFTAG_CFAPATTERN               = TIFFTAG(41730) /* CFA pattern */
	EXIFTAG_CUSTOMRENDERED           = TIFFTAG(41985) /* Custom image processing */
	EXIFTAG_EXPOSUREMODE             = TIFFTAG(41986) /* Exposure mode */
	EXIFTAG_WHITEBALANCE             = TIFFTAG(41987) /* White balance */
	EXIFTAG_DIGITALZOOMRATIO         = TIFFTAG(41988) /* Digital zoom ratio */
	EXIFTAG_FOCALLENGTHIN35MMFILM    = TIFFTAG(41989) /* Focal length in 35 mm film */
	EXIFTAG_SCENECAPTURETYPE         = TIFFTAG(41990) /* Scene capture type */
	EXIFTAG_GAINCONTROL              = TIFFTAG(41991) /* Gain control */
	EXIFTAG_CONTRAST                 = TIFFTAG(41992) /* Contrast */
	EXIFTAG_SATURATION               = TIFFTAG(41993) /* Saturation */
	EXIFTAG_SHARPNESS                = TIFFTAG(41994) /* Sharpness */
	EXIFTAG_DEVICESETTINGDESCRIPTION = TIFFTAG(41995) /* Device settings description */
	EXIFTAG_SUBJECTDISTANCERANGE     = TIFFTAG(41996) /* Subject distance range */
	EXIFTAG_IMAGEUNIQUEID            = TIFFTAG(42016) /* Unique image ID */

	/*--: New for EXIF-Version 2.32, May 2019 ... */
	EXIFTAG_SENSITIVITYTYPE           = TIFFTAG(34864) /* The SensitivityType tag indicates which one of the parameters of ISO12232 is the PhotographicSensitivity tag. */
	EXIFTAG_STANDARDOUTPUTSENSITIVITY = TIFFTAG(34865) /* This tag indicates the standard output sensitivity value of a camera or input device defined in ISO 12232. */
	EXIFTAG_RECOMMENDEDEXPOSUREINDEX  = TIFFTAG(34866) /* recommended exposure index   */
	EXIFTAG_ISOSPEED                  = TIFFTAG(34867) /* ISO speed value */
	EXIFTAG_ISOSPEEDLATITUDEYYY       = TIFFTAG(34868) /* ISO speed latitude yyy */
	EXIFTAG_ISOSPEEDLATITUDEZZZ       = TIFFTAG(34869) /* ISO speed latitude zzz */
	EXIFTAG_OFFSETTIME                = TIFFTAG(36880) /* offset from UTC of the time of DateTime tag. */
	EXIFTAG_OFFSETTIMEORIGINAL        = TIFFTAG(36881) /* offset from UTC of the time of DateTimeOriginal tag. */
	EXIFTAG_OFFSETTIMEDIGITIZED       = TIFFTAG(36882) /* offset from UTC of the time of DateTimeDigitized tag. */
	EXIFTAG_TEMPERATURE               = TIFFTAG(37888) /* Temperature as the ambient situation at the shot in dergee Celsius */
	EXIFTAG_HUMIDITY                  = TIFFTAG(37889) /* Humidity as the ambient situation at the shot in percent */
	EXIFTAG_PRESSURE                  = TIFFTAG(37890) /* Pressure as the ambient situation at the shot hecto-Pascal (hPa) */
	EXIFTAG_WATERDEPTH                = TIFFTAG(37891) /* WaterDepth as the ambient situation at the shot in meter (m) */
	EXIFTAG_ACCELERATION              = TIFFTAG(37892) /* Acceleration (a scalar regardless of direction) as the ambientsituation at the shot in units of mGal (10-5 m/s^2) */
	/* EXIFTAG_CAMERAELEVATIONANGLE: Elevation/depression. angle of the orientation of the  camera(imaging optical axis)
	 *                               as the ambient situation at the shot in degree from -180deg to +180deg. */
	EXIFTAG_CAMERAELEVATIONANGLE = TIFFTAG(37893)
	EXIFTAG_CAMERAOWNERNAME      = TIFFTAG(42032) /* owner of a camera */
	EXIFTAG_BODYSERIALNUMBER     = TIFFTAG(42033) /* serial number of the body of the camera */
	/* EXIFTAG_LENSSPECIFICATION: minimum focal length (in mm), maximum focal length (in mm),minimum F number in the minimum focal length,
	 *                            and minimum F number in the maximum focal length, */
	EXIFTAG_LENSSPECIFICATION                   = TIFFTAG(42034)
	EXIFTAG_LENSMAKE                            = TIFFTAG(42035) /* the lens manufacturer */
	EXIFTAG_LENSMODEL                           = TIFFTAG(42036) /* the lens model name and model number */
	EXIFTAG_LENSSERIALNUMBER                    = TIFFTAG(42037) /* the serial number of the interchangeable lens */
	EXIFTAG_GAMMA                               = TIFFTAG(42240) /* value of coefficient gamma */
	EXIFTAG_COMPOSITEIMAGE                      = TIFFTAG(42080) /* composite image */
	EXIFTAG_SOURCEIMAGENUMBEROFCOMPOSITEIMAGE   = TIFFTAG(42081) /* source image number of composite image */
	EXIFTAG_SOURCEEXPOSURETIMESOFCOMPOSITEIMAGE = TIFFTAG(42082) /* source exposure times of composite image */

	/*
	 * EXIF-GPS tags  (Version 2.31, July 2016)
	 */
	GPSTAG_VERSIONID            = TIFFTAG(0)  /* Indicates the version of GPSInfoIFD. */
	GPSTAG_LATITUDEREF          = TIFFTAG(1)  /* Indicates whether the latitude is north or south latitude. */
	GPSTAG_LATITUDE             = TIFFTAG(2)  /* Indicates the latitude. */
	GPSTAG_LONGITUDEREF         = TIFFTAG(3)  /* Indicates whether the longitude is east or west longitude. */
	GPSTAG_LONGITUDE            = TIFFTAG(4)  /* Indicates the longitude. */
	GPSTAG_ALTITUDEREF          = TIFFTAG(5)  /* Indicates the altitude used as the reference altitude. */
	GPSTAG_ALTITUDE             = TIFFTAG(6)  /* Indicates the altitude based on the reference in GPSAltitudeRef. */
	GPSTAG_TIMESTAMP            = TIFFTAG(7)  /*Indicates the time as UTC (Coordinated Universal Time). */
	GPSTAG_SATELLITES           = TIFFTAG(8)  /*Indicates the GPS satellites used for measurements. */
	GPSTAG_STATUS               = TIFFTAG(9)  /* Indicates the status of the GPS receiver when the image is  recorded. */
	GPSTAG_MEASUREMODE          = TIFFTAG(10) /* Indicates the GPS measurement mode. */
	GPSTAG_DOP                  = TIFFTAG(11) /* Indicates the GPS DOP (data degree of precision). */
	GPSTAG_SPEEDREF             = TIFFTAG(12) /* Indicates the unit used to express the GPS receiver speed of movement. */
	GPSTAG_SPEED                = TIFFTAG(13) /* Indicates the speed of GPS receiver movement. */
	GPSTAG_TRACKREF             = TIFFTAG(14) /* Indicates the reference for giving the direction of GPS receiver movement. */
	GPSTAG_TRACK                = TIFFTAG(15) /* Indicates the direction of GPS receiver movement. */
	GPSTAG_IMGDIRECTIONREF      = TIFFTAG(16) /* Indicates the reference for giving the direction of the image when it is captured. */
	GPSTAG_IMGDIRECTION         = TIFFTAG(17) /* Indicates the direction of the image when it was captured. */
	GPSTAG_MAPDATUM             = TIFFTAG(18) /* Indicates the geodetic survey data used by the GPS receiver. (e.g. WGS-84) */
	GPSTAG_DESTLATITUDEREF      = TIFFTAG(19) /* Indicates whether the latitude of the destination point is north or south latitude. */
	GPSTAG_DESTLATITUDE         = TIFFTAG(20) /* Indicates the latitude of the destination point. */
	GPSTAG_DESTLONGITUDEREF     = TIFFTAG(21) /* Indicates whether the longitude of the destination point is east or west longitude. */
	GPSTAG_DESTLONGITUDE        = TIFFTAG(22) /* Indicates the longitude of the destination point. */
	GPSTAG_DESTBEARINGREF       = TIFFTAG(23) /* Indicates the reference used for giving the bearing to the destination point. */
	GPSTAG_DESTBEARING          = TIFFTAG(24) /* Indicates the bearing to the destination point. */
	GPSTAG_DESTDISTANCEREF      = TIFFTAG(25) /* Indicates the unit used to express the distance to the destination point. */
	GPSTAG_DESTDISTANCE         = TIFFTAG(26) /* Indicates the distance to the destination point. */
	GPSTAG_PROCESSINGMETHOD     = TIFFTAG(27) /* A character string recording the name of the method used for location finding. */
	GPSTAG_AREAINFORMATION      = TIFFTAG(28) /* A character string recording the name of the GPS area. */
	GPSTAG_DATESTAMP            = TIFFTAG(29) /* A character string recording date and time information relative to UTC (Coordinated Universal Time). */
	GPSTAG_DIFFERENTIAL         = TIFFTAG(30) /* Indicates whether differential correction is applied to the GPS receiver. */
	GPSTAG_GPSHPOSITIONINGERROR = TIFFTAG(31) /* Indicates horizontal positioning errors in meters. */
)
