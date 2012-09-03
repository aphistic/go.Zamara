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

package mpq

import (
	"compress/bzip2"
	"io"
)

type File struct {
	Filename string

	CompressedSize uint32
	FileSize       uint32
	Flags          uint32
	Language       uint16
	Platform       uint16

	compressionType byte
	block           *BlockEntry
	hash            *HashEntry
	mpq             *Mpq
}

func newFile(filename string, hash *HashEntry, block *BlockEntry, mpq *Mpq) (file *File) {
	file = new(File)

	file.Filename = filename
	file.FileSize = block.FileSize
	file.Flags = block.Flags
	file.Language = hash.Language
	file.Platform = hash.Platform

	file.block = block
	file.hash = hash
	file.mpq = mpq

	compressionTypeBuffer := make([]byte, 1)
	mpq.inFile.Seek(file.fileOffset(), 0)
	mpq.inFile.Read(compressionTypeBuffer)
	file.compressionType = compressionTypeBuffer[0]

	return
}

func (file *File) fileOffset() (offset int64) {
	return int64(file.mpq.ArchiveOffset + file.block.FilePosition)
}

func (file *File) GetReader() (reader io.Reader) {
	file.mpq.inFile.Seek(file.fileOffset()+1, 0)
	switch file.compressionType {
	case MPQ_BZIP2:
		reader = bzip2.NewReader(file.mpq.inFile)
		return
	case MPQ_UNCOMPRESSED:
		fallthrough
	default:
		reader = file.mpq.inFile
		return
	}

	return nil
}
