# yarr

**yarr** (yet another rss reader) is a web-based feed aggregator which can be used both
as a desktop application and a personal self-hosted server.

The app is a single binary with an embedded database (SQLite).

![screenshot](etc/promo.png)

## usage

The latest prebuilt binaries for Linux/MacOS/Windows AMD64 are available
[here](https://github.com/nkanaev/yarr/releases/latest). Installation instructions:

* Command Arges
  
  ```
    -addr string
          address to run server on (default "127.0.0.1:7070")
    -auth-file path
          path to a file containing username:password
    -base string
          base path of the service url
    -cert-file path
          path to cert file for https
    -db path
          storage file path
    -key-file path
          path to key file for https
    -log-file path
          path to log file to use instead of stdout
    -open
          open the server in browser
    -version
          print application version
  ```

* MacOS

  Download `yarr-*-macos64.zip`, unzip it, place `yarr.app` in `/Applications` folder, [open the app][macos-open], click the anchor menu bar icon, select "Open".

* Windows

  Download `yarr-*-windows64.zip`, unzip it, open `yarr.exe`, click the anchor system tray icon, select "Open".

* Linux

  Download `yarr-*-linux64.zip`, unzip it, place `yarr` in `$HOME/.local/bin`
and run [the script](etc/install-linux.sh).

[macos-open]: https://support.apple.com/en-gb/guide/mac-help/mh40616/mac

* Docker environment
  
  You can use docker or docker-compose to run yarr, and you can also use environment variables to configure startup parameters.
  
  - `YARR_ADDR` ：address to run server on (default "127.0.0.1:7070")
  - `YARR_BASE` ：base path of the service url
  - `YARR_AUTHFILE` ：path to a file containing username:password
  - `YARR_CERTFILE` ：path to cert file for https
  - `YARR_KEYFILE` ：path to key file for https
  - `YARR_DB` ：storage file path
  - `YARR_LOGFILE` ：path to log file to use instead of stdout
  
* Docker run：
  ```
  docker run -d \
  --name yarr \
  -p 25255:7070 \
  -e YARR_AUTHFILE="/data/.auth.list" \
  -v /data/yarr-data:/data \
  --restart always \
  arsfeld/yarr:latest
  ```
  
* Docker-Compose Run

  Create a file named `.auth.list` under the `/data/` directory, and the content format should be: `username:password`.
  Then start by running docker-compose up -d and enjoy!

  ```yaml
  version: '3.3'
  services:
      yarr:
          container_name: yarr
          image: 'arsfeld/yarr:latest'
          restart: always
          ports:
              - '25255:7070'
          environment:
            YARR_AUTHFILE: "/data/.auth.list"
          volumes:
              - '/data/yarr-data:/data'
  ```
  
* See more:

  * [Building from source code](doc/build.md)
  * [Fever API support](doc/fever.md)

## credits

[Feather](http://feathericons.com/) for icons.
