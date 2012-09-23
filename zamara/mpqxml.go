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
	"encoding/xml"
	"fmt"
	"github.com/aphistic/go.Zamara/mpq"
	"os"
)

func mpqXml(flags zamaraFlags) {
	if len(flags.outputAbs) <= 0 {
		fmt.Printf("An output file must be specified.")
		usage()
	}

	reader, err := os.Open(flags.inputAbs)
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

	if flags.verbose {
		fmt.Printf("Reading MPQ: %v\n", flags.input)
	}

	xmlOutput, err := xml.MarshalIndent(mpq, "", "    ")
	if err != nil {
		fmt.Printf("Error: %v\n", err.Error())
	}

	writer, err := os.Create(flags.outputAbs)
	if err != nil {
		fmt.Fprintf(os.Stderr,
			"Unable to open output file: %v", flags.outputAbs)
		os.Exit(1)
	}

	_, err = writer.WriteString("<?xml version=\"1.0\" encoding=\"utf-8\" ?>")
	if err != nil {
		fmt.Printf("There was a problem writing the XML file.\n%v", err.Error())
		os.Exit(1)
	}
	_, err = writer.Write(xmlOutput)
	if err != nil {
		fmt.Printf("There was a problem writing the XML file.\n%v", err.Error())
		os.Exit(1)
	}
	_, err = writer.WriteString("\n")
	if err != nil {
		fmt.Printf("There was a problem writing the XML file.\n%v", err.Error())
		os.Exit(1)
	}

	if flags.verbose {
		fmt.Printf("The MPQ's XML was written to %v", flags.output)
	}

	os.Exit(0)
}
