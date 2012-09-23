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
)

type attributeFile struct {
	attributes []*replayAttribute
}

func newAttributeFile(data []byte) (file *attributeFile, err error) {
	file = new(attributeFile)
	err = file.load(data)
	if err != nil {
		return nil, err
	}

	return file, nil
}

func (file *attributeFile) load(data []byte) (err error) {
	// Find the beginning of the actual data
	offset := int64(0)
	for data[offset] == 0x00 {
		offset++
	}

	count := int64(binary.LittleEndian.Uint32(data[offset : offset+4]))
	offset += 4

	file.attributes = make([]*replayAttribute, count)
	for idx := int64(0); idx < count; idx++ {
		attr, size, err := newReplayAttribute(data[offset:])
		if err != nil {
			return err
		}
		offset += size
		file.attributes[idx] = attr
	}

	return
}
