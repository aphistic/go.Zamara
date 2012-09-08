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
	"testing"
)

func TestLoadHashEntry(t *testing.T) {
	// Hash table after it's been decrypted
	decryptedTable := []byte{
		0xCB, 0x37, 0x84, 0xD3,
		0xEC, 0xEA, 0xDF, 0x07,
		0x00, 0x00, 0x00, 0x00,
		0x09, 0x00, 0x00, 0x00,
	}

	entry := newHashEntry(decryptedTable)

	if entry.FilePathHashA != 0xD38437CB {
		t.Errorf("FilePathHashA - Expected: %#v Actual: %#v",
			uint32(0xD38437CB), entry.FilePathHashA)
	}
	if entry.FilePathHashB != 0x07DFEAEC {
		t.Errorf("FilePathHashB - Expected: %#v Actual: %#v",
			uint32(0x07DFEAEC), entry.FilePathHashB)
	}
	if entry.Language != 0x0000 {
		t.Errorf("Language - Expected: %#v Actual: %#v",
			0x0000, entry.Language)
	}
	if entry.Platform != 0x0000 {
		t.Errorf("Platform - Expected: %#v Actual: %#v",
			0x0000, entry.Platform)
	}
	if entry.BlockIndex != 0x00000009 {
		t.Errorf("BlockIndex - Expected: %#v Actual: %#v",
			uint32(0x00000009), entry.BlockIndex)
	}
}
