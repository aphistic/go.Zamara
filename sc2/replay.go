package sc2

import (
	"github.com/aphistic/go.Zamara/mpq"
	"io"
)

const (
	GAME_UNKNOWN int = iota
	GAME_1V1
	GAME_2V2
	GAME_3V3
	GAME_4V4
	GAME_FFA
	GAME_6V6
	GAME_CUSTOM
)

const (
	SPEED_UNKNOWN = iota
	SPEED_SLOWER
	SPEED_SLOW
	SPEED_NORMAL
	SPEED_FAST
	SPEED_FASTER
)

const (
	CATEGORY_UNKNOWN = iota
	CATEGORY_PRIVATE
	CATEGORY_LADDER
	CATEGORY_PUBLIC
)

type Replay struct {
	mpq     *mpq.Mpq
	players []*Player

	MapName string

	Timestamp      uint64
	TimezoneOffset int64

	GameType     int
	GameSpeed    int
	GameCategory int
}

func NewReplay(reader io.ReadSeeker) (replay *Replay, err error) {
	replay = new(Replay)
	err = replay.load(reader)
	if err != nil {
		return nil, err
	}

	return replay, nil
}

func (replay *Replay) load(reader io.ReadSeeker) (err error) {
	replay.players = make([]*Player, 0)

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

	replay.MapName = value.i(1).asString()

	return
}

func (replay *Replay) loadAttributes() (err error) {
	return
}
