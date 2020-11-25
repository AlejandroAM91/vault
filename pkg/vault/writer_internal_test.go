package vault

import (
	"bytes"
	"encoding/binary"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type WriterInternalTestSuite struct {
	suite.Suite
	pass    []byte
	bheader []byte
	vheader Header
}

func TestWriterInternalTestSuite(t *testing.T) {
	suite.Run(t, new(WriterInternalTestSuite))
}

func (s *WriterInternalTestSuite) SetupTest() {
	s.pass = []byte("secret pass")
	s.vheader = Header{
		Magic: MagicNumber,
		Version: Version{
			Major: VersionMajor,
			Minor: VersionMinor,
		},
	}

	buf := new(bytes.Buffer)
	if err := binary.Write(buf, binary.BigEndian, s.vheader); err != nil {
		panic(err)
	}
	s.bheader = buf.Bytes()
}

func (s *WriterInternalTestSuite) TestNewWriter() {
	b := &bytes.Buffer{}

	result, err := NewWriter(b, s.pass)
	if assert.Nil(s.T(), err) {
		assert.Equal(s.T(), s.vheader.Magic, result.header.Magic)
		assert.Equal(s.T(), s.vheader.Version, result.header.Version)
		assert.Equal(s.T(), s.bheader[:8], b.Bytes()[:8]) // First 8 bytes should match (Magic, Version)
	}
}
