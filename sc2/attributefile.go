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
