#include <emscripten.h>
#include <stdlib.h>
#include <string.h>
#include <stdio.h>
#include <tiffio.h>

extern tmsize_t TIFFReadProcGoCB(thandle_t, void*, tmsize_t);
extern tmsize_t TIFFWriteProcGoCB(thandle_t, void*, tmsize_t);
extern toff_t TIFFSeekProcGoCB(thandle_t, toff_t, int);
extern int TIFFCloseProcGoCB(thandle_t);
extern toff_t TIFFSizeProcGoCB(thandle_t);
extern int TIFFMapFileProcGoCB(thandle_t, void **base, toff_t *size);
extern void TIFFUnmapFileProcGoCB(thandle_t, void *base, toff_t size);
extern int TIFFOpenOptionsSetErrorHandlerExtRGoCB(TIFF *tif, void *user_data, const char *module, const char *fmt, va_list ap);
extern int TIFFOpenOptionsSetWarningHandlerExtRGoCB(TIFF *tif, void *user_data, const char *module, const char *fmt, va_list ap);

EMSCRIPTEN_KEEPALIVE
int TIFFGetFieldUint16_t(TIFF *tif, uint32_t tag, uint16_t *val) {
  return TIFFGetField(tif, tag, val);
}

EMSCRIPTEN_KEEPALIVE
int TIFFGetFieldUint32_t(TIFF *tif, uint32_t tag, uint32_t *val) {
  return TIFFGetField(tif, tag, val);
}

EMSCRIPTEN_KEEPALIVE
int TIFFGetFieldInt(TIFF *tif, uint32_t tag, int *val) {
  return TIFFGetField(tif, tag, val);
}

EMSCRIPTEN_KEEPALIVE
int TIFFGetFieldFloat(TIFF *tif, uint32_t tag, float *val) {
  return TIFFGetField(tif, tag, val);
}

EMSCRIPTEN_KEEPALIVE
int TIFFGetFieldDouble(TIFF *tif, uint32_t tag, double *val) {
  return TIFFGetField(tif, tag, val);
}

EMSCRIPTEN_KEEPALIVE
int TIFFGetFieldConstChar(TIFF *tif, uint32_t tag, const char** val) {
  return TIFFGetField(tif, tag, val);
}

EMSCRIPTEN_KEEPALIVE
TIFF* TIFFOpenExtGo(const char *filename, const char *mode, TIFFOpenOptions *opts) {
  return TIFFOpenExt(filename, mode, opts);
}

EMSCRIPTEN_KEEPALIVE
TIFF* TIFFClientOpenExtGo(const char *filename, const char *mode, thandle_t clientdata, TIFFOpenOptions *opts) {
  return TIFFClientOpenExt(filename, mode, clientdata, TIFFReadProcGoCB, TIFFWriteProcGoCB, TIFFSeekProcGoCB, TIFFCloseProcGoCB, TIFFSizeProcGoCB, TIFFMapFileProcGoCB, TIFFUnmapFileProcGoCB, opts);
}

EMSCRIPTEN_KEEPALIVE
void TIFFOpenOptionsSetErrorHandlerExtRGo(TIFFOpenOptions *opts, void *errorhandler_user_data) {
  TIFFOpenOptionsSetErrorHandlerExtR(opts, TIFFOpenOptionsSetErrorHandlerExtRGoCB, errorhandler_user_data);
}

EMSCRIPTEN_KEEPALIVE
void TIFFOpenOptionsSetWarningHandlerExtRGo(TIFFOpenOptions *opts, void *warnhandler_user_data) {
  TIFFOpenOptionsSetWarningHandlerExtR(opts, TIFFOpenOptionsSetWarningHandlerExtRGoCB, warnhandler_user_data);
}

EMSCRIPTEN_KEEPALIVE
int TIFFSetFieldUint16_t(TIFF *tif, uint32_t tag, uint16_t val) {
  return TIFFSetField(tif, tag, val);
}

EMSCRIPTEN_KEEPALIVE
int TIFFSetFieldUint32_t(TIFF *tif, uint32_t tag, uint32_t val) {
  return TIFFSetField(tif, tag, val);
}

EMSCRIPTEN_KEEPALIVE
int TIFFSetFieldInt(TIFF *tif, uint32_t tag, int val) {
  return TIFFSetField(tif, tag, val);
}

EMSCRIPTEN_KEEPALIVE
int TIFFSetFieldFloat(TIFF *tif, uint32_t tag, float val) {
  return TIFFSetField(tif, tag, val);
}

EMSCRIPTEN_KEEPALIVE
int TIFFSetFieldDouble(TIFF *tif, uint32_t tag, double val) {
  return TIFFSetField(tif, tag, val);
}

EMSCRIPTEN_KEEPALIVE
int TIFFSetFieldString(TIFF *tif, uint32_t tag, const char* val) {
  return TIFFSetField(tif, tag, val);
}

EMSCRIPTEN_KEEPALIVE
int TIFFSetFieldExtraSamples(TIFF *tif, uint16_t count, uint16_t *types) {
  return TIFFSetField(tif, TIFFTAG_EXTRASAMPLES, count, types);
}

EMSCRIPTEN_KEEPALIVE
int TIFFSetFieldTwoUint16(TIFF *tif, uint32_t tag, uint16_t val1, uint16_t val2) {
  return TIFFSetField(tif, tag, val1, val2);
}

EMSCRIPTEN_KEEPALIVE
int TIFFGetFieldTwoUint16(TIFF *tif, uint32_t tag, uint16_t *val1, uint16_t *val2) {
  return TIFFGetField(tif, tag, val1, val2);
}

EMSCRIPTEN_KEEPALIVE
int TIFFGetFieldUint64_t(TIFF *tif, uint32_t tag, uint64_t *val) {
  return TIFFGetField(tif, tag, val);
}

EMSCRIPTEN_KEEPALIVE
int TIFFSetFieldUint64_t(TIFF *tif, uint32_t tag, uint64_t val) {
  return TIFFSetField(tif, tag, val);
}

// TIFFGetFieldDefaulted wrappers - like TIFFGetField but returns TIFF spec
// defaults for tags that are not explicitly set.
EMSCRIPTEN_KEEPALIVE
int TIFFGetFieldDefaultedUint16_t(TIFF *tif, uint32_t tag, uint16_t *val) {
  return TIFFGetFieldDefaulted(tif, tag, val);
}

EMSCRIPTEN_KEEPALIVE
int TIFFGetFieldDefaultedUint32_t(TIFF *tif, uint32_t tag, uint32_t *val) {
  return TIFFGetFieldDefaulted(tif, tag, val);
}

EMSCRIPTEN_KEEPALIVE
int TIFFGetFieldDefaultedUint64_t(TIFF *tif, uint32_t tag, uint64_t *val) {
  return TIFFGetFieldDefaulted(tif, tag, val);
}

EMSCRIPTEN_KEEPALIVE
int TIFFGetFieldDefaultedInt(TIFF *tif, uint32_t tag, int *val) {
  return TIFFGetFieldDefaulted(tif, tag, val);
}

EMSCRIPTEN_KEEPALIVE
int TIFFGetFieldDefaultedFloat(TIFF *tif, uint32_t tag, float *val) {
  return TIFFGetFieldDefaulted(tif, tag, val);
}

EMSCRIPTEN_KEEPALIVE
int TIFFGetFieldDefaultedDouble(TIFF *tif, uint32_t tag, double *val) {
  return TIFFGetFieldDefaulted(tif, tag, val);
}

EMSCRIPTEN_KEEPALIVE
int TIFFGetFieldDefaultedConstChar(TIFF *tif, uint32_t tag, const char** val) {
  return TIFFGetFieldDefaulted(tif, tag, val);
}

EMSCRIPTEN_KEEPALIVE
int TIFFGetFieldDefaultedTwoUint16(TIFF *tif, uint32_t tag, uint16_t *val1, uint16_t *val2) {
  return TIFFGetFieldDefaulted(tif, tag, val1, val2);
}

// TIFFPrintDirectoryToBuffer writes the directory info to a memory buffer
// instead of a FILE*. Returns the number of bytes written, or -1 on error.
EMSCRIPTEN_KEEPALIVE
int TIFFPrintDirectoryToBuffer(TIFF *tif, char *buf, int bufsize, long flags) {
  char *membuf = NULL;
  size_t membufsize = 0;
  FILE *f = open_memstream(&membuf, &membufsize);
  if (!f) return -1;
  TIFFPrintDirectory(tif, f, flags);
  fclose(f);
  int len = (int)(membufsize < (size_t)(bufsize - 1) ? membufsize : (size_t)(bufsize - 1));
  memcpy(buf, membuf, len);
  buf[len] = '\0';
  free(membuf);
  return len;
}
