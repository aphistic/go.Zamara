package sc2

import (
	. "launchpad.net/gocheck"
	"os"
)

type ReplaySuite struct{}

var _ = Suite(&ReplaySuite{})

func (s *ReplaySuite) TestNewReplay(c *C) {
	reader, err := os.Open("../mpq/testdata/replay1.SC2Replay")
	c.Assert(err, IsNil)

	replay, err := NewReplay(reader)
	c.Assert(err, IsNil)

	c.Check(replay.MapName, Equals, "Discord IV")
}
