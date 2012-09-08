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
	fmt.Printf("%v", mpq.ArchiveOffset)
}
