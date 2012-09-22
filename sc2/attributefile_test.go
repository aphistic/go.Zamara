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
