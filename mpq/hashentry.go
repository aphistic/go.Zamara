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

package mpq

import (
	"encoding/binary"
	"encoding/xml"
)

type HashEntry struct {
	XMLName xml.Name `xml:"hashEntry"`

	FilePathHashA uint32 `xml:"filePathHashA"`
	FilePathHashB uint32 `xml:"filePathHashB"`
	Language      uint16 `xml:"language"`
	Platform      uint16 `xml:"platform"`
	BlockIndex    uint32 `xml:"blockIndex"`
}

func newHashEntry(data []byte) (entry *HashEntry) {
	entry = new(HashEntry)

	entry.FilePathHashA = binary.LittleEndian.Uint32(data[:4])
	entry.FilePathHashB = binary.LittleEndian.Uint32(data[0x04 : 0x04+4])
	entry.Language = binary.LittleEndian.Uint16(data[0x08 : 0x08+2])
	entry.Platform = binary.LittleEndian.Uint16(data[0x0A : 0x0A+2])
	entry.BlockIndex = binary.LittleEndian.Uint32(data[0x0C : 0x0C+4])

	return
}
