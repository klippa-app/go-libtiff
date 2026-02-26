# Execute in a clean ubuntu:24.04 docker.
apt install git software-properties-common curl wget build-essential

git clone https://github.com/emscripten-core/emsdk.git
cd emsdk
git pull
git checkout 4.0.20
./emsdk install 4.0.20
./emsdk activate 4.0.20
source ./emsdk_env.sh
cd upstream/emscripten
rm -Rf cache
patch -p1 < ../../../emscripten.patch
cd ../../../

embuilder build libjpeg
embuilder build zlib

LIBTIFF_PKGVER=4.7.1
curl -o "tiff-${LIBTIFF_PKGVER}.tar.gz" "http://download.osgeo.org/libtiff/tiff-${LIBTIFF_PKGVER}.tar.gz"
tar xzvf tiff-${LIBTIFF_PKGVER}.tar.gz
rm tiff-${LIBTIFF_PKGVER}.tar.gz
cd tiff-${LIBTIFF_PKGVER}
emcmake cmake \
  -DCMAKE_EXE_LINKER_FLAGS="-sERROR_ON_UNDEFINED_SYMBOLS=0 -sWASM=1 -sALLOW_MEMORY_GROWTH=1 -sSTANDALONE_WASM=1" \
  -DCMAKE_CXX_FLAGS="-O2"

# Build the tools and shared library
emmake make

# Build the WASM file for libtiff.
emcc -O2 -s ALLOW_MEMORY_GROWTH=1 -s ALLOW_TABLE_GROWTH=1 -s STANDALONE_WASM=1 -s ERROR_ON_UNDEFINED_SYMBOLS=0 -s EXPORTED_FUNCTIONS="_TIFFGetField,_TIFFClose,_TIFFReadDirectory,_TIFFSetDirectory,_TIFFReadRGBAImageOriented,_TIFFCurrentDirectory,_TIFFLastDirectory,_TIFFNumberOfDirectories,_TIFFOpen,_TIFFGetFieldUint16_t,_TIFFGetFieldUint32_t,_TIFFGetFieldInt,_TIFFGetFieldFloat,_TIFFGetFieldDouble,_TIFFGetFieldConstChar,_TIFFOpenOptionsAlloc,_TIFFOpenOptionsFree,_TIFFOpenOptionsSetMaxSingleMemAlloc,_TIFFOpenOptionsSetMaxCumulatedMemAlloc,_TIFFOpenOptionsSetErrorHandlerExtRGo,_TIFFOpenOptionsSetWarningHandlerExtRGo,_TIFFOpenOptionsSetWarnAboutUnknownTags,_TIFFOpenExt,_TIFFClientOpenExt,_TIFFClientOpenExtGo,_TIFFGetVersion,_TIFFOpenExtGo,_TIFFSetFieldUint16_t,_TIFFSetFieldUint32_t,_TIFFSetFieldInt,_TIFFSetFieldFloat,_TIFFSetFieldDouble,_TIFFSetFieldString,_TIFFSetFieldExtraSamples,_TIFFWriteEncodedStrip,_TIFFWriteDirectory,_TIFFDefaultStripSize,_free,_malloc,_calloc,_realloc,_vsprintf" -s EXPORTED_RUNTIME_METHODS="ccall,cwrap,addFunction,removeFunction" -s LLD_REPORT_UNDEFINED -s WASM=1 -o "build/libtiff.html" -I/build/tiff-4.7.1/libtiff libtiff/libtiff.a ../emsdk/upstream/emscripten/cache/sysroot/lib/wasm32-emscripten/libjpeg.a ../emsdk/upstream/emscripten/cache/sysroot/lib/wasm32-emscripten/libz.a ../extra.c --no-entry

# Copy files to the right locations.
cd ../../
find build/tiff-4.7.1/tools build/tiff-4.7.1/build -name '*.wasm' -exec sh -c 'echo "Copying $(basename {})"; cp {} $(basename {} .wasm)/$(basename {})'  \;