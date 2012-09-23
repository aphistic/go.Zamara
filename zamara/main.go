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
	"flag"
	"fmt"
	"os"
	"strings"
)

type zamaraFlags struct {
	input     string // Input file
	inputAbs  string // Absolute path of the input file
	output    string // Output file or directory
	outputAbs string // Absolute path of the output file or directory
	format    string // Output format for output
	runType   string // Output type

	verbose bool // Verbose output
}

var flags zamaraFlags

func init() {
	flag.StringVar(&flags.input, "in", "", "Input file.")
	flag.StringVar(&flags.output, "out", "", "Output file or directory.")
	flag.StringVar(&flags.format, "format", "stdout", "Output format. [stdout, json, xml, extract]")
	flag.StringVar(&flags.runType, "type", "sc2", "Output type, see below for options.")
	flag.BoolVar(&flags.verbose, "v", false, "Verbose output.")
}

func usage() {
	fmt.Fprintf(os.Stderr, "Usage:\n\n")
	fmt.Fprintf(os.Stderr, "  zamara [arguments]\n\n")

	flag.PrintDefaults()
	os.Exit(2)
}

func validateUsage() {
	// Make sure the input file exists
	_, err := os.Stat(flags.input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to find input file %v\n", flags.input)
		usage()
	}

	flags.format = strings.ToLower(flags.format)
	flags.runType = strings.ToLower(flags.runType)

	// Make sure we recognize the output format
	switch flags.format {
	case "stdout":
		break
	case "json":
		break
	case "xml":
		break
	case "extract":
		break
	default:
		fmt.Fprintf(os.Stderr,
			"Unrecognized output format: %v\n", flags.format)
		os.Exit(1234)
	}
}

func main() {
	flag.Usage = usage
	flag.Parse()

	fmt.Println(expandPath(flags.output))
	flags.inputAbs = expandPath(flags.input)
	flags.outputAbs = expandPath(flags.output)

	validateUsage()

	switch flags.format {
	case "extract":
		extractMpq(flags)
		break
	case "stdout":
		handleStdout(flags)
		break
	case "xml":
		handleXml(flags)
		break
	default:
		break
	}

	fmt.Fprintf(os.Stderr,
		"Unknown arguments lead to an unknown state, file a bug with your command line arguments to get this fixed!\n")
	os.Exit(1)
}
