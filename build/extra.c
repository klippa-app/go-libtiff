#include <emscripten.h>
#include <stdlib.h>
#include <tiffio.h>

extern int TIFFGetFieldUint32_t(TIFF *tif, uint32_t tag, uint32_t *val);
extern int TIFFGetFieldFloat(TIFF *tif, uint32_t tag, float *val);

EMSCRIPTEN_KEEPALIVE
int TIFFGetFieldUint32_t(TIFF *tif, uint32_t tag, uint32_t *val) {
  return TIFFGetField(tif, tag, val);
}

EMSCRIPTEN_KEEPALIVE
int TIFFGetFieldFloat(TIFF *tif, uint32_t tag, float *val) {
  return TIFFGetField(tif, tag, val);
}
