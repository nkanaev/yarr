package main

import (
	"os"
	"path"
	"io/ioutil"
)

var plist = `<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
	<key>CFBundleName</key>
	<string>yarr</string>
	<key>CFBundleDisplayName</key>
	<string>yarr</string>
	<key>CFBundleIdentifier</key>
	<string>nkanaev.yarr</string>
	<key>CFBundleVersion</key>
	<string>1.0</string>
	<key>CFBundlePackageType</key>
	<string>APPL</string>
	<key>CFBundleExecutable</key>
	<string>yarr</string>

	<key>CFBundleIconFile</key>
	<string>AppIcon</string>
	<key>CFBundleIconName</key>
	<string>AppIcon</string>
	<key>LSApplicationCategoryType</key>
	<string>public.app-category.news</string>

	<key>NSHighResolutionCapable</key>
	<string>True</string>

	<key>CFBundleInfoDictionaryVersion</key>
	<string>6.0</string>
	<key>CFBundleShortVersionString</key>
	<string>1.1</string>

	<key>LSMinimumSystemVersion</key>
	<string>10.13</string>
	<key>LSUIElement</key>
	<true/>
	<key>NSHumanReadableCopyright</key>
	<string>Copyright © 2020 nkanaev. All rights reserved.</string>
</dict>
</plist>
`

func main() {
	outdir := os.Args[1]
	outfile := "yarr"

	binDir := path.Join(outdir, "yarr.app", "Contents/MacOS")
	resDir := path.Join(outdir, "yarr.app", "Contents/Resources")
	plistPath := path.Join(outdir, "yarr.app", "Contents/Info.plist")

	os.MkdirAll(binDir, 0700)
	os.MkdirAll(resDir, 0700)

	f, _ := ioutil.ReadFile(path.Join(outdir, outfile))
	ioutil.WriteFile(path.Join(binDir, outfile), f, 0700)

	ioutil.WriteFile(plistPath, []byte(plist), 0600)
}
