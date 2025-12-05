# go-libtiff

[![Go Reference](https://pkg.go.dev/badge/github.com/klippa-app/go-libtiff/go-libtiff.svg)](https://pkg.go.dev/github.com/klippa-app/go-libtiff)

:rocket: *Easy to use TIFF library using Go, libtiff and Wazero* :rocket:

**A fast, multi-threaded and easy to use TIFF library for Go applications.**

## Features

* WebAssembly build of libtiff, so no need for local dependencies
* Contains libtiff and the binary tools that come with it (like tiff2pdf)
* Ability to run the binary tools through the library or through CLI
* This library will handle all complicated cgo/WebAssembly gymnastics for you, no direct WebAssembly usage/knowledge
  required
* Ability to open files from a Go file reader
* Helper method render tiff file to Go image and binary image (JPEG or PNG)
* libjpeg-turbo implementation to speed up JEPG compresssion (CGO + native library required)

Note: this library does not provide direct access to all libtiff methods yet,
please open an issue if you want to see a method added.

### Todo

* Implement more libtiff methods (please open issues if you need something)

## libtiff

This project uses the libtiff library (https://libtiff.gitlab.io/libtiff/) to process the TIFF files . 
Therefor this project could also be called a binding.

Please be aware that libtiff comes with its [own license](https://gitlab.com/libtiff/libtiff/-/blob/master/LICENSE.md).

## WebAssembly

This project uses WebAssembly, we use the [Wazero runtime](https://wazero.io/) for running WebAssembly within Go. 
The comes with quite some advantages:

- WebAssembly is one binary for every platform, which means we can embed the WebAssembly version of libtiff in this
  repository
- Because we have `go:embed`, we can embed the WebAssembly binary inside the Go binary
- Because of this, you won't have to download and distribute libtiff yourself, making deployments simpler
- Wazero is pure Go, and thus runs on any platform where you can run Go on (a lot), which also allows you to run
  go-libtiff and libtiff on all of those platforms once
- Because it's running in Go, you can directly access the memory in Wazero, for example to render a tiff image
  directly to a Go Image
- Since libtiff is compiled to WebAssembly and runs inside the Wazero runtime, it basically runs in a sandbox:
    - No chance of crashing the Go process like with cgo
    - No access to other local resources in case of attacks on libtiff (disk, network, memory)
    - Full control over file access (you decide which folders Wazero exposes to libtiff, by default it exposes the whole
      disk)

Please be aware that Wazero comes with the `Apache License 2.0` license.

## File handling

Because you can tell Wazero which folders have to be mounted in WebAssembly, you have full control over the filesystem.

By default, go-libtiff will mount the full root disk in Wazero on non-Windows environments.
On Windows environments, go-libtiff will get the volume of the current working directory and mount that as the root.

You can change this behaviour by overwriting FSConfig in the `libtiff.GetInstance()` call.

All paths given to go-libtiff in WebAssembly mode have to be in POSIX style and have to be absolute, so for
example: `/home/user/Downloads/file.pdf`. If you have mounted `/home/user/`on the root, then the path you would have to
give is `/Downloads/file.pdf`, this is the same on Windows, so no backward slashes or volume names in paths.

You can set your own mounts by overwriting FSConfig in the `libtiff.GetInstance()` call.

## Getting started

The following example shows how to open a tiff file and render it to a Go image:

```go
ctx := context.Background()
instance, err := libtiff.GetInstance(ctx, &libtiff.Config{})
if err != nil {
	log.Fatal(err)
}
defer instance.Close(ctx)

file, err := instance.TIFFOpenFileFromPath(ctx, input, nil)
if err != nil {
	log.Fatal(err)
}
defer file.Close(ctx)

// Loop over the images in the TIFF directory.
for i := range file.Directories(ctx) {
    renderedImage, cleanup, err := file.ToGoImage(ctx)
    if err != nil {
        log.Fatal(fmt.Errorf("could not convert tiff image %d to go image: %w", i, err))
    }
    defer cleanup(ctx)
}
```

### Instance re-use

Since libtiff allows you to open multiple files at the same time and operate on them, you can re-use the instance
for multiple files. This is more efficient than opening a new instance for every file.

## libtiff tools

You can use any of the libtiff tools by importing a tool package, for example:

```go
import "github.com/klippa-app/go-libtiff/tiff2pdf"
```

You can then run the tool using:

```go
ctx := context.Background()
err := tiff2pdf.Run(ctx, []string{"input.tiff", "output.pdf"})
if err != nil {
	log.Fatal(err)
}
```

It's also possible to override the instance config:

```go
ctx := context.Background()
ctx = libtiff.ConfigInContext(ctx, &libtiff.Config{
	CompilationCache: compilationCache,
})
err := tiff2pdf.Run(ctx, []string{"input.tiff", "output.pdf"})
if err != nil {
log.Fatal(err)
}
```

## CLI tool

You can install the CLI tool using:

```bash
go install github.com/klippa-app/go-libtiff@latest
go-libtiff tiff2pdf input.tiff output.pdf
```

This will provide you access to the following tools:

- fax2ps
- fax2tiff
- pal2rgb
- ppm2tiff
- raw2tiff
- rgb2ycbcr
- thumbnail
- tiff2bw
- tiff2pdf
- tiff2ps
- tiff2rgba
- tiffcmp
- tiffcp
- tiffcrop
- tiffdither
- tiffdump
- tiffinfo
- tiffmedian
- tiffset
- tiffsplit
- tiff2img (tool of this project to render tiff to images (JPEG and PNG))

Please be aware that these tools mount your own filesystem inside the Wazero runtime to give the tools access to the
files, since they can't access the files from Go itself, the only difference is in the tiff2img tool.

Note: the tool mkg3states is not available because it won't run in Wazero, the tool also seems to be more of a build
tool than a runtime tool.

## Compilation cache

When using the CLI tool, you can set the `LIBTIFF_COMPILATION_CACHE_DIR` environment variable to enable compilation cache.
The value will be the directory that the compilation cache will be written to. For libtiff usage you can provide a Wazero
compilation cache instance to config in the `libtiff.GetInstance()` call. For example:

```go
ctx := context.Background()
ctx = libtiff.ConfigInContext(ctx, &libtiff.Config{
	CompilationCache: wazero.NewCompilationCache(), // Memory cache
})
instance, err := libtiff.GetInstance(ctx)
if err != nil {
	log.Fatal(err)
}
```

Or for a file cache:

```go
ctx := context.Background()
compilationCache, err = wazero.NewCompilationCacheWithDir(".libtiff-compilation-cache")
if err != nil {
	log.Fatal(err)
}
ctx = libtiff.ConfigInContext(ctx, &libtiff.Config{
	CompilationCache: compilationCache,
})
instance, err := libtiff.GetInstance(ctx)
if err != nil {
	log.Fatal(err)
}
```


## Improving JPEG rendering speed

By default, this library renders images with the `image/jpeg` package that comes with Go to make distribution as simple
as possible. However, this package is quite slow compared to other native libraries like libjpeg and libjpeg-turbo, you
can enable the usage of libjpeg-turbo by using the build tag `libtiff_use_turbojpeg`, this will require you to have the
package `libturbojpeg-dev` installed during build time and the `libturbojpeg` package during runtime and build time.

Speed improvements that can be expected are significant, for example: on a simple PDF the full process of rendering a
page is 3x as fast compared to a build without libjpeg-turbo.

Please note we use CGO for libjpeg-turbo for now. There are plans to compile libjpeg-turbo to WebAssembly in the future.

## Support Policy

We offer an API stability promise with semantic versioning. In other words, we promise to not break any exported
function signature without incrementing the major version. New features and behaviors happen with a minor version 
increment, e.g. 1.0.11 to 1.1.0. We also fix bugs or change internal details with a patch version, e.g. 1.0.0 to 1.0.1.

### Go

This project will support the last 2 version of Go, this means that if the last version of Go is 1.24, our `go.mod`
will be set to Go 1.23, and our CI tests will be run on Go 1.23 and 1.24. This is in line with Go's
[Release Policy](https://go.dev/doc/devel/release).

It won't mean that the library won't work with older versions of Go, but it will tell you what to expect of the
supported  Go versions. If we change the supported Go versions, we will make that a minor version upgrade. This policy
allows you to not be forced to the latest Go version for a pretty long time, but it still allows us to use new language
features in a pretty reasonable time-frame.

## About Klippa

Founded in 2015, [Klippa](https://www.klippa.com/en)'s goal is to digitize & automate administrative processes with
modern technologies. We help clients enhance the effectiveness of their organization by using machine learning and OCR.
Since 2015, more than a thousand happy clients have used Klippa's software solutions. Klippa currently has an
international team of 50 people, with offices in Groningen, Amsterdam and Brasov.

## License

The MIT License (MIT)
