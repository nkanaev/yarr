package main

import (
	"flag"
	"io/ioutil"
	"strings"
)

var rsrc = `1 VERSIONINFO
FILEVERSION     {VERSION_COMMA},0,0
PRODUCTVERSION  {VERSION_COMMA},0,0
BEGIN
  BLOCK "StringFileInfo"
  BEGIN
    BLOCK "080904E4"
    BEGIN
      VALUE "CompanyName", "Old MacDonald's Farm"
      VALUE "FileDescription", "Yet another RSS reader"
      VALUE "FileVersion", "{VERSION}"
      VALUE "InternalName", "yarr"
      VALUE "LegalCopyright", "nkanaev"
      VALUE "OriginalFilename", "yarr.exe"
      VALUE "ProductName", "yarr"
      VALUE "ProductVersion", "{VERSION}"
    END
  END
  BLOCK "VarFileInfo"
  BEGIN
    VALUE "Translation", 0x809, 1252
  END
END

1 ICON "icon.ico"
`

func main() {
	var version, outfile string
	flag.StringVar(&version, "version", "0.0", "")
	flag.StringVar(&outfile, "outfile", "versioninfo.rc", "")
	flag.Parse()

	version_comma := strings.ReplaceAll(version, ".", ",")

	out := strings.ReplaceAll(rsrc, "{VERSION}", version)
	out = strings.ReplaceAll(out, "{VERSION_COMMA}", version_comma)

	ioutil.WriteFile(outfile, []byte(out), 0644)
}
