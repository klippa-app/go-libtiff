#include <emscripten.h>
#include <stdlib.h>
#include <tiffio.h>

extern tmsize_t TIFFReadWriteProcGoCB(thandle_t, void*, tmsize_t);
extern toff_t TIFFSeekProcGoCB(thandle_t, toff_t, int);
extern int TIFFCloseProcGoCB(thandle_t);
extern toff_t TIFFSizeProcGoCB(thandle_t);
extern int TIFFMapFileProcGoCB(thandle_t, void **base, toff_t *size);
extern void TIFFUnmapFileProcGoCB(thandle_t, void *base, toff_t size);

extern int TIFFGetFieldUint32_t(TIFF *tif, uint32_t tag, uint32_t *val);
extern int TIFFGetFieldFloat(TIFF *tif, uint32_t tag, float *val);
extern TIFF* TIFFClientOpenGo(const char *filename, const char *mode, thandle_t clientdata);

EMSCRIPTEN_KEEPALIVE
int TIFFGetFieldUint32_t(TIFF *tif, uint32_t tag, uint32_t *val) {
  return TIFFGetField(tif, tag, val);
}

EMSCRIPTEN_KEEPALIVE
int TIFFGetFieldFloat(TIFF *tif, uint32_t tag, float *val) {
  return TIFFGetField(tif, tag, val);
}

EMSCRIPTEN_KEEPALIVE
TIFF* TIFFClientOpenGo(const char *filename, const char *mode, thandle_t clientdata) {
  return TIFFClientOpen(filename, mode, clientdata, TIFFReadWriteProcGoCB, TIFFReadWriteProcGoCB, TIFFSeekProcGoCB, TIFFCloseProcGoCB, TIFFSizeProcGoCB, TIFFMapFileProcGoCB, TIFFUnmapFileProcGoCB);
}
