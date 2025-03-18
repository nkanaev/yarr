#/bin/sh

set -e

usage() {
  echo "usage: $0 VERSION path/to/icon.icns path/to/binary output/dir"
}

if [ $# -eq 0 ]; then
    usage
    exit
fi

VERSION=$1
ICNFILE=$2
BINFILE=$3
OUTPATH=$4

mkdir -p $OUTPATH/yarr.app/Contents/MacOS
mkdir -p $OUTPATH/yarr.app/Contents/Resources

mv $BINFILE $OUTPATH/yarr.app/Contents/MacOS/yarr
cp $ICNFILE $OUTPATH/yarr.app/Contents/Resources/icon.icns

chmod u+x  $OUTPATH/yarr.app/Contents/MacOS/yarr

echo -n 'APPL????' >$OUTPATH/yarr.app/Contents/PkgInfo
cat <<EOF >$OUTPATH/yarr.app/Contents/Info.plist
<?xml version="1.0" encoding="UTF-8"?>
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
	<string>$VERSION</string>
	<key>CFBundlePackageType</key>
	<string>APPL</string>
	<key>CFBundleExecutable</key>
	<string>yarr</string>

	<key>CFBundleIconFile</key>
	<string>icon</string>
	<key>LSApplicationCategoryType</key>
	<string>public.app-category.news</string>

	<key>NSHighResolutionCapable</key>
	<string>True</string>

	<key>LSMinimumSystemVersion</key>
	<string>10.13</string>
	<key>LSUIElement</key>
	<true/>
	<key>NSHumanReadableCopyright</key>
	<string>Copyright © 2020 nkanaev. All rights reserved.</string>
</dict>
</plist>
EOF
