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
	"fmt"
	"github.com/aphistic/go.Zamara/mpq"
	"io"
	"time"
)

const (
	GameUnknown int = iota
	Game1v1
	Game2v2
	Game3v3
	Game4v4
	GameFfa
	Game6v6
	GameCustom
)

const (
	SpeedUnknown = iota
	SpeedSlower
	SpeedSlow
	SpeedNormal
	SpeedFast
	SpeedFaster
)

const (
	CategoryUnknown = iota
	CategoryPrivate
	CategoryLadder
	CategoryPublic
)

type Replay struct {
	mpq *mpq.Mpq

	MapName string

	Timestamp time.Time

	GameType     int
	GameSpeed    int
	GameCategory int

	Players []*Player
}

func NewReplay(reader io.ReadSeeker) (replay *Replay, err error) {
	fmt.Printf("")
	replay = new(Replay)
	err = replay.load(reader)
	if err != nil {
		return nil, err
	}

	return replay, nil
}

func (replay *Replay) load(reader io.ReadSeeker) (err error) {
	replay.Players = make([]*Player, 0)

	replay.mpq, err = mpq.NewMpq(reader)
	if err != nil {
		return
	}

	err = replay.loadDetails()
	if err != nil {
		return
	}
	replay.loadAttributes()
	if err != nil {
		return
	}

	return
}

func (replay *Replay) loadDetails() (err error) {
	file, err := replay.mpq.File("replay.details")
	if err != nil {
		return
	}

	buffer := make([]byte, file.FileSize)
	_, err = replay.mpq.Read(buffer)
	if err != nil && err != mpq.EOF {
		return
	}

	value, _, err := newSerializedValue(buffer)
	if err != nil {
		return
	}

	totalPlayers := value.i(0).size()
	replay.Players = make([]*Player, totalPlayers)
	for idx := int64(0); idx < totalPlayers; idx++ {
		player, _ := newPlayer(value.i(0).i(idx))
		replay.Players[idx] = player
	}

	replay.MapName = value.i(1).asString()

	tzo := value.i(6).asInt64()
	tzo = tzo / 10000000
	loc := time.FixedZone("unknown", int(tzo))

	ts := value.i(5).asInt64()
	ts = (ts - 116444735995904000) / 10000000
	utc := time.Unix(ts, 0).UTC()
	t := time.Date(utc.Year(), utc.Month(), utc.Day(),
		utc.Hour(), utc.Minute(), utc.Second(),
		utc.Nanosecond(), loc)
	replay.Timestamp = t

	return
}

func (replay *Replay) loadAttributes() (err error) {
	file, err := replay.mpq.File("replay.attributes.events")
	if err != nil {
		return
	}

	buffer := make([]byte, file.FileSize)
	_, err = replay.mpq.Read(buffer)
	if err != nil && err != mpq.EOF {
		return
	}

	attrs, err := newAttributeFile(buffer)
	if err != nil {
		return
	}

	for idx := 0; idx < len(attrs.attributes); idx++ {
		attr := attrs.attributes[idx]
		if attr.isPlayer() {
			replay.processPlayerAttribute(attr)
		} else {
			replay.processGlobalAttribute(attr)
		}
	}

	return
}

func (replay *Replay) processPlayerAttribute(attr *replayAttribute) (err error) {
	playerIdx := attr.playerId - 1

	typeMap := map[string]int{
		"Humn": PlayerHuman,
		"Comp": PlayerComputer,
	}
	raceMap := map[string]int{
		"RAND": RaceRandom,
		"Terr": RaceTerran,
		"Prot": RaceProtoss,
		"Zerg": RaceZerg,
	}
	difficultyMap := map[string]int{
		"VyEy": DifficultyVeryEasy,
		"Easy": DifficultyEasy,
		"Medi": DifficultyMedium,
		"Hard": DifficultyHard,
		"VyHd": DifficultyVeryHard,
		"Insa": DifficultyInsane,
	}
	colorMap := map[string]int{
		"tc01": ColorRed,
		"tc02": ColorBlue,
		"tc03": ColorTeal,
		"tc04": ColorPurple,
		"tc05": ColorYellow,
		"tc06": ColorOrange,
		"tc07": ColorGreen,
		"tc08": ColorLightPink,
		"tc09": ColorViolet,
		"tc10": ColorLightGrey,
		"tc11": ColorDarkGreen,
		"tc12": ColorBrown,
		"tc13": ColorLightGreen,
		"tc14": ColorDarkGrey,
		"tc15": ColorPink,
	}

	switch attr.id {
	case AttrPType:
		val, exists := typeMap[attr.strValue]
		if !exists {
			val = PlayerUnknown
		}
		replay.Players[playerIdx].Type = val
		break
	case AttrPChosenRace:
		val, exists := raceMap[attr.strValue]
		if !exists {
			val = RaceUnknown
		}
		replay.Players[playerIdx].ChosenRace = val
		break
	case AttrPDifficulty:
		val, exists := difficultyMap[attr.strValue]
		if !exists {
			val = DifficultyUnknown
		}
		replay.Players[playerIdx].Difficulty = val
		break
	case AttrPNamedColor:
		val, exists := colorMap[attr.strValue]
		if !exists {
			val = ColorUnknown
		}
		replay.Players[playerIdx].NamedColor = val
		break
	}

	return
}

func (replay *Replay) processGlobalAttribute(attr *replayAttribute) (err error) {
	typeMap := map[string]int{
		"1v1":  Game1v1,
		"2v2":  Game2v2,
		"3v3":  Game3v3,
		"4v4":  Game4v4,
		"FFA":  GameFfa,
		"6v6":  Game6v6,
		"Cust": GameCustom,
	}
	speedMap := map[string]int{
		"Slor": SpeedSlower,
		"Slow": SpeedSlow,
		"Norm": SpeedNormal,
		"Fast": SpeedFast,
		"Fasr": SpeedFaster,
	}
	categoryMap := map[string]int{
		"Priv": CategoryPrivate,
		"Amm":  CategoryLadder,
		"Pub":  CategoryPublic,
	}

	switch attr.id {
	case AttrGGameType:
		val, exists := typeMap[attr.strValue]
		if !exists {
			val = GameUnknown
		}
		replay.GameType = val
		break
	case AttrGGameSpeed:
		val, exists := speedMap[attr.strValue]
		if !exists {
			val = SpeedUnknown
		}
		replay.GameSpeed = val
		break
	case AttrGGameCategory:
		val, exists := categoryMap[attr.strValue]
		if !exists {
			val = CategoryUnknown
		}
		replay.GameCategory = val
		break
	}

	return
}
