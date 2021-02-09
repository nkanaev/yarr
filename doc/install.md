# Linux desktop

Grab the latest linux binary, then run:

```
$ sudo mv /path/to/yarr /usr/local/bin
$ sudo tee /usr/local/share/applications/yarr.desktop >/dev/null <<EOF
[Desktop Entry]
Name=yarr
Exec=yarr -open
Icon=rss
Type=Application
Categories=Internet;
EOF
```
