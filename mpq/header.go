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
	"encoding/xml"
)

type Header struct {
	XMLName xml.Name `xml:"header"`

	HeaderSize    uint32 `xml:"headerSize"`
	ArchiveSize   uint32 `xml:"archiveSize"`
	FormatVersion uint16 `xml:"formatVersion"`

	BlockSize uint16 `xml:"blockSize"`

	HashTableOffset   uint32 `xml:"hashTableOffset"`
	BlockTableOffset  uint32 `xml:"blockTableOffset"`
	HashTableEntries  uint32 `xml:"hashTableEntries"`
	BlockTableEntries uint32 `xml:"blockTableEntries"`

	ExtendedBlockTableOffset uint64 `xml:"extendedBlockTableOffset"`
	HashTableOffsetHigh      uint16 `xml:"hashTableOffsetHigh"`
	BlockTableOffsetHigh     uint16 `xml:"blockTableOffsetHigh"`
}
