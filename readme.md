# yarr

yet another rss reader.

![screenshot](https://github.com/nkanaev/yarr/blob/master/artwork/promo.png?raw=true)

*yarr* is a server written in Go with the frontend in Vue.js. The storage is backed by SQLite.

The goal of the project is to provide a desktop application accessible via web browser.
Longer-term plans include a self-hosted solution for individuals.

[download](https://github.com/nkanaev/yarr/releases/latest)

## Building With Docker

To build with Docker

1. Clone this git repo to your machine and `cd` into that directory

2. Run `docker build -t yarr .`

3. Create a data directory to store persistent yarr data `mkdir $HOME/yarr-data`

4. Run `docker run -p 7070:7070 -v $HOME/yarr-data:/data/ yarr`

## credits

[Feather](http://feathericons.com/) for icons.
