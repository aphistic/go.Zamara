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
	"errors"
	"fmt"
)

const (
	ValueString byte = 0x02
	ValueArray       = 0x04
	ValueKey         = 0x05
	ValueInt8        = 0x06
	ValueInt32       = 0x07
	ValueIntVlf      = 0x09
)

var OutOfRange error = errors.New("Out of range")

type serializedValue struct {
	members []*serializedValue
	keys    []byte

	valueType byte

	stringValue string
	intValue    int64
}

func newSerializedValue(data []byte) (value *serializedValue, size int64, err error) {
	value = new(serializedValue)
	size, err = value.load(data)
	if err != nil {
		return nil, 0, err
	}

	return value, size, nil
}

func (value *serializedValue) load(data []byte) (size int64, err error) {
	value.members = []*serializedValue{}
	value.keys = []byte{}

	if len(data) < 1 {
		return
	}

	switch data[0] {
	case ValueString:
		dataSize, err := value.loadString(data[1:])
		if err != nil {
			return 0, err
		}
		size += dataSize
		break
	case ValueArray:
		dataSize, err := value.loadArray(data[1:])
		if err != nil {
			return 0, err
		}
		size += dataSize
		break
	case ValueKey:
		dataSize, err := value.loadKey(data[1:])
		if err != nil {
			return 0, err
		}
		size += dataSize
		break
	case ValueInt8:
		dataSize, err := value.loadInt8(data[1:])
		if err != nil {
			return 0, err
		}
		size += dataSize
		break
	case ValueInt32:
		dataSize, err := value.loadInt32(data[1:])
		if err != nil {
			return 0, err
		}
		size += dataSize
		break
	case ValueIntVlf:
		dataSize, err := value.loadIntVlf(data[1:])
		if err != nil {
			return 0, err
		}
		size += dataSize
		break
	default:
		err = fmt.Errorf("Unknown serialization type: %v", data[0])
		break
	}
	size++ // Add one for the "type" byte that's skipped

	return size, nil
}

func (value *serializedValue) loadString(data []byte) (size int64, err error) {
	value.valueType = ValueString

	var length int64
	length, amountRead := value.readIntVlf(data)
	value.stringValue = string(data[amountRead : amountRead+length])

	return length + 1, nil
}

func (value *serializedValue) loadArray(data []byte) (size int64, err error) {
	value.valueType = ValueArray

	offset := int64(2)
	elements, readSize := value.readIntVlf(data[offset:])
	offset += readSize
	value.keys = make([]byte, elements)
	value.members = make([]*serializedValue, elements)

	for idx := int64(0); idx < elements; idx++ {
		newValue, readSize, err := newSerializedValue(data[offset:])
		if err != nil {
			return 0, err
		}
		offset += readSize
		value.members[idx] = newValue
	}

	return offset, nil
}

func (value *serializedValue) loadKey(data []byte) (size int64, err error) {
	value.valueType = ValueKey

	elements, offset := value.readIntVlf(data)
	value.keys = make([]byte, elements)
	value.members = make([]*serializedValue, elements)

	for idx := int32(0); int64(idx) < elements; idx++ {
		value.keys[idx] = data[offset]
		offset++

		newValue, readSize, err := newSerializedValue(data[offset:])
		if err != nil {
			return 0, err
		}
		offset += readSize
		value.members[idx] = newValue
	}

	return offset, nil
}

func (value *serializedValue) loadInt8(data []byte) (size int64, err error) {
	value.valueType = ValueInt8
	value.intValue = int64(data[0])

	return 1, nil
}

func (value *serializedValue) loadInt32(data []byte) (size int64, err error) {
	value.valueType = ValueInt32

	value.intValue = int64(binary.LittleEndian.Uint32(data[:4]))

	return 4, nil
}

func (value *serializedValue) loadIntVlf(data []byte) (size int64, err error) {
	value.valueType = ValueIntVlf

	value.intValue, size = value.readIntVlf(data)

	return size, nil
}

func (value *serializedValue) readIntVlf(data []byte) (result int64, size int64) {
	var currentByte byte
	var byteCount uint32 = 0

	for {
		currentByte = data[byteCount]
		result += (int64(currentByte & 0x7F)) << (7 * byteCount)
		byteCount++
		if (currentByte & 0x80) <= 0 {
			break
		}
	}

	if result&1 == 1 {
		result = -(result >> 1)
	} else {
		result = (result >> 1)
	}

	return result, int64(byteCount)
}

// Type checks
func (value *serializedValue) isString() (result bool) {
	return value.valueType == ValueString
}
func (value *serializedValue) isArray() (result bool) {
	return value.valueType == ValueArray
}
func (value *serializedValue) isKey() (result bool) {
	return value.valueType == ValueKey
}
func (value *serializedValue) isInt8() (result bool) {
	return value.valueType == ValueInt8
}
func (value *serializedValue) isInt32() (result bool) {
	return value.valueType == ValueInt32
}
func (value *serializedValue) isInt64() (result bool) {
	return value.valueType == ValueIntVlf
}

// Get values
func (value *serializedValue) size() (size int64) {
	switch value.valueType {
	case ValueArray:
		fallthrough
	case ValueKey:
		return int64(len(value.members))
	case ValueString:
		return int64(len(value.stringValue))
	case ValueInt8:
		return 1
	case ValueInt32:
		return 4
	case ValueIntVlf:
		return 8
	}

	return 0
}

func (value *serializedValue) i(index int64) (item *serializedValue) {
	return value.item(index)
}

func (value *serializedValue) item(index int64) (item *serializedValue) {
	if index < 0 || index >= int64(len(value.members)) {
		return nil
	}

	return value.members[index]
}

func (value *serializedValue) asString() (result string) {
	return value.stringValue
}

func (value *serializedValue) asInt8() (result int8) {
	return int8(value.intValue)
}

func (value *serializedValue) asInt32() (result int32) {
	return int32(value.intValue)
}

func (value *serializedValue) asInt64() (result int64) {
	return value.intValue
}
