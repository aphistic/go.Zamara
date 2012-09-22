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
	. "launchpad.net/gocheck"
	"math"
	"os"
	"testing"
)

// Hook into gotest
func Test(t *testing.T) { TestingT(t) }

type MpqSuite struct{}

var _ = Suite(&MpqSuite{})

func (s *MpqSuite) TestLoadingNonMpqFile(c *C) {
	reader, err := os.Open("testdata/not_a_replay.SC2Replay")
	c.Assert(err, IsNil)

	mpq, err := NewMpq(reader)
	c.Check(err, NotNil)
	c.Check(mpq, IsNil)
}

func (s *MpqSuite) TestUserDataHeaderFlag(c *C) {
	reader, err := os.Open("testdata/replay1.SC2Replay")
	c.Assert(err, IsNil)

	mpq, err := NewMpq(reader)
	c.Assert(err, IsNil)

	c.Check(mpq.HasUserData, Equals, true)
	c.Check(mpq.UserData, NotNil)
}

func (s *MpqSuite) TestUserDataHeader(c *C) {
	reader, err := os.Open("testdata/replay1.SC2Replay")
	c.Assert(err, IsNil)

	mpq, err := NewMpq(reader)
	c.Assert(err, IsNil)
	c.Check(mpq.UserData.Header.MaxUserDataSize, Equals, uint32(512))
	c.Check(mpq.UserData.Header.UserDataSize, Equals, uint32(60))
}

func (s *MpqSuite) TestReadMpqHeader(c *C) {
	reader, err := os.Open("testdata/replay1.SC2Replay")
	c.Assert(err, IsNil)

	mpq, err := NewMpq(reader)
	c.Assert(err, IsNil)

	c.Check(mpq.Header.HeaderSize, Equals, uint32(44))
	c.Check(mpq.Header.ArchiveSize, Equals, uint32(109012))
	c.Check(mpq.Header.FormatVersion, Equals, uint16(0x01))
	c.Check(mpq.Header.BlockSize, Equals, uint16(0x03))
	c.Check(mpq.Header.HashTableOffset, Equals, uint32(108596))
	c.Check(mpq.Header.BlockTableOffset, Equals, uint32(108852))
	c.Check(mpq.Header.HashTableEntries, Equals, uint32(16))
	c.Check(mpq.Header.BlockTableEntries, Equals, uint32(10))
	c.Check(mpq.Header.ExtendedBlockTableOffset, Equals, uint64(0))
	c.Check(mpq.Header.HashTableOffsetHigh, Equals, uint16(0))
	c.Check(mpq.Header.BlockTableOffsetHigh, Equals, uint16(0))
}

var expectedFiles = []string {
	"(listfile)",
	"(attributes)",
	"replay.attributes.events",
	"replay.details",
	"replay.game.events",
	"replay.initData",
	"replay.load.info",
	"replay.message.events",
	"replay.smartcam.events",
	"replay.sync.events",
}

func (s *MpqSuite) TestReadFiles(c *C) {
	reader, err := os.Open("testdata/replay1.SC2Replay")
	c.Assert(err, IsNil)

	mpq, err := NewMpq(reader)
	c.Assert(err, IsNil)

	c.Check(len(mpq.Files()), Equals, len(expectedFiles))
	for _, filename := range expectedFiles {
		_, err := mpq.File(filename)
		c.Check(err, IsNil)
		if err != nil {
			c.Errorf("Could not find file %v in MPQ.",
				filename)
		}
	}
}

func (s *MpqSuite) TestReadBeyondFile(c *C) {
	reader, err := os.Open("testdata/replay1.SC2Replay")
	c.Assert(err, IsNil)

	// Read the listfile with a buffer exactly the size
	// of the file
	mpq, _ := NewMpq(reader)
	file, err := mpq.File("(listfile)")
	if err != nil {
		c.Fatalf("Could not load list file: %v", err.Error())
	}

	// Read the whole file
	buffer := make([]byte, file.FileSize)
	read, err := mpq.Read(buffer)
	if read != int(file.FileSize) {
		c.Errorf("Read the incorrect number of bytes. Expected: %v Actual: %v", file.FileSize, read)
	}
	if err == nil || err != EOF {
		c.Errorf("Read to end of file did not return EOF error.")
	}

	// Read the file in two reads
	mpq.File("(listfile)")
	partialSize := int(math.Ceil(float64(file.FileSize / 2)))
	buffer = make([]byte, partialSize)
	read, err = mpq.Read(buffer)
	if read != partialSize {
		c.Errorf("Partial read of the wrong size. Expected: %v Actual: %v", partialSize, read)
	}
	if err != nil {
		c.Errorf("Partial read returned an unexpected error: %v", err.Error())
	}

	finalSize := int(file.FileSize) - partialSize
	buffer = make([]byte, finalSize)
	read, err = mpq.Read(buffer)
	if read != finalSize {
		c.Error("Partial read of the wrong size. Expected: %v Actual: %v", finalSize, read)
	}
	if err == nil || err != EOF {
		c.Errorf("Partial read expected an EOF error but did not receive one.")
	}

	// Read past end of file
	mpq.File("(listfile)")
	buffer = make([]byte, file.FileSize+10)

	read, err = mpq.Read(buffer)
	if read != int(file.FileSize) {
		c.Errorf("Expected read of size %v but got %v instead.", file.FileSize, read)
	}
	if err == nil || err != EOF {
		c.Errorf("Did not get expected EOF error for read beyond file.")
	}
}

func (s *MpqSuite) TestSelectFileFromClosedReader(c *C) {
	reader, err := os.Open("testdata/replay1.SC2Replay")
	c.Assert(err, IsNil)

	mpq, err := NewMpq(reader)
	c.Assert(err, IsNil)

	reader.Close()

	_, err = mpq.File("(listfile)")
	if err == nil {
		c.Error("Attempt to select file from closed reader succeeded.")
	}
}

func (s *MpqSuite) TestReadFromClosedReader(c *C) {
	reader, err := os.Open("testdata/replay1.SC2Replay")
	c.Assert(err, IsNil)

	mpq, err := NewMpq(reader)
	c.Assert(err, IsNil)

	mpq.File("(listfile)")

	reader.Close()

	buffer := make([]byte, 8)
	_, err = mpq.Read(buffer)
	if err == nil {
		c.Error("Attempt to read from a closed io.Reader succeeded.")
	}
}
