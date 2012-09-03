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
	"encoding/binary"
	"errors"
	"fmt"
	"os"
	"strings"
)

const (
	MPQ_UNCOMPRESSED byte = 0x00
	MPQ_BZIP2        byte = 0x10
)

type Mpq struct {
	IsLoaded bool
	inFile   *os.File

	ArchiveOffset uint32
	Header        Header

	HasUserData bool
	UserData    *UserData

	Files        map[string]*File
	HashEntries  []*HashEntry
	BlockEntries []*BlockEntry
}

func LoadMpq(filePath string) (mpq *Mpq, err error) {
	mpq = new(Mpq)
	err = mpq.Load(filePath)
	if err != nil {
		return nil, err
	}
	return mpq, nil
}

func (mpq *Mpq) Load(filePath string) (err error) {
	if mpq.IsLoaded {
		mpq.Close()
	}

	mpq.Files = make(map[string]*File, 0)
	mpq.inFile, err = os.Open(filePath)
	if err != nil {
		return err
	}

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

	mpq.IsLoaded = true
	return
}

func (mpq *Mpq) Close() {
	if mpq.inFile != nil {
		mpq.inFile.Close()
		mpq.inFile = nil
	}
}

func (mpq *Mpq) GetHashEntry(filename string) (entry *HashEntry, err error) {
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

func (mpq *Mpq) GetFile(filename string) (file *File, err error) {
	file = mpq.Files[filename]
	if file != nil {
		return
	}

	hash, err := mpq.GetHashEntry(filename)
	if err != nil {
		return
	}

	block := mpq.BlockEntries[hash.BlockIndex]
	if block == nil {
		return nil, fmt.Errorf("Unable to find block entry %v for file %v",
			hash.BlockIndex, filename)
	}

	file = newFile(filename, hash, block, mpq)
	mpq.Files[filename] = file

	return
}

func (mpq *Mpq) readHeader() error {
	buf := make([]byte, 64)

	_, _ = mpq.inFile.Read(buf[:4])
	if !strings.HasPrefix(string(buf[:3]), "MPQ") {
		mpq.Close()
		return errors.New("File is not an MPQ")
	}

	if buf[3] == 0x1b {
		// This is the user data portion of the file

		// Header Size
		_, _ = mpq.inFile.Read(buf[:4])
		_ = binary.LittleEndian.Uint32(buf[:4])

		// Archive Size
		_, _ = mpq.inFile.Read(buf[:4])
		userArchiveSize := binary.LittleEndian.Uint32(buf[:4])

		mpq.inFile.Seek(0, 0)

		user_buf := make([]byte, userArchiveSize)

		mpq.inFile.Read(user_buf)
		mpq.UserData = readUserData(user_buf)
		mpq.HasUserData = true
		mpq.inFile.Seek(4, 1)
	}

	_, _ = mpq.inFile.Read(buf[:4])
	mpq.Header.HeaderSize = binary.LittleEndian.Uint32(buf[:4])

	buf = make([]byte, mpq.Header.HeaderSize-4)
	mpq.inFile.Read(buf)

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
	hashEntries := mpq.Header.HashTableEntries

	mpq.HashEntries = make([]*HashEntry, hashEntries)

	mpq.inFile.Seek(
		int64(mpq.ArchiveOffset+mpq.Header.HashTableOffset), 0)

	// Each entry is the size of 4x uint32, giving 16 bytes.
	// Create a buffer and read the entire hash table from
	// the file, then decrypt it.
	buffer := make([]byte, hashEntries*16)
	mpq.inFile.Read(buffer)

	encryptor := newBlockEncryptor("(hash table)", 0x300)
	encryptor.decrypt(&buffer)

	offset := 0
	for idx := uint32(0); idx < hashEntries; idx++ {
		entry := newHashEntry(buffer[offset : offset+16])
		mpq.HashEntries[idx] = entry
		offset += 16
	}

	return
}

func (mpq *Mpq) readBlockTable() (err error) {
	blockEntries := mpq.Header.BlockTableEntries

	mpq.BlockEntries = make([]*BlockEntry, blockEntries)

	mpq.inFile.Seek(int64(mpq.ArchiveOffset+mpq.Header.BlockTableOffset), 0)

	// Each entry is the size of 4x uint32, giving 16 bytes.
	// Create a buffer and read the entire hash table from
	// the file, then decrypt it.
	buffer := make([]byte, blockEntries*16)
	mpq.inFile.Read(buffer)

	encryptor := newBlockEncryptor("(block table)", 0x300)
	encryptor.decrypt(&buffer)

	offset := 0
	for idx := uint32(0); idx < blockEntries; idx++ {
		entry := newBlockEntry(buffer[offset : offset+16])
		mpq.BlockEntries[idx] = entry
		offset += 16
	}

	return
}

func (mpq *Mpq) readFiles() (err error) {
	listfile, err := mpq.GetFile("(listfile)")
	if err != nil {
		return
	}

	outBuffer := make([]byte, listfile.FileSize)
	outReader := listfile.GetReader()
	_, err = outReader.Read(outBuffer)
	if err != nil {
		return
	}

	listFile := string(outBuffer)
	files := strings.Split(listFile, "\r\n")
	for _, filename := range files {
		hash, _ := mpq.GetHashEntry(filename)
		if hash == nil {
			continue
		}
		file, err := mpq.GetFile(filename)
		if err != nil {
			return fmt.Errorf("Unable to find file '%v' in MPQ: %v", filename, err.Error())
		}
		mpq.Files[filename] = file
	}

	return
}
