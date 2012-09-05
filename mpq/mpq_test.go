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
	"math"
	"os"
	"testing"
)

func TestLoadingNonMpqFile(t *testing.T) {
	reader, err := os.Open("testdata/not_a_replay.SC2Replay")
	if err != nil {
		t.Errorf("Could not load test data: %v", err.Error())
		return
	}

	mpq, err := NewMpq(reader)
	if err == nil {
		t.Error("Error was not returned for a non-MPQ file")
	}
	if mpq != nil {
		t.Error("MPQ was returned when loading a non-MPQ file")
	}
}

func TestUserDataHeaderFlag(t *testing.T) {
	reader, err := os.Open("testdata/replay1.SC2Replay")
	if err != nil {
		t.Errorf("Could not load test data: %v", err.Error())
		return
	}

	mpq, err := NewMpq(reader)
	if err != nil {
		t.Errorf("Error loading MPQ: %v", err.Error())
		return
	}
	if !mpq.hasUserData {
		t.Error("User data flag was not set")
	}
	if mpq.userData == nil {
		t.Error("User data was not loaded")
	}
}

func TestUserDataHeader(t *testing.T) {
	reader, err := os.Open("testdata/replay1.SC2Replay")
	if err != nil {
		t.Errorf("Could not load test data: %v", err.Error())
		return
	}

	mpq, err := NewMpq(reader)
	if err != nil {
		t.Errorf("Error loading MPQ: %v", err.Error())
		return
	}
	if mpq.userData.header.maxUserDataSize != 512 {
		t.Errorf("MaxUserDataSize = %v (Expected: 512)",
			mpq.userData.header.maxUserDataSize)
	}
	if mpq.userData.header.archiveOffset != 1024 {
		t.Errorf("ArchiveOffset = %v (Expected: 1024)",
			mpq.userData.header.archiveOffset)
	}
	if mpq.userData.header.userDataSize != 60 {
		t.Errorf("UserDataSize = %v (Expected: 60)",
			mpq.userData.header.userDataSize)
	}
}

func TestReadMpqHeader(t *testing.T) {
	reader, err := os.Open("testdata/replay1.SC2Replay")
	if err != nil {
		t.Errorf("Could not load test data: %v", err.Error())
		return
	}

	mpq, err := NewMpq(reader)
	if err != nil {
		t.Errorf("Error loading MPQ: %v", err.Error())
		return
	}
	if mpq.header.headerSize != 44 {
		t.Errorf("HeaderSize = %v (Expected: 44)",
			mpq.header.headerSize)
	}
	if mpq.header.archiveSize != 109012 {
		t.Errorf("ArchiveSize = %v (Expected: 109012)",
			mpq.header.archiveSize)
	}
	if mpq.header.formatVersion != 0x01 {
		t.Errorf("FormatVersion = %#v (Expected: 0x01)",
			mpq.header.formatVersion)
	}
	if mpq.header.blockSize != 0x03 {
		t.Errorf("BlockSize = %#v (Expected: 0x03)",
			mpq.header.blockSize)
	}
	if mpq.header.hashTableOffset != 108596 {
		t.Errorf("HashTableOffset = %v (Expected: 108596)",
			mpq.header.hashTableOffset)
	}
	if mpq.header.blockTableOffset != 108852 {
		t.Errorf("BlockTableOffset = %v (Expected: 108852)",
			mpq.header.blockTableOffset)
	}
	if mpq.header.hashTableEntries != 16 {
		t.Errorf("HashTableEntries = %v (Expected: 16)",
			mpq.header.hashTableEntries)
	}
	if mpq.header.blockTableEntries != 10 {
		t.Errorf("BlockTableEntries = %v (Expected: 10)",
			mpq.header.blockTableEntries)
	}
	if mpq.header.extendedBlockTableOffset != 0 {
		t.Errorf("ExtendedBlockTableOffset = %v (Expected: 0)",
			mpq.header.extendedBlockTableOffset)
	}
	if mpq.header.hashTableOffsetHigh != 0 {
		t.Errorf("HashTableOffsetHigh = %v (Expected: 0)",
			mpq.header.hashTableOffsetHigh)
	}
	if mpq.header.blockTableOffsetHigh != 0 {
		t.Errorf("BlockTableOffsetHigh = %v (Expected: 0)",
			mpq.header.blockTableOffsetHigh)
	}
}

var expectedFiles = []struct {
	filename string
}{
	{"replay.attributes.events"},
	{"replay.details"},
	{"replay.game.events"},
	{"replay.initData"},
	{"replay.load.info"},
	{"replay.message.events"},
	{"replay.smartcam.events"},
	{"replay.sync.events"},
}

func TestReadFiles(t *testing.T) {
	reader, err := os.Open("testdata/replay1.SC2Replay")
	if err != nil {
		t.Errorf("Could not load test data: %v", err.Error())
		return
	}

	mpq, err := NewMpq(reader)
	if err != nil {
		t.Errorf("Error loading MPQ: %v", err.Error())
		return
	}

	if len(mpq.Files()) < 1 {
		t.Error("No files were found!")
	} else {
		for _, f := range expectedFiles {
			_, err := mpq.File(f.filename)
			if err != nil {
				t.Errorf("Could not find file %v in MPQ.",
					f.filename)
			}
		}
	}
}

func TestReadBeyondFile(t *testing.T) {
	reader, err := os.Open("testdata/replay1.SC2Replay")
	if err != nil {
		t.Errorf("Could not load test data: %v", err.Error())
		return
	}

	// Read the listfile with a buffer exactly the size
	// of the file
	mpq, _ := NewMpq(reader)
	file, err := mpq.File("(listfile)")
	if err != nil {
		t.Errorf("Could not load list file: %v", err.Error())
		return
	}

	// Read the whole file
	buffer := make([]byte, file.FileSize)
	read, err := mpq.Read(buffer)
	if read != int(file.FileSize) {
		t.Errorf("Read the incorrect number of bytes. Expected: %v Actual: %v", file.FileSize, read)
	}
	if err == nil || err != EOF {
		t.Errorf("Read to end of file did not return EOF error.")
	}

	// Read the file in two reads
	mpq.File("(listfile)")
	partialSize := int(math.Ceil(float64(file.FileSize / 2)))
	buffer = make([]byte, partialSize)
	read, err = mpq.Read(buffer)
	if read != partialSize {
		t.Errorf("Partial read of the wrong size. Expected: %v Actual: %v", partialSize, read)
	}
	if err != nil {
		t.Errorf("Partial read returned an unexpected error: %v", err.Error())
	}

	finalSize := int(file.FileSize) - partialSize
	buffer = make([]byte, finalSize)
	read, err = mpq.Read(buffer)
	if read != finalSize {
		t.Error("Partial read of the wrong size. Expected: %v Actual: %v", finalSize, read)
	}
	if err == nil || err != EOF {
		t.Errorf("Partial read expected an EOF error but did not receive one.")
	}

	// Read past end of file
	mpq.File("(listfile)")
	buffer = make([]byte, file.FileSize+10)

	read, err = mpq.Read(buffer)
	if read != int(file.FileSize) {
		t.Errorf("Expected read of size %v but got %v instead.", file.FileSize, read)
	}
	if err == nil || err != EOF {
		t.Errorf("Did not get expected EOF error for read beyond file.")
	}
}
