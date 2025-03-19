#!/bin/bash

set -e

# Function to display usage information
usage() {
  echo "Usage: $0 [-version VERSION] [-outfile FILENAME]"
  echo "  -version VERSION   Set the version number (default: 0.0)"
  echo "  -outfile FILENAME  Set the output file name (default: versioninfo.rc)"
  echo ""
  echo "This script generates a Windows resource file with version information."
  exit 1
}

# Default values
version="0.0"
outfile="versioninfo.rc"

# Check if help is requested
if [[ "$1" == "-h" || "$1" == "--help" ]]; then
  usage
fi

if [ $# -eq 0 ]; then
  usage
fi

# Parse command-line options
while [[ $# -gt 0 ]]; do
  case $1 in
    -version)
      if [[ -z "$2" || "$2" == -* ]]; then
        echo "Error: Missing value for -version parameter"
        usage
      fi
      version="$2"
      shift 2
      ;;
    -outfile)
      if [[ -z "$2" || "$2" == -* ]]; then
        echo "Error: Missing value for -outfile parameter"
        usage
      fi
      outfile="$2"
      shift 2
      ;;
    *)
      echo "Error: Unknown parameter: $1"
      usage
      ;;
  esac
done

# Replace dots with commas for version_comma
version_comma="${version//./,}"

# Use a here document for the template with ENDFILE delimiter
cat <<ENDFILE > "$outfile"
1 VERSIONINFO
FILEVERSION     $version_comma,0,0
PRODUCTVERSION  $version_comma,0,0
BEGIN
  BLOCK "StringFileInfo"
  BEGIN
    BLOCK "080904E4"
    BEGIN
      VALUE "CompanyName", "Old MacDonald's Farm"
      VALUE "FileDescription", "Yet another RSS reader"
      VALUE "FileVersion", "$version"
      VALUE "InternalName", "yarr"
      VALUE "LegalCopyright", "nkanaev"
      VALUE "OriginalFilename", "yarr.exe"
      VALUE "ProductName", "yarr"
      VALUE "ProductVersion", "$version"
    END
  END
  BLOCK "VarFileInfo"
  BEGIN
    VALUE "Translation", 0x809, 1252
  END
END

1 ICON "icon.ico"
ENDFILE

# Set the correct permissions
chmod 644 "$outfile"

echo "Generated $outfile with version $version"
