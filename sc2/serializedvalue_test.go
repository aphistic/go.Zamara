package sc2

import (
	. "launchpad.net/gocheck"
	"os"
)

// Test Suite
type SerializedValueSuite struct{}

var _ = Suite(&SerializedValueSuite{})

func (s *SerializedValueSuite) TestDeserializeData(c *C) {
	fi, err := os.Stat("testdata/serialized.dat")
	c.Assert(err, IsNil)

	reader, err := os.Open("testdata/serialized.dat")
	c.Assert(err, IsNil)

	buffer := make([]byte, fi.Size())
	reader.Read(buffer)
	value, _, err := newSerializedValue(buffer)
	c.Assert(err, IsNil)

	// Root Node
	c.Check(value.isKey(), Equals, true)

	// Player Array
	c.Check(value.i(0).isArray(), Equals, true)
	c.Check(value.i(0).size(), Equals, int64(4))

	// First Player
	c.Check(value.i(0).i(0).isKey(), Equals, true)
	c.Check(value.i(0).i(0).i(0).isString(), Equals, true)
	c.Check(value.i(0).i(0).i(0).asString(), Equals, "TehPartE") // Player Name
	c.Check(value.i(0).i(0).i(1).i(3).isInt64(), Equals, true)
	c.Check(value.i(0).i(0).i(1).i(3).asInt64(), Equals, int64(278960)) // Player Id
	c.Check(value.i(0).i(0).i(2).isString(), Equals, true)
	c.Check(value.i(0).i(0).i(2).asString(), Equals, "Protoss") // Player Race
	c.Check(value.i(0).i(0).i(3).isKey(), Equals, true)
	c.Check(value.i(0).i(0).i(3).i(0).isInt64(), Equals, true)
	c.Check(value.i(0).i(0).i(3).i(0).asInt64(), Equals, int64(255)) // Color - A
	c.Check(value.i(0).i(0).i(3).i(1).isInt64(), Equals, true)
	c.Check(value.i(0).i(0).i(3).i(1).asInt64(), Equals, int64(180)) // Color - R
	c.Check(value.i(0).i(0).i(3).i(2).isInt64(), Equals, true)
	c.Check(value.i(0).i(0).i(3).i(2).asInt64(), Equals, int64(20)) // Color - G
	c.Check(value.i(0).i(0).i(3).i(3).isInt64(), Equals, true)
	c.Check(value.i(0).i(0).i(3).i(3).asInt64(), Equals, int64(30)) // Color - B
	c.Check(value.i(0).i(0).i(5).isInt64(), Equals, true)
	c.Check(value.i(0).i(0).i(5).asInt64(), Equals, int64(0)) // Team
	c.Check(value.i(0).i(0).i(6).isInt64(), Equals, true)
	c.Check(value.i(0).i(0).i(6).asInt64(), Equals, int64(100)) // Handicap
	c.Check(value.i(0).i(0).i(8).isInt64(), Equals, true)
	c.Check(value.i(0).i(0).i(8).asInt64(), Equals, int64(0)) // Outcome

	// Map Name
	c.Check(value.i(1).isString(), Equals, true)
	c.Check(value.i(1).asString(), Equals, "Discord IV")

	// Datetime of game
	c.Check(value.i(5).isInt64(), Equals, true)
	c.Check(value.i(5).asInt64(), Equals, int64(129250222145273475))

	// Timezone offset
	c.Check(value.i(6).isInt64(), Equals, true)
	c.Check(value.i(6).asInt64(), Equals, int64(-144000000000))
}
