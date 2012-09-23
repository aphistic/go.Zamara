/* go.Zamara Library
 * Copyright (c) 2012, Erik Davidson
 * All rights reserved.
 * 
 * Redistribution and use in source and binary forms, with or without
 * modification, are permitted provided that the following conditions
 * are met:
 *
 * 1. Redistributions of source code must retain the above copyright notice,
 *    this list of conditions and the following disclaimer.
 * 
 * 2. Redistributions in binary form must reproduce the above copyright notice,
 *    this list of conditions and the following disclaimer in the documentation
 *    and/or other materials provided with the distribution.
 *
 * THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
 * AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
 * IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE
 * ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE
 * LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR
 * CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF
 * SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS
 * INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN
 * CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE)
 * ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE
 * POSSIBILITY OF SUCH DAMAGE.
 */

package main

import (
	"fmt"
	"github.com/aphistic/go.Zamara/mpq"
	"os"
)

func mpqStdout(flags zamaraFlags) {
	reader, err := os.Open(flags.input)
	if err != nil {
		fmt.Fprintf(os.Stderr,
			"Unable to open MPQ (%v): %v", flags.input, err.Error())
		os.Exit(1)
	}

	mpq, err := mpq.NewMpq(reader)
	if err != nil {
		fmt.Fprintf(os.Stderr,
			"Unable to read MPQ: %v", err.Error())
		os.Exit(1)
	}

	fmt.Printf("Reading MPQ: %v\n", flags.input)
	mpqStdoutHeader(flags, mpq)
}

func mpqStdoutHeader(flags zamaraFlags, mpq *mpq.Mpq) {
	fmt.Printf("Archive Offset: %v\n\n", mpq.ArchiveOffset)
	fmt.Printf("\n")

	fmt.Printf("Header\n")
	fmt.Printf("======\n")
	fmt.Printf("Header Size: %v\n", mpq.Header.HeaderSize)
	fmt.Printf("Archive Size: %v\n", mpq.Header.ArchiveSize)
	fmt.Printf("Format Version: %v\n", mpq.Header.FormatVersion)
	fmt.Printf("Block Size: %v\n", mpq.Header.BlockSize)
	fmt.Printf("Hash Table Offset: %v\n", mpq.Header.HashTableOffset)
	fmt.Printf("Block Table Offset: %v\n", mpq.Header.BlockTableOffset)
	fmt.Printf("Hash Table Entries: %v\n", mpq.Header.HashTableEntries)
}
