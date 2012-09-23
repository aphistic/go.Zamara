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
	"time"
)

type ReplaySuite struct{}

var _ = Suite(&ReplaySuite{})

func (s *ReplaySuite) TestNewReplay(c *C) {
	reader, err := os.Open("../mpq/testdata/replay1.SC2Replay")
	c.Assert(err, IsNil)

	replay, err := NewReplay(reader)
	c.Assert(err, IsNil)

	c.Check(replay.MapName, Equals, "Discord IV")

	tz := time.FixedZone("unknown", -14400)
	time := time.Date(2010, 7, 31, 3, 56, 54, 0, tz)
	c.Check(replay.Timestamp, DeepEquals, time)

	c.Check(replay.GameType, Equals, Game2v2)
	c.Check(replay.GameSpeed, Equals, SpeedFaster)
	c.Check(replay.GameCategory, Equals, CategoryLadder)

	// Players
	c.Check(len(replay.Players), Equals, 4)

	// First Player
	player := replay.Players[0]
	c.Check(player.Name, Equals, "TehPartE")
	c.Check(player.Type, Equals, PlayerHuman)
	c.Check(player.ChosenRace, Equals, RaceProtoss)
	c.Check(player.ActualRace, Equals, RaceProtoss)
	c.Check(player.Difficulty, Equals, DifficultyMedium)
	c.Check(player.NamedColor, Equals, ColorRed)
	c.Check(player.Color.A, Equals, 255)
	c.Check(player.Color.R, Equals, 180)
	c.Check(player.Color.G, Equals, 20)
	c.Check(player.Color.B, Equals, 30)
	c.Check(player.Team, Equals, 0)
	c.Check(player.Handicap, Equals, 100)
	c.Check(player.Outcome, Equals, 0)

	// Second Player
	player = replay.Players[1]
	c.Check(player.Name, Equals, "totsgerber")
	c.Check(player.Type, Equals, PlayerHuman)
	c.Check(player.ChosenRace, Equals, RaceZerg)
	c.Check(player.ActualRace, Equals, RaceZerg)
	c.Check(player.Difficulty, Equals, DifficultyMedium)
	c.Check(player.NamedColor, Equals, ColorBlue)
	c.Check(player.Color.A, Equals, 255)
	c.Check(player.Color.R, Equals, 0)
	c.Check(player.Color.G, Equals, 66)
	c.Check(player.Color.B, Equals, 255)
	c.Check(player.Team, Equals, 1)
	c.Check(player.Handicap, Equals, 100)
	c.Check(player.Outcome, Equals, 0)

	// Third Player
	player = replay.Players[2]
	c.Check(player.Name, Equals, "David")
	c.Check(player.Type, Equals, PlayerHuman)
	c.Check(player.ChosenRace, Equals, RaceTerran)
	c.Check(player.ActualRace, Equals, RaceTerran)
	c.Check(player.Difficulty, Equals, DifficultyMedium)
	c.Check(player.NamedColor, Equals, ColorTeal)
	c.Check(player.Color.A, Equals, 255)
	c.Check(player.Color.R, Equals, 28)
	c.Check(player.Color.G, Equals, 167)
	c.Check(player.Color.B, Equals, 234)
	c.Check(player.Team, Equals, 1)
	c.Check(player.Handicap, Equals, 100)
	c.Check(player.Outcome, Equals, 0)

	// Fourth Player
	player = replay.Players[3]
	c.Check(player.Name, Equals, "Steven")
	c.Check(player.Type, Equals, PlayerHuman)
	c.Check(player.ChosenRace, Equals, RaceTerran)
	c.Check(player.ActualRace, Equals, RaceTerran)
	c.Check(player.Difficulty, Equals, DifficultyMedium)
	c.Check(player.NamedColor, Equals, ColorPurple)
	c.Check(player.Color.A, Equals, 255)
	c.Check(player.Color.R, Equals, 84)
	c.Check(player.Color.G, Equals, 0)
	c.Check(player.Color.B, Equals, 129)
	c.Check(player.Team, Equals, 0)
	c.Check(player.Handicap, Equals, 0) // TODO: Confirm value
	c.Check(player.Outcome, Equals, 0)
}
