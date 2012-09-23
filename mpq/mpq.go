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
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"strings"
)

var EOF = errors.New("EOF")

const (
	CompressNone  byte = 0x00
	CompressBzip2 byte = 0x10
)

type Mpq struct {
	XMLName xml.Name `xml:"mpq"`

	reader io.ReadSeeker

	ArchiveOffset uint32 `xml:"archiveOffset"`
	Header        Header `xml:"header"`

	HasUserData bool      `xml:"-"`
	UserData    *UserData `xml:"userData"`

	files        map[string]*File
	HashEntries  []*HashEntry  `xml:"hashEntries>hashEntry"`
	BlockEntries []*BlockEntry `xml:"blockEntries>blockEntry"`

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
		fileBlock := mpq.BlockEntries[fileHash.BlockIndex]

		file = newFile(filename, fileHash, fileBlock)
		mpq.files[filename] = file
	}

	_, err = mpq.reader.Seek(int64(mpq.ArchiveOffset+
		file.block.FilePosition), 0)
	if err != nil {
		return
	}
	compType := make([]byte, 1)
	_, err = mpq.reader.Read(compType)
	if err != nil {
		return
	}

	file.compressionType = compType[0]
	switch file.compressionType {
	case CompressBzip2:
		mpq.fileReader = bzip2.NewReader(mpq.reader)
		break
	case CompressNone:
		// Explicitly found that this doesn't have
		// compression, so don't go back to the
		// original file position as in the default.
		mpq.fileReader = mpq.reader
		break
	default:
		// Don't know this compression type, so just
		// reset back one byte and read from there as
		// if this file doesn't have a header byte.  This
		// can happen if the file flags say the file is
		// compressed but the header byte doesn't exist.
		_, err = mpq.reader.Seek(int64(mpq.ArchiveOffset+
			file.block.FilePosition), 0)
		if err != nil {
			return
		}
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

func (mpq *Mpq) getHashEntry(filename string) (entry *HashEntry, err error) {
	hashA := hashString(filename, 0x100)
	hashB := hashString(filename, 0x200)

	for _, entry := range mpq.HashEntries {
		if entry.FilePathHashA == hashA &&
			entry.FilePathHashB == hashB {
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
		mpq.UserData = readUserData(user_buf)
		mpq.HasUserData = true
		mpq.reader.Seek(4, 1)
	}

	_, _ = mpq.reader.Read(buf[:4])
	mpq.Header.HeaderSize = binary.LittleEndian.Uint32(buf[:4])

	buf = make([]byte, mpq.Header.HeaderSize-4)
	mpq.reader.Read(buf)

	mpq.Header.ArchiveSize = binary.LittleEndian.Uint32(buf[:4])
	mpq.Header.FormatVersion = binary.LittleEndian.Uint16(buf[0x04 : 0x04+2])
	mpq.Header.BlockSize = binary.LittleEndian.Uint16(buf[0x06 : 0x06+2])
	mpq.Header.HashTableOffset = binary.LittleEndian.Uint32(buf[0x08 : 0x08+4])
	mpq.Header.BlockTableOffset = binary.LittleEndian.Uint32(buf[0x0c : 0x0c+4])
	mpq.Header.HashTableEntries = binary.LittleEndian.Uint32(buf[0x10 : 0x10+4])
	mpq.Header.BlockTableEntries = binary.LittleEndian.Uint32(buf[0x14 : 0x14+4])
	mpq.Header.ExtendedBlockTableOffset = binary.LittleEndian.Uint64(buf[0x18 : 0x18+8])
	mpq.Header.HashTableOffsetHigh = binary.LittleEndian.Uint16(buf[0x20 : 0x20+2])
	mpq.Header.BlockTableOffsetHigh = binary.LittleEndian.Uint16(buf[0x22 : 0x22+2])

	mpq.ArchiveOffset = 0x00
	if mpq.HasUserData {
		mpq.ArchiveOffset = mpq.UserData.Header.ArchiveOffset
	}

	return nil
}

func (mpq *Mpq) readHashTable() (err error) {
	HashEntries := mpq.Header.HashTableEntries

	mpq.HashEntries = make([]*HashEntry, HashEntries)

	mpq.reader.Seek(
		int64(mpq.ArchiveOffset+mpq.Header.HashTableOffset), 0)

	// Each entry is the size of 4x uint32, giving 16 bytes.
	// Create a buffer and read the entire hash table from
	// the file, then decrypt it.
	buffer := make([]byte, HashEntries*16)
	mpq.reader.Read(buffer)

	encryptor := newBlockEncryptor("(hash table)", 0x300)
	encryptor.decrypt(&buffer)

	offset := 0
	for idx := uint32(0); idx < HashEntries; idx++ {
		entry := newHashEntry(buffer[offset : offset+16])
		mpq.HashEntries[idx] = entry
		offset += 16
	}

	return nil
}

func (mpq *Mpq) readBlockTable() (err error) {
	BlockEntries := mpq.Header.BlockTableEntries

	mpq.BlockEntries = make([]*BlockEntry, BlockEntries)

	mpq.reader.Seek(int64(mpq.ArchiveOffset+mpq.Header.BlockTableOffset), 0)

	// Each entry is the size of 4x uint32, giving 16 bytes.
	// Create a buffer and read the entire hash table from
	// the file, then decrypt it.
	buffer := make([]byte, BlockEntries*16)
	mpq.reader.Read(buffer)

	encryptor := newBlockEncryptor("(block table)", 0x300)
	encryptor.decrypt(&buffer)

	offset := 0
	for idx := uint32(0); idx < BlockEntries; idx++ {
		entry := newBlockEntry(buffer[offset : offset+16])
		mpq.BlockEntries[idx] = entry
		offset += 16
	}

	return nil
}

func (mpq *Mpq) readFiles() (err error) {
	// Attempt to read the special files just
	// to get them in the file list, since they
	// won't be in the list file
	mpq.File("(attributes)")
	mpq.File("(signature)")
	mpq.File("(user data)")

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
