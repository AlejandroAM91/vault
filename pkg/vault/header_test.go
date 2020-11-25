package vault

import (
	"bytes"
	"encoding/binary"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type HeaderTestSuite struct {
	suite.Suite
	bheader []byte
	vheader Header
}

func TestHeaderTestSuite(t *testing.T) {
	suite.Run(t, new(HeaderTestSuite))
}

func (s *HeaderTestSuite) SetupTest() {
	s.vheader = Header{
		Magic: MagicNumber,
		Version: Version{
			Major: VersionMajor,
			Minor: VersionMinor,
		},
		Checksum: 0xB7FFFE4F3D0A2F80,
	}

	buf := new(bytes.Buffer)
	if err := binary.Write(buf, binary.BigEndian, s.vheader); err != nil {
		panic(err)
	}
	s.bheader = buf.Bytes()
}

func (s *HeaderTestSuite) TestNewHeader() {
	result, err := NewHeader()
	if assert.Nil(s.T(), err) {
		assert.Equal(s.T(), s.vheader.Magic, result.Magic)
		assert.Equal(s.T(), s.vheader.Version, result.Version)
	}
}

func (s *HeaderTestSuite) TestRead() {
	buf := bytes.NewBuffer(s.bheader)
	header := Header{}
	if err := header.Read(buf); assert.Nil(s.T(), err) {
		assert.Equal(s.T(), s.vheader, header)
	}
}

func (s *HeaderTestSuite) TestWrite() {
	buf := new(bytes.Buffer)
	if err := s.vheader.Write(buf); assert.Nil(s.T(), err) {
		assert.Equal(s.T(), s.bheader, buf.Bytes())
	}
}

func (s *HeaderTestSuite) TestCheck() {
	tests := []struct {
		header Header
		result bool
	}{
		{
			header: s.vheader,
			result: true,
		},
		{
			header: Header{Magic: s.vheader.Magic},
			result: false,
		},
		{
			header: Header{Version: s.vheader.Version},
			result: false,
		},
		{
			header: Header{Checksum: s.vheader.Checksum},
			result: false,
		},
		{
			header: Header{
				Magic:   s.vheader.Magic,
				Version: s.vheader.Version,
			},
			result: false,
		},
	}

	for _, test := range tests {
		result := test.header.Check()
		assert.Equal(s.T(), test.result, result)
	}
}

func (s *HeaderTestSuite) TestUpdateChecksum() {
	header := Header{
		Magic:   s.vheader.Magic,
		Version: s.vheader.Version,
	}
	err := header.UpdateChecksum()
	assert.Nil(s.T(), err)

	assert.Equal(s.T(), s.vheader.Checksum, header.Checksum)
}
