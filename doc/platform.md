Incomplete & inaccurate platform-specific notes.

# MacOS Icon

The format for desktop apps is [.icns][icns].
AFAIK, the format is not open (even though it had been [reverse-engineered][icns-re]),
and I couldn't find any 3rd party tool that'd fully support it.

The easiest way for creating icon file is either via `Xcode`,
or by using built-in `iconutil` command that ships with MacOS.

The steps are provided below:

    $ sips -s format png --resampleWidth 1024 source.png --out /path/to/icons/icon_512x512@2x.png
    $ sips -s format png --resampleWidth  512 source.png --out /path/to/icons/icon_512x512.png
    $ sips -s format png --resampleWidth  256 source.png --out /path/to/icons/icon_256x256.png
    $ sips -s format png --resampleWidth  128 source.png --out /path/to/icons/icon_128x128.png
    $ sips -s format png --resampleWidth   64 source.png --out /path/to/icons/icon_32x32@2x.png
    $ sips -s format png --resampleWidth   32 source.png --out /path/to/icons/icon_32x32.png
    $ sips -s format png --resampleWidth   16 source.png --out /path/to/icons/icon_16x16.png
    $ iconutil -c icns /path/to/icons -o icon.icns

[icns]: https://en.wikipedia.org/wiki/Apple_Icon_Image_format
[icns-re]: https://www.macdisk.com/maciconen.php#RLE

# Windows Icon

Terminology:

- coff: precursor to pe format (portable executable). pe is an extension of coff.
- manifest: xml file with platform requirements needed during runtime
  - https://docs.microsoft.com/en-us/windows/win32/sbscs/application-manifests
  - https://www.samlogic.net/articles/manifest.htm
- rc: dsl file that describes the application metadata & resources
  - https://docs.microsoft.com/en-gb/windows/win32/menurc/about-resource-files
  - https://github.com/josephspurrier/goversioninfo/blob/master/testdata/rc/versioninfo.rc (sample rc)

Windows Icons are directly embedded to the binary.
To do so one needs to provide `.syso` file prior to compiling Go code,
which will be passed to the linker. So, basically `.syso` is any
[object file][obj-file] that the linker understands.

More info here: [ticket][syso-ticket] & [commit][syso-commit].

Note to self: running `go build main.go` [won't embed][syso-quirk]
.syso file if it isn't located in a package directory.

Tools to create `.syso` files:

- [windres][windres]: ships with mingw (gnu tools for windows)
- [rsrc][rsrc]: written in Go, wasn't considered at the time
  due to the critical bug with icon alignment
- [goversioninfo][goversioninfo]: rsrc wrapper
  with manifest file creation via json

[obj-file]: https://en.wikipedia.org/wiki/Object_file
[syso-linker]: https://github.com/golang/go/issues/23278#issuecomment-354567634
[syso-ticket]: https://github.com/golang/go/issues/1552
[syso-commit]: https://github.com/golang/go/commit/b0996334
[syso-quirk]: https://github.com/golang/go/issues/16090
[mingw]: https://en.wikipedia.org/wiki/MinGW
[coff]: https://en.wikipedia.org/wiki/COFF
[windres]: https://sourceware.org/binutils/docs/binutils/windres.html
[rsrs]: https://github.com/akavel/rsrc
[rsrc-bug]: https://github.com/akavel/rsrc/issues/12
[goversioninfo]: github.com/josephspurrier/goversioninfo

[winicon-guide]: https://docs.microsoft.com/en-us/windows/win32/uxguide/vis-icons#size-requirements
[res-vs-coff]: http://www.mingw.org/wiki/MS_resource_compiler
[versioninfo-resource]: https://docs.microsoft.com/en-us/windows/win32/menurc/versioninfo-resource
