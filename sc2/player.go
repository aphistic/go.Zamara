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

const (
	PlayerUnknown = iota
	PlayerHuman
	PlayerComputer
)

const (
	RaceUnknown = iota
	RaceRandom
	RaceTerran
	RaceProtoss
	RaceZerg
)

const (
	DifficultyUnknown = iota
	DifficultyVeryEasy
	DifficultyEasy
	DifficultyMedium
	DifficultyHard
	DifficultyVeryHard
	DifficultyInsane
)

const (
	ColorUnknown = iota
	ColorRed
	ColorBlue
	ColorTeal
	ColorPurple
	ColorYellow
	ColorOrange
	ColorGreen
	ColorLightPink
	ColorViolet
	ColorLightGrey
	ColorDarkGreen
	ColorBrown
	ColorLightGreen
	ColorDarkGrey
	ColorPink
)

type Player struct {
	Name string
	Id   int64

	Type int

	Team       int
	Color      Color
	NamedColor int

	ChosenRace int
	ActualRace int
	Difficulty int
	Handicap   int

	Outcome int
}

func newPlayer(value *serializedValue) (player *Player, err error) {
	player = new(Player)
	player.load(value)

	return
}

func (player *Player) load(value *serializedValue) (err error) {
	raceMap := map[string]int{
		"Protoss": RaceProtoss,
		"Terran":  RaceTerran,
		"Zerg":    RaceZerg,
	}

	player.Name = value.i(0).asString()
	player.Id = value.i(1).i(3).asInt64()
	player.Color.A = int(value.i(3).i(0).asInt64())
	player.Color.R = int(value.i(3).i(1).asInt64())
	player.Color.G = int(value.i(3).i(2).asInt64())
	player.Color.B = int(value.i(3).i(3).asInt64())
	player.Team = int(value.i(5).asInt64())
	player.Handicap = int(value.i(6).asInt64())
	player.Outcome = int(value.i(8).asInt64())

	val, exists := raceMap[value.i(2).asString()]
	if !exists {
		val = RaceUnknown
	}
	player.ActualRace = val

	return
}
