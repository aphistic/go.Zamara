package main

import (
	"fmt"
	pmpq "github.com/aphistic/go.Zamara/mpq"
	"log"
	"os"
)

func extractMpq(flags zamaraFlags) {
	err := os.MkdirAll(flags.output, 0755)
	if err != nil {
		fmt.Printf("Could not access output location: %v\n", err.Error())
		os.Exit(1)
	}

	reader, err := os.Open(flags.input)
	if err != nil {
		log.Printf("Error opening MPQ: %v\n", err.Error())
		os.Exit(1)
	}
	mpq, err := pmpq.NewMpq(reader)
	if err != nil {
		log.Printf("Error reading MPQ: %v\n", err.Error())
		os.Exit(1)
	}

	for _, file := range mpq.Files() {
		cleanPath := cleanDirectory(flags.output) +
			cleanFilename(file.Filename)
		if flags.verbose {
			fmt.Printf("Extracting %v to %v\n",
				file.Filename, cleanPath)
		}

		_, err = mpq.File(file.Filename)
		if err != nil {
			fmt.Printf("Error extracting file %v\n%v\n",
				file.Filename, err.Error())
			continue
		}

		osFile, err := os.Create(cleanPath)
		defer osFile.Close()
		if err != nil {
			fmt.Printf("Error creating file %v\n%v\n",
				cleanPath, err.Error())
			continue
		}

		buffer := make([]byte, file.FileSize)
		_, err = mpq.Read(buffer)
		if err != nil && err != pmpq.EOF {
			fmt.Printf("Error reading file %v from MPQ\n%v\n",
				file.Filename, err.Error())
			continue
		}

		_, err = osFile.Write(buffer)
		if err != nil {
			fmt.Printf("Error writing file %v to %v\n%v\n",
				file.Filename, cleanPath, err.Error())
			continue
		}
	}

	os.Exit(0)
}
