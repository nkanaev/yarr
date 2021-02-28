## install (linux)

The Linux version doesn't come with the desktop environment integration.
For easy access you can manually create a desktop menu entry by
by following the steps below:

    # get the latest linux binary zip
    unzip -x yarr*.zip
    sudo mv yarr /usr/local/bin/yarr
    sudo nano /usr/local/share/applications/yarr.desktop

and paste the content below:

    [Desktop Entry]
    Name=yarr
    Exec=/usr/loca/bin -open
    Icon=rss
    Type=Application
    Categories=Internet;
