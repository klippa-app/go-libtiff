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
  -DCMAKE_CXX_FLAGS="-O2" \
  -Dlibdeflate=ON -DDeflate_INCLUDE_DIR=/build/libdeflate -DDeflate_LIBRARY=/build/libdeflate/libdeflate.a \
  -Djbig=ON -DJBIG_INCLUDE_DIR=/build/jbigkit-2.1/libjbig -DJBIG_LIBRARY=/build/jbigkit-2.1/libjbig/libjbig85.a

# Build the tools and shared library
emmake make

# Build the WASM file for libtiff.
emcc -O2 -s ALLOW_MEMORY_GROWTH=1 -s ALLOW_TABLE_GROWTH=1 -s STANDALONE_WASM=1 -s ERROR_ON_UNDEFINED_SYMBOLS=0 -s EXPORTED_FUNCTIONS="_TIFFGetField,_TIFFClose,_TIFFReadDirectory,_TIFFSetDirectory,_TIFFReadRGBAImageOriented,_TIFFOpen,_TIFFGetFieldUint32_t,_TIFFGetFieldFloat,_free,_malloc,_calloc,_realloc" -s EXPORTED_RUNTIME_METHODS="ccall,cwrap,addFunction,removeFunction" -s LLD_REPORT_UNDEFINED -s WASM=1 -o "build/libtiff.html" -I/build/tiff-${LIBTIFF_PKGVER}/libtiff libtiff/libtiff.a ../emsdk/upstream/emscripten/cache/sysroot/lib/wasm32-emscripten/libjpeg.a ../emsdk/upstream/emscripten/cache/sysroot/lib/wasm32-emscripten/libz.a ../extra.c --no-entry

# Todo: copy WASM files to the right location.