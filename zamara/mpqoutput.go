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
	"../mpq"
	"fmt"
	"os"
)

func mpqOutput(flags *zamaraFlags) {
	mpq, err := mpq.LoadMpq(flags.input)
	if err != nil {
		fmt.Printf("There was an error loading the MPQ:\n%v\n", err.Error())
		os.Exit(0)
	}

	printMpqDetails(mpq)
	fmt.Print("\n")
	printHashTable(mpq)
	fmt.Print("\n")
	printBlockTable(mpq)

	os.Exit(1234)
}

func printMpqDetails(mpq *mpq.Mpq) {
	fmt.Print("MPQ Details\n")
	fmt.Print("===========\n")

	fmt.Printf("ArchiveOffset: %v (%#v)\n",
		mpq.ArchiveOffset, mpq.ArchiveOffset)
	fmt.Printf("HeaderSize: %v (%#v)\n",
		mpq.Header.HeaderSize, mpq.Header.HeaderSize)
	fmt.Printf("ArchiveSize: %v (%#v)\n",
		mpq.Header.ArchiveSize, mpq.Header.ArchiveSize)
	fmt.Printf("FormatVersion: %v (%#v)\n",
		mpq.Header.FormatVersion, mpq.Header.FormatVersion)
	fmt.Printf("BlockSize: %v (%#v)\n",
		mpq.Header.BlockSize, mpq.Header.BlockSize)
	fmt.Printf("HashTableOffset: %v (%#v)\n",
		mpq.Header.HashTableOffset, mpq.Header.HashTableOffset)
	fmt.Printf("BlockTableOffset: %v (%#v)\n",
		mpq.Header.BlockTableOffset, mpq.Header.BlockTableOffset)
	fmt.Printf("HashTableEntries: %v (%#v)\n",
		mpq.Header.HashTableEntries, mpq.Header.HashTableEntries)
	fmt.Printf("BlockTableEntries: %v (%#v)\n",
		mpq.Header.BlockTableEntries, mpq.Header.BlockTableEntries)
	fmt.Printf("ExtendedBlockTableOffset: %v (%#v)\n",
		mpq.Header.ExtendedBlockTableOffset,
		mpq.Header.ExtendedBlockTableOffset)
	fmt.Printf("HashTableOffsetHigh: %v (%#v)\n",
		mpq.Header.HashTableOffsetHigh,
		mpq.Header.HashTableOffsetHigh)
	fmt.Printf("BlockTableOffsetHigh: %v (%#v)\n",
		mpq.Header.BlockTableOffsetHigh,
		mpq.Header.BlockTableOffsetHigh)
}

func printHashTable(mpq *mpq.Mpq) {
	fmt.Print("Hash Table\n")
	fmt.Print("==========\n")

	fmt.Print("\tFilePathiHashA\tFilePathHashB\tLanguage\tPlatform\tBlockIndex\n")
	fmt.Print("\t-------------\t-------------\t--------\t--------\t----------\n")
	for idx, entry := range mpq.HashEntries {
		fmt.Printf("%v:\t%#v\t%#v\t%#v\t\t%#v\t\t%v\n",
			idx, entry.FilePathHashA, entry.FilePathHashB,
			entry.Language, entry.Platform, entry.BlockIndex)
	}
}

func printBlockTable(mpq *mpq.Mpq) {
	fmt.Print("Block Table\n")
	fmt.Print("===========\n")

	fmt.Printf("\t%12v\t%14v\t%8v\t%-10v\n",
		"FilePosition", "CompressedSize", "FileSize", "Flags")
	fmt.Printf("\t%12v\t%14v\t%8v\t%-10v\n",
		"------------", "--------------", "--------", "-----")
	for idx, entry := range mpq.BlockEntries {
		fmt.Printf("%v:\t%#-12v\t%-14v\t%-8v\t%#-10v\n",
			idx, entry.FilePosition, entry.CompressedSize,
			entry.FileSize, entry.Flags)
	}
}
