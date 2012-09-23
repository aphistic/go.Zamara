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

type BlockEncryptorSuite struct{}

var _ = Suite(&BlockEncryptorSuite{})

func (s *BlockEncryptorSuite) TestBlockEncryptionTable(c *C) {
	c.Check(blockEncryptionTable[0], Equals, uint32(1439053538))
	c.Check(blockEncryptionTable[100], Equals, uint32(2690928833))
	c.Check(blockEncryptionTable[1000], Equals, uint32(2859196621))
	c.Check(blockEncryptionTable[1279], Equals, uint32(1929586796))
}

func (s *BlockEncryptorSuite) TestStringHashing(c *C) {
	hash := hashString("this is a string", 0x100)
	c.Check(hash, Equals, uint32(450484832))

	hash = hashString("this is a string", 0x200)
	c.Check(hash, Equals, uint32(2082422408))
}

func (s *BlockEncryptorSuite) TestDecryptHashTable(c *C) {
	hashTable := []byte{
		0x07, 0xf8, 0xB8, 0x55, 0x4F, 0xB4, 0x8E, 0x3C, 0x7C, 0xA8, 0x7B, 0xAC, 0xAE, 0x1A, 0x00, 0xE0, // Hash Entry 1
		0xC7, 0xC9, 0xDC, 0xC5, 0x3E, 0x6C, 0xFE, 0xC3, 0xA2, 0x02, 0x33, 0xA7, 0xB8, 0x1B, 0x6D, 0xB7, // Hash Entry 2
	}

	expectedResults := []byte{
		0xCB, 0x37, 0x84, 0xD3, 0xEC, 0xEA, 0xDF, 0x07, 0x00, 0x00, 0x00, 0x00, 0x09, 0x00, 0x00, 0x00, // Hash Entry 1
		0x4B, 0xA5, 0xC2, 0xAA, 0x95, 0x2B, 0x76, 0xF4, 0x00, 0x00, 0x00, 0x00, 0x02, 0x00, 0x00, 0x00, // Hash Entry 2
	}

	encryptor := newBlockEncryptor("(hash table)", 0x300)
	encryptor.decrypt(&hashTable)

	for idx, val := range hashTable {
		if expectedResults[idx] != byte(val) {
			c.Errorf("hashTable[%v]: Actual: %#v Expected: %#v",
				idx, val, expectedResults[idx])
		}
	}
}

func (s *BlockEncryptorSuite) TestDecryptBlockTable(c *C) {
	blockTable := []byte{
		0xA7, 0x67, 0x48, 0x3D, 0xFC, 0xD1, 0x08, 0xCA, 0xE1, 0xBC, 0x35, 0xF8, 0x97, 0xF1, 0x33, 0xE9, // Block Entry 1
		0x13, 0x52, 0xB3, 0xB3, 0x07, 0x7F, 0xC0, 0x10, 0x94, 0xF8, 0xD8, 0x0D, 0xD6, 0x1E, 0xA4, 0xD3, // Block Entry 2
	}

	expectedResults := []byte{
		0x2C, 0x00, 0x00, 0x00, 0x51, 0x02, 0x00, 0x00, 0x51, 0x02, 0x00, 0x00, 0x00, 0x02, 0x00, 0x81, // Block Entry 1
		0x7D, 0x02, 0x00, 0x00, 0xA0, 0x02, 0x00, 0x00, 0x82, 0x04, 0x00, 0x00, 0x00, 0x02, 0x00, 0x81, // Block Entry 2
	}

	encryptor := newBlockEncryptor("(block table)", 0x300)
	encryptor.decrypt(&blockTable)

	for idx, val := range blockTable {
		if expectedResults[idx] != byte(val) {
			c.Errorf("blockTable[%v]: Actual: %#v Expected: %#v",
				idx, val, expectedResults[idx])
		}
	}
}
