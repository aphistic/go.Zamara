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
