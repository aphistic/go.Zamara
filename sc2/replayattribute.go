package sc2

import (
	"encoding/binary"
	"strings"
)

const (
	AttrGGameType     = 0x07D1
	AttrGGameSpeed    = 0x0BB8
	AttrGGameCategory = 0x0BC1
)

const (
	AttrPType       = 0x01F4
	AttrPChosenRace = 0x0BB9
	AttrPDifficulty = 0x0BBC
	AttrPHandicap   = 0x0BBB
	AttrPNamedColor = 0x0BBA
)

type replayAttribute struct {
	header   uint32
	id       uint32
	playerId uint8
	value    uint32
	strValue string
}

func newReplayAttribute(data []byte) (attr *replayAttribute, size int64, err error) {
	attr = new(replayAttribute)
	size, err = attr.load(data)
	if err != nil {
		return nil, 0, err
	}

	return attr, size, nil
}

func readAttributeString(data []byte) (value string) {
	result := make([]byte, len(data))

	result[0] = data[3]
	result[1] = data[2]
	result[2] = data[1]
	result[3] = data[0]

	value = string(result)
	value = strings.Replace(value, "\x00", "", -1)
	value = strings.TrimSpace(value)
	return
}

func (attr *replayAttribute) load(data []byte) (size int64, err error) {
	attr.header = binary.LittleEndian.Uint32(data[:0x04])
	attr.id = binary.LittleEndian.Uint32(data[0x04 : 0x04+4])
	attr.playerId = uint8(data[0x08])
	attr.value = binary.LittleEndian.Uint32(data[0x09 : 0x09+4])

	attr.strValue = readAttributeString(data[0x09 : 0x09+4])

	return 13, nil
}

func (attr *replayAttribute) isGlobal() (result bool) {
	return attr.playerId == 0x10
}

func (attr *replayAttribute) isPlayer() (result bool) {
	return !attr.isGlobal()
}
