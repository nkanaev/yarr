// Copyright (c) 2024 The Decred developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.
//
// Feature detection originally written by Dave Collins Feb 2019.

//go:build !purego

#include "textflag.h"

// func supportsCPUID() bool
TEXT ·supportsCPUID(SB), $8-4
	// Per the Intel 64 and IA-32 Architectures Software Developer's Manual,
	// CPUID is supported if bit 21 of the EFLAGS register can be modified.
	//
	// To that end, this works as follows:
	//
	// 1. Get the current value of EFLAGS by pushing it and popping it into AX.
	// 2. Make a copy into BX for later comparison.
	// 3. Toggle bit 21 (the EFLAGS ID bit) of AX.
	// 4. Put the modified value back into EFLAGS by pushing it and popping it
	//    into EFLAGS.  The CPU will either update bit 21 of the EFLAGS with the
	//    modified value when it supports CPUID or leave it unmodified when it
	//    does not.
	// 5. Get the potentially modified value of EFLAGS by pushing it and popping
	//    it into AX.
	// 6. Compare the original and potentially modified value (aka AX vs BX)
	// 7. CPUID is supported when they do not match since bit 21 was able to be
	//    modified.
	PUSHFQ
	POPQ AX
	MOVQ AX, BX
	XORQ $0x200000, AX
	PUSHQ AX
	POPFQ
	PUSHFQ
	POPQ AX
	CMPQ AX, BX
	JE nocpuid
	MOVB $1, ret+0(FP)
	RET
nocpuid:
	MOVB $0, ret+0(FP)
	RET

// func cpuid(eaxIn, ecxIn uint32) (eax, ebx, ecx, edx uint32)
TEXT ·cpuid(SB), NOSPLIT, $0-24
	MOVL eaxIn+0(FP), AX
	MOVL ecxIn+4(FP), CX
	CPUID
	MOVL AX, eax+8(FP)
	MOVL BX, ebx+12(FP)
	MOVL CX, ecx+16(FP)
	MOVL DX, edx+20(FP)
	RET

// func xgetbv() (eax uint32)
TEXT ·xgetbv(SB), NOSPLIT, $0-4
	MOVL $0, CX
	XGETBV
	MOVL AX, eax+0(FP)
	RET
