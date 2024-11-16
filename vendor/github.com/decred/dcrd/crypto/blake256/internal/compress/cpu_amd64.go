// Copyright (c) 2024 The Decred developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.
//
// Feature detection originally written by Dave Collins Feb 2019.  Additional
// cleanup and comments added Jul 2024.

//go:build !purego

package compress

import (
	"os"
)

var (
	// features houses the result of querying the CPU and OS for supported
	// features.
	features = querySupportedFeatures()

	hasSSE2  = features.SSE2 && os.Getenv("BLAKE256_DISABLE_SSE2") != "1"
	hasSSE41 = features.SSE41 && os.Getenv("BLAKE256_DISABLE_SSE41") != "1"
	hasAVX   = features.AVX && os.Getenv("BLAKE256_DISABLE_AVX") != "1"
)

// supportsCPUID returns true when the CPU supports the CPUID opcode.
//
//go:noescape
func supportsCPUID() bool

// cpuid provides access to the CPUID opcode.
//
//go:noescape
func cpuid(eaxIn, ecxIn uint32) (eax, ebx, ecx, edx uint32)

// xgetbv provides access to the XGETBV opcode to read the contents of the
// extended control register with ECX = 0x00.
//
//go:noescape
func xgetbv() (eax uint32)

// isBitSet returns whether or not the provided bit is set in the given test
// value.
func isBitSet(testVal uint32, bit uint8) bool {
	return testVal>>bit&1 == 1
}

// supportedFeatures houses flags that specify whether or not various features
// are supported by the CPU.
type supportedFeatures struct {
	SSE2  bool
	SSE41 bool
	AVX   bool
	AVX2  bool
}

// querySupportedFeatures returns the result of querying the CPU and OS to
// determine supported features.
func querySupportedFeatures() supportedFeatures {
	// Per CPUID—CPU Identification in Chapter 3 of the Intel 64 and IA-32
	// Architectures Software Developer's Manual, Volume 2A:
	//
	// "The ID flag (bit 21) in the EFLAGS register indicates support for the
	// CPUID instruction. If a software procedure can set and clear this flag,
	// the processor executing the procedure supports the CPUID instruction.
	// This instruction operates the same in non-64-bit modes and 64-bit mode.
	//
	// CPUID returns processor identification and feature information in the
	// EAX, EBX, ECX, and EDX registers.  The output is dependent on the
	// contents of the EAX register upon execution (in some cases, ECX as
	// well)."
	//
	// The inputs and outputs for determining various levels of SIMD support
	// that are likely relevant to BLAKE are:
	//
	// Initial EAX Value | Output
	// ------------------|------------------------------------------------
	// 0x00              | EAX = Maximum Input Value for Basic CPUID Info.
	// -------------------------------------------------------------------
	// 0x01              | ECX = Feature Information
	//                   |  Bit 0 = Streaming SIMD Extensions 3 (SSE3)
	//                   |  Bit 9 = Supplemental SSE3 (SSSE3)
	//                   |  Bit 19 = Streaming SIMD Extensions 4.1 (SSE4.1)
	//                   |  Bit 20 = Streaming SIMD Extensions 4.2 (SSE4.2)
	//                   |  Bit 27 = OS sets to enable XSAVE features (OSXSAVE)
	//                   |  Bit 28 = Advanced Vector Extensions (AVX)
	//                   | EDX = Feature Information
	//                   |  Bit 25 = Streaming SIMD Extensions (SSE)
	//                   |  Bit 26 = Streaming SIMD Extensions 2 (SSE2)
	// -------------------------------------------------------------------
	// 0x07              | EBX = Feature Information
	//                   |  Bit 5 = Advanced Vector Extensions 2 (AVX2)
	//                   |  Bit 16 = AVX-512 Foundation (AVX512F)
	//                   |  Bit 17 = AVX-512 Double and Quadword (AVX512DQ)
	//                   |  Bit 30 = AVX-512 Byte and Word (AVX512BW)
	//                   |  Bit 31 = AVX-512 Vector Length Extensions (AVX512VL)
	//
	// Note that all SSE and AVX variants also require operating system support
	// in order to properly save the additional state when doing context
	// switches.  Starting with AVX, this is signaled by the OS by setting bits
	// in the extended control register which itself requires CPU support as
	// specified by the OSXSAVE bit in the table above.
	//
	// Per Chapter 13 of the Intel 64 and IA-32 Architectures Software
	// Developer’s Manual, Volume 1, the XGETBV opcode is used to obtain the
	// aforementioned extended control register (XCR).  Per the "XSAVE-SUPPORTED
	// FEATURES AND STATE-COMPONENT BITMAPS" section, the relevant bits for
	// AVX/AVX-512 support are:
	//
	// XCR  | Output
	// -----|------------------------------------------------
	// 0x00 | EAX
	//      |  Bit 1 = SSE state (XMM registers)
	//      |  Bit 2 = AVX state (YMM registers)
	//      |  Bits 5-7 = AVX-512 state components
	//      |   Bit 5 = Opmask state (K0-K7 registers)
	//      |   Bit 6 = ZMM high 256 state (upper 256 bits of ZMM0-ZMM15 registers)
	//      |   Bit 7 = High 16 ZMM state (ZMM16-ZMM31 registers)
	const (
		eaxInputQueryMax          = 0x00
		eaxInputQueryFeatureInfo  = 0x01
		eaxInputQueryExtFeatFlags = 0x07

		ecx1OutputOSXSAVEBit = 27
		edx1OutputSSE2Bit    = 26
		ecx1OutputSSE41Bit   = 19
		ecx1OutputAVXBit     = 28
		ebx7OutputAVX2Bit    = 5

		xgbvEaxOutputSSEStateBit = 1
		xgbvEaxOutputAVXStateBit = 2
	)

	// Nothing to do if the CPU somehow does not support CPUID.  Go probably
	// won't even run on such a CPU, but as the Intel manual states, it is
	// technically required to check if CPUID is supported before querying it
	// and it's best to be safe.
	var features supportedFeatures
	if !supportsCPUID() {
		return features
	}

	// Perform initial query to determine the max CPUID input value since the
	// remaining checks are only valid if the CPU at least supports querying
	// them to begin with.
	maxEAXInput, _, _, _ := cpuid(eaxInputQueryMax, 0)
	if maxEAXInput < eaxInputQueryFeatureInfo {
		return features
	}

	// Query basic feature info to determine if the CPU supports SSE2/SSE4.1.
	//
	// Note that SSE2 is always active on amd64, so checking for it could
	// probably technically be skipped, but it doesn't really cost anything
	// extra to check for it and checking is more correct.
	_, _, ecx, edx := cpuid(eaxInputQueryFeatureInfo, 0)
	features.SSE2 = isBitSet(edx, edx1OutputSSE2Bit)
	features.SSE41 = isBitSet(ecx, ecx1OutputSSE41Bit)
	hasOSXSAVE := isBitSet(ecx, ecx1OutputOSXSAVEBit)

	// Query basic feature info to determine AVX support as well as if the OS
	// supports AVX/AVX2.  See the description above for details.
	var osSupportsAVX bool
	if hasOSXSAVE {
		eax := xgetbv()
		osSupportsSSE := isBitSet(eax, xgbvEaxOutputSSEStateBit)
		osSupportsAVX = osSupportsSSE && isBitSet(eax, xgbvEaxOutputAVXStateBit)
	}
	features.AVX = isBitSet(ecx, ecx1OutputAVXBit) && osSupportsAVX

	// Querying the supported feature info for AVX2 is only valid if the CPU at
	// least supports querying it to begin with.
	if maxEAXInput < eaxInputQueryExtFeatFlags {
		return features
	}

	// Query extended feature info to determine AVX2 support.
	_, ebx, _, _ := cpuid(eaxInputQueryExtFeatFlags, 0)
	features.AVX2 = isBitSet(ebx, ebx7OutputAVX2Bit) && osSupportsAVX

	return features
}
