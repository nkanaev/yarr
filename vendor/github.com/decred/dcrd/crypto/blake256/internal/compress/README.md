compress
========

[![Build Status](https://github.com/decred/dcrd/workflows/Build%20and%20Test/badge.svg)](https://github.com/decred/dcrd/actions)
[![ISC License](https://img.shields.io/badge/license-ISC-blue.svg)](http://copyfree.org)
[![Doc](https://img.shields.io/badge/doc-reference-blue.svg)](https://pkg.go.dev/github.com/decred/dcrd/crypto/blake256/internal/compress)

## Overview

Package `compress` implements the BLAKE-224 and BLAKE-256 block compression
function.  It provides a pure Go implementation as well as specialized
implementations that take advantage of vector extensions (SSE2, SSE4.1, and AVX)
on the `amd64` architecture when they are supported.

The package detects hardware support and arranges for the exported `Blocks`
function to automatically use the fastest available supported hardware
extensions that are not disabled.

## Tests and Benchmarks

The package also provides full tests for all implementations as well as
benchmarks.  However, do note that since the specialized implementations require
hardware support, the tests and benchmarks for them will be skipped when running
on hardware that does not support the required extensions.

It is possible to test all implementations without hardware support by using
software such as the [Intel Software Development Emulator](https://www.intel.com/content/www/us/en/developer/articles/tool/software-development-emulator.html).

Some relevant flags for testing purposes with the Intel SDE are:

* SSE2:  `-p4p  Set chip-check and CPUID for Intel(R) Pentium4 Prescott CPU`
* SSE41: `-pnr  Set chip-check and CPUID for Intel(R) Penryn CPU`
* AVX:   `-snb  Set chip-check and CPUID for Intel(R) Sandy Bridge CPU`

## Disabling Assembler Optimizations

The `purego` build tag may be used to disable all assembly code.

Additionally, when built normally without the `purego` build tag, the assembly
optimizations for each of the supported vector extensions can individually be
disabled at runtime by setting the following environment variables to `1`.

* `BLAKE256_DISABLE_AVX=1`: Disable Advanced Vector Extensions (AVX) optimizations
* `BLAKE256_DISABLE_SSE41=1`: Disable Streaming SIMD Extensions 4.1 (SSE4.1) optimizations
* `BLAKE256_DISABLE_SSE2=1`: Disable Streaming SIMD Extensions 2 (SSE2) optimizations

## Installation and Updating

This package is internal and therefore is neither directly installed nor needs
to be manually updated.

## License

Package compress is licensed under the [copyfree](http://copyfree.org) ISC
License.
