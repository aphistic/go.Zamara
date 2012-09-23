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
