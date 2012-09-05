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

// Package mpq provides the ability to read MPQ files
package mpq

import (
	"compress/bzip2"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"strings"
)

var EOF = errors.New("EOF")

const (
	COMPRESS_NONE  byte = 0x00
	COMPRESS_BZIP2 byte = 0x10
)

type Mpq struct {
	reader io.ReadSeeker

	archiveOffset uint32
	header        header

	hasUserData bool
	userData    *userData

	files        map[string]*File
	hashEntries  []*hashEntry
	blockEntries []*blockEntry

	file          *File
	fileReader    io.Reader
	fileBytesRead int
}

func NewMpq(reader io.ReadSeeker) (mpq *Mpq, err error) {
	mpq = new(Mpq)
	err = mpq.readHeaders(reader)
	if err != nil {
		return nil, err
	}
	return mpq, nil
}

func (mpq *Mpq) readHeaders(reader io.ReadSeeker) (err error) {
	mpq.files = make(map[string]*File)

	mpq.reader = reader
	mpq.reader.Seek(0, 0)

	err = mpq.readHeader()
	if err != nil {
		return err
	}

	err = mpq.readHashTable()
	if err != nil {
		return err
	}

	err = mpq.readBlockTable()
	if err != nil {
		return err
	}

	err = mpq.readFiles()
	if err != nil {
		fmt.Printf("qwefasdf")
		return err
	}

	return
}

func (mpq *Mpq) File(filename string) (file *File, err error) {
	// First see if the file is already in the map
	file, found := mpq.files[filename]
	if !found {
		// If it's not, try to load it
		fileHash, hashErr := mpq.getHashEntry(filename)
		if hashErr != nil {
			err = fmt.Errorf("Unable to find file: %v",
				filename)
			return
		}
		fileBlock := mpq.blockEntries[fileHash.blockIndex]

		file = newFile(filename, fileHash, fileBlock)
		mpq.files[filename] = file
	}

	mpq.reader.Seek(int64(mpq.archiveOffset+
		file.block.filePosition), 0)
	compType := make([]byte, 1)
	mpq.reader.Read(compType)
	file.compressionType = compType[0]

	switch file.compressionType {
	case COMPRESS_BZIP2:
		mpq.fileReader = bzip2.NewReader(mpq.reader)
		break
	case COMPRESS_NONE:
		fallthrough
	default:
		mpq.fileReader = mpq.reader
		break
	}

	mpq.file = file
	// Reset the number of bytes read from the file
	mpq.fileBytesRead = 0

	return
}

func (mpq *Mpq) Files() (files map[string]*File) {
	return mpq.files
}

func (mpq *Mpq) Read(p []byte) (n int, err error) {
	// Don't allow the read to go past the end of the
	// file, even if there's more data in the MPQ.
	// This will simulate reading a file to the end as
	// if it were reading it off the file system.
	bytesLeft := int(mpq.file.FileSize) - mpq.fileBytesRead
	readBuffer := p
	if bytesLeft <= 0 {
		return 0, EOF
	}
	if len(p) > bytesLeft {
		readBuffer = p[:bytesLeft]
	}
	n, err = mpq.fileReader.Read(readBuffer)
	mpq.fileBytesRead += n
	if err != nil {
		return n, err
	}

	if n == bytesLeft {
		return n, EOF
	}

	return
}

func (mpq *Mpq) getHashEntry(filename string) (entry *hashEntry, err error) {
	hashA := hashString(filename, 0x100)
	hashB := hashString(filename, 0x200)

	for _, entry := range mpq.hashEntries {
		if entry.filePathHashA == hashA &&
			entry.filePathHashB == hashB {
			return entry, nil
		}
	}

	return nil, fmt.Errorf("Could not find hash entry: %v", filename)
}

func (mpq *Mpq) readHeader() (err error) {
	buf := make([]byte, 64)

	_, err = mpq.reader.Read(buf[:4])
	if err != nil {
		return fmt.Errorf("Could not read MPQ header")
	}
	if !strings.HasPrefix(string(buf[:3]), "MPQ") {
		return fmt.Errorf("File is not an MPQ")
	}

	if buf[3] == 0x1b {
		// This is the user data portion of the file

		// Header Size
		_, _ = mpq.reader.Read(buf[:4])
		_ = binary.LittleEndian.Uint32(buf[:4])

		// Archive Size
		_, _ = mpq.reader.Read(buf[:4])
		userArchiveSize := binary.LittleEndian.Uint32(buf[:4])

		mpq.reader.Seek(0, 0)

		user_buf := make([]byte, userArchiveSize)

		mpq.reader.Read(user_buf)
		mpq.userData = readUserData(user_buf)
		mpq.hasUserData = true
		mpq.reader.Seek(4, 1)
	}

	_, _ = mpq.reader.Read(buf[:4])
	mpq.header.headerSize = binary.LittleEndian.Uint32(buf[:4])

	buf = make([]byte, mpq.header.headerSize-4)
	mpq.reader.Read(buf)

	mpq.header.archiveSize = binary.LittleEndian.Uint32(buf[:4])
	mpq.header.formatVersion = binary.LittleEndian.Uint16(buf[0x04 : 0x04+2])
	mpq.header.blockSize = binary.LittleEndian.Uint16(buf[0x06 : 0x06+2])
	mpq.header.hashTableOffset = binary.LittleEndian.Uint32(buf[0x08 : 0x08+4])
	mpq.header.blockTableOffset = binary.LittleEndian.Uint32(buf[0x0c : 0x0c+4])
	mpq.header.hashTableEntries = binary.LittleEndian.Uint32(buf[0x10 : 0x10+4])
	mpq.header.blockTableEntries = binary.LittleEndian.Uint32(buf[0x14 : 0x14+4])
	mpq.header.extendedBlockTableOffset = binary.LittleEndian.Uint64(buf[0x18 : 0x18+8])
	mpq.header.hashTableOffsetHigh = binary.LittleEndian.Uint16(buf[0x20 : 0x20+2])
	mpq.header.blockTableOffsetHigh = binary.LittleEndian.Uint16(buf[0x22 : 0x22+2])

	mpq.archiveOffset = 0x00
	if mpq.hasUserData {
		mpq.archiveOffset = mpq.userData.header.archiveOffset
	}

	return nil
}

func (mpq *Mpq) readHashTable() (err error) {
	hashEntries := mpq.header.hashTableEntries

	mpq.hashEntries = make([]*hashEntry, hashEntries)

	mpq.reader.Seek(
		int64(mpq.archiveOffset+mpq.header.hashTableOffset), 0)

	// Each entry is the size of 4x uint32, giving 16 bytes.
	// Create a buffer and read the entire hash table from
	// the file, then decrypt it.
	buffer := make([]byte, hashEntries*16)
	mpq.reader.Read(buffer)

	encryptor := newBlockEncryptor("(hash table)", 0x300)
	encryptor.decrypt(&buffer)

	offset := 0
	for idx := uint32(0); idx < hashEntries; idx++ {
		entry := newHashEntry(buffer[offset : offset+16])
		mpq.hashEntries[idx] = entry
		offset += 16
	}

	return nil
}

func (mpq *Mpq) readBlockTable() (err error) {
	blockEntries := mpq.header.blockTableEntries

	mpq.blockEntries = make([]*blockEntry, blockEntries)

	mpq.reader.Seek(int64(mpq.archiveOffset+mpq.header.blockTableOffset), 0)

	// Each entry is the size of 4x uint32, giving 16 bytes.
	// Create a buffer and read the entire hash table from
	// the file, then decrypt it.
	buffer := make([]byte, blockEntries*16)
	mpq.reader.Read(buffer)

	encryptor := newBlockEncryptor("(block table)", 0x300)
	encryptor.decrypt(&buffer)

	offset := 0
	for idx := uint32(0); idx < blockEntries; idx++ {
		entry := newBlockEntry(buffer[offset : offset+16])
		mpq.blockEntries[idx] = entry
		offset += 16
	}

	return nil
}

func (mpq *Mpq) readFiles() (err error) {
	listfile, err := mpq.File("(listfile)")
	if err != nil {
		return
	}

	outBuffer := make([]byte, listfile.FileSize)
	read, err := mpq.Read(outBuffer)
	if read <= 0 && err != nil {
		return
	}

	listFile := string(outBuffer)
	files := strings.Split(listFile, "\r\n")

	for _, filename := range files {
		file, err := mpq.File(filename)
		if err != nil {
			continue
		}

		mpq.files[filename] = file
	}

	return nil
}
