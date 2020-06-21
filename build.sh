#!/bin/sh

CGO_ENABLED=1 go build -tags sqlite_foreign_keys
