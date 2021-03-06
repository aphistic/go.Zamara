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

package sc2

import (
	. "launchpad.net/gocheck"
	"os"
)

type AttributeFileSuite struct{}

var _ = Suite(&AttributeFileSuite{})

func (s *AttributeFileSuite) TestReadAttributes(c *C) {
	fi, err := os.Stat("testdata/attributes.dat")
	c.Assert(err, IsNil)

	reader, err := os.Open("testdata/attributes.dat")
	c.Assert(err, IsNil)

	buffer := make([]byte, fi.Size())
	reader.Read(buffer)
	file, err := newAttributeFile(buffer)
	c.Assert(err, IsNil)

	c.Check(len(file.attributes), Equals, 80)

	c.Check(file.attributes[0].header, Equals, uint32(0x000003E7))
	c.Check(file.attributes[0].id, Equals, uint32(0x0BBF))
	c.Check(file.attributes[0].playerId, Equals, uint8(0x01))
	c.Check(file.attributes[0].value, Equals, uint32(0x50617274))

	c.Check(file.attributes[4].header, Equals, uint32(0x000003E7))
	c.Check(file.attributes[4].id, Equals, uint32(0x01F4))
	c.Check(file.attributes[4].playerId, Equals, uint8(0x01))
	c.Check(file.attributes[4].value, Equals, uint32(0x48756D6E))
	c.Check(file.attributes[4].strValue, Equals, "Humn")

	c.Check(file.attributes[5].header, Equals, uint32(0x000003E7))
	c.Check(file.attributes[5].id, Equals, uint32(0x01F4))
	c.Check(file.attributes[5].playerId, Equals, uint8(0x02))
	c.Check(file.attributes[5].value, Equals, uint32(0x48756D6E))
	c.Check(file.attributes[5].strValue, Equals, "Humn")

	c.Check(file.attributes[8].header, Equals, uint32(0x000003E7))
	c.Check(file.attributes[8].id, Equals, uint32(0x0BB9))
	c.Check(file.attributes[8].playerId, Equals, uint8(0x01))
	c.Check(file.attributes[8].value, Equals, uint32(0x50726F74))
	c.Check(file.attributes[8].strValue, Equals, "Prot")

	c.Check(file.attributes[24].header, Equals, uint32(0x000003E7))
	c.Check(file.attributes[24].id, Equals, uint32(0x0BBB))
	c.Check(file.attributes[24].playerId, Equals, uint8(0x01))
	c.Check(file.attributes[24].value, Equals, uint32(0x20313030))
	c.Check(file.attributes[24].strValue, Equals, "100")

	c.Check(file.attributes[28].header, Equals, uint32(0x000003E7))
	c.Check(file.attributes[28].id, Equals, uint32(0x07D6))
	c.Check(file.attributes[28].playerId, Equals, uint8(0x01))
	c.Check(file.attributes[28].value, Equals, uint32(0x00005431))
	c.Check(file.attributes[28].strValue, Equals, "T1")

	c.Check(file.attributes[75].header, Equals, uint32(0x000003E7))
	c.Check(file.attributes[75].id, Equals, uint32(0x0BB8))
	c.Check(file.attributes[75].playerId, Equals, uint8(0x10))
	c.Check(file.attributes[75].value, Equals, uint32(0x46617372))
	c.Check(file.attributes[75].strValue, Equals, "Fasr")
}
