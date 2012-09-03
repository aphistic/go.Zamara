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
	"testing"
)

func TestErrorOnFileNotFound(t *testing.T) {
	mpq, err := LoadMpq("testdata/does_not_exist.SC2Replay")
	if err == nil {
		t.Error("Error was not set on file not found.")
	}
	if mpq != nil {
		t.Error("MPQ was returned when loading failed.")
	}
}

func TestLoadedFlag(t *testing.T) {
	mpq, err := LoadMpq("testdata/replay1.SC2Replay")
	if mpq == nil {
		t.Error("MPQ was not loaded")
	} else {
		if !mpq.IsLoaded {
			t.Error("Mpq.IsLoaded is false")
		}
	}
	if err != nil {
		t.Error("Error was returned on successful load")
	}
}

func TestLoadingNonMpqFile(t *testing.T) {
	mpq, err := LoadMpq("testdata/not_a_replay.SC2Replay")
	if err == nil {
		t.Error("Error was not returned for a non-MPQ file")
	}
	if mpq != nil {
		t.Error("MPQ was returned when loading a non-MPQ file")
	}
}

func TestUserDataHeaderFlag(t *testing.T) {
	mpq, err := LoadMpq("testdata/replay1.SC2Replay")
	if err != nil {
		t.Errorf("Error loading MPQ: %v", err.Error())
		return
	}
	if !mpq.HasUserData {
		t.Error("User data flag was not set")
	}
	if mpq.UserData == nil {
		t.Error("User data was not loaded")
	}
}

func TestUserDataHeader(t *testing.T) {
	mpq, err := LoadMpq("testdata/replay1.SC2Replay")
	if err != nil {
		t.Errorf("Error loading MPQ: %v", err.Error())
		return
	}
	if mpq.UserData.Header.MaxUserDataSize != 512 {
		t.Errorf("MaxUserDataSize = %v (Expected: 512)",
			mpq.UserData.Header.MaxUserDataSize)
	}
	if mpq.UserData.Header.ArchiveOffset != 1024 {
		t.Errorf("ArchiveOffset = %v (Expected: 1024)",
			mpq.UserData.Header.ArchiveOffset)
	}
	if mpq.UserData.Header.UserDataSize != 60 {
		t.Errorf("UserDataSize = %v (Expected: 60)",
			mpq.UserData.Header.UserDataSize)
	}
}

func TestReadMpqHeader(t *testing.T) {
	mpq, err := LoadMpq("testdata/replay1.SC2Replay")
	if err != nil {
		t.Errorf("Error loading MPQ: %v", err.Error())
		return
	}
	if mpq.Header.HeaderSize != 44 {
		t.Errorf("HeaderSize = %v (Expected: 44)",
			mpq.Header.HeaderSize)
	}
	if mpq.Header.ArchiveSize != 109012 {
		t.Errorf("ArchiveSize = %v (Expected: 109012)",
			mpq.Header.ArchiveSize)
	}
	if mpq.Header.FormatVersion != 0x01 {
		t.Errorf("FormatVersion = %#v (Expected: 0x01)",
			mpq.Header.FormatVersion)
	}
	if mpq.Header.BlockSize != 0x03 {
		t.Errorf("BlockSize = %#v (Expected: 0x03)",
			mpq.Header.BlockSize)
	}
	if mpq.Header.HashTableOffset != 108596 {
		t.Errorf("HashTableOffset = %v (Expected: 108596)",
			mpq.Header.HashTableOffset)
	}
	if mpq.Header.BlockTableOffset != 108852 {
		t.Errorf("BlockTableOffset = %v (Expected: 108852)",
			mpq.Header.BlockTableOffset)
	}
	if mpq.Header.HashTableEntries != 16 {
		t.Errorf("HashTableEntries = %v (Expected: 16)",
			mpq.Header.HashTableEntries)
	}
	if mpq.Header.BlockTableEntries != 10 {
		t.Errorf("BlockTableEntries = %v (Expected: 10)",
			mpq.Header.BlockTableEntries)
	}
	if mpq.Header.ExtendedBlockTableOffset != 0 {
		t.Errorf("ExtendedBlockTableOffset = %v (Expected: 0)",
			mpq.Header.ExtendedBlockTableOffset)
	}
	if mpq.Header.HashTableOffsetHigh != 0 {
		t.Errorf("HashTableOffsetHigh = %v (Expected: 0)",
			mpq.Header.HashTableOffsetHigh)
	}
	if mpq.Header.BlockTableOffsetHigh != 0 {
		t.Errorf("BlockTableOffsetHigh = %v (Expected: 0)",
			mpq.Header.BlockTableOffsetHigh)
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
	mpq, err := LoadMpq("testdata/replay1.SC2Replay")
	if err != nil {
		t.Errorf("Error loading MPQ: %v", err.Error())
		return
	}

	if len(mpq.Files) < 1 {
		t.Error("No files were found!")
	} else {
		for _, f := range expectedFiles {
			_, err := mpq.GetFile(f.filename)
			if err != nil {
				t.Errorf("Could not find file %v in MPQ.",
					f.filename)
			}
		}
	}
}
