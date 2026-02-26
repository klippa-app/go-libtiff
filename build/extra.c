#include <emscripten.h>
#include <stdlib.h>
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
