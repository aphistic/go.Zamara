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
	. "launchpad.net/gocheck"
)

type HashEntrySuite struct {}
var _ = Suite(&HashEntrySuite{})

func (s *HashEntrySuite) TestLoadHashEntry(c *C) {
	// Hash table after it's been decrypted
	decryptedTable := []byte{
		0xCB, 0x37, 0x84, 0xD3,
		0xEC, 0xEA, 0xDF, 0x07,
		0x00, 0x00, 0x00, 0x00,
		0x09, 0x00, 0x00, 0x00,
	}

	entry := newHashEntry(decryptedTable)

	c.Check(entry.FilePathHashA, Equals, uint32(0xD38437CB))
	c.Check(entry.FilePathHashB, Equals, uint32(0x07DFEAEC))
	c.Check(entry.Language, Equals, uint16(0x0000))
	c.Check(entry.Platform, Equals, uint16(0x0000))
	c.Check(entry.BlockIndex, Equals, uint32(0x00000009))
}
