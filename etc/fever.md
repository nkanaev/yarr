# Fever API support

Fever API is a kind of RSS HTTP API interface, because the Fever API definition is not very clear, so the implementation of Fever server and Client may have some compatibility problems.

The Fever API implemented by Yarr is based on the Fever API spec: https://github.com/DigitalDJ/tinytinyrss-fever-plugin/blob/master/fever-api.md.

Here are some Apps that have been tested to work with yarr.  Feel free to test other Clients/Apps and update the list here.

>  Different apps support different URL/Address formats.  Please note whether the URL entered has `http://` scheme and `/` suffix.

| App                                                                       | Platforms        | Config Server URL                                   |
|:------------------------------------------------------------------------- | ---------------- |:--------------------------------------------------- |
| [Reeder](https://reederapp.com/)                                          | MacOS<br>iOS     | 127.0.0.1:7070/fever<br>http://127.0.0.1:7070/fever |
| [ReadKit](https://readkit.app/)                                           | MacOS<br>iOS     | http://127.0.0.1:7070/fever                         |
| [Fluent Reader](https://github.com/yang991178/fluent-reader)              | MacOS<br>Windows | http://127.0.0.1:7070/fever/                        |
| [Unread](https://apps.apple.com/us/app/unread-an-rss-reader/id1363637349) | iOS              | http://127.0.0.1:7070/fever                         |
| [Fiery Feeds](https://voidstern.net/fiery-feeds)                          | MacOS<br>iOS     | http://127.0.0.1:7070/fever                         |

If you are having trouble using Fever, please open an issue and @icefed, thanks.
