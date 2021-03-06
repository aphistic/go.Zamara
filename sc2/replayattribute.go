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
