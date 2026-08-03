---
title: Options
description: Command line flags and their environment variable equivalents.
weight: 3
---

The server accepts options as command line arguments and/or environment variables.
A command line flag takes precedence over its environment variable.

| Flag         | Environment variable | Description                                                                  |
| ------------ | -------------------- | ---------------------------------------------------------------------------- |
| `-addr`      | `YARR_ADDR`          | Address to run the server on (default `127.0.0.1:7070`)                      |
| `-base`      | `YARR_BASE`          | Base path of the service URL                                                 |
| `-auth`      | `YARR_AUTH`          | Username and password in the format `username:password`                      |
| `-auth-file` | `YARR_AUTHFILE`      | Path to a file containing `username:password`. Takes precedence over `-auth` |
| `-cert-file` | `YARR_CERTFILE`      | Path to the TLS certificate file                                             |
| `-key-file`  | `YARR_KEYFILE`       | Path to the TLS key file                                                     |
| `-db`        | `YARR_DB`            | Storage file path                                                            |
| `-log-file`  | `YARR_LOGFILE`       | Path to the log file                                                         |
| `-open`      | —                    | Open the server in the browser                                               |

## HTTPS

Both `-cert-file` and `-key-file` are required to enable HTTPS.
