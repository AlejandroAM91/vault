package vault

import (
	"bytes"
	"encoding/binary"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type ReaderInternalTestSuite struct {
	suite.Suite
	pass    []byte
	bheader []byte
	vheader Header
}

func TestReaderInternalTestSuite(t *testing.T) {
	suite.Run(t, new(ReaderInternalTestSuite))
}

func (s *ReaderInternalTestSuite) SetupTest() {
	s.pass = []byte("secret pass")
	s.vheader = Header{
		Magic: MagicNumber,
		Version: Version{
			Major: VersionMajor,
			Minor: VersionMinor,
		},
	}
	s.vheader.UpdateChecksum()

	buf := new(bytes.Buffer)
	if err := binary.Write(buf, binary.BigEndian, s.vheader); err != nil {
		panic(err)
	}
	s.bheader = buf.Bytes()
}

func (s *ReaderInternalTestSuite) TestNewReader() {
	b := bytes.NewBuffer(s.bheader)

	result, err := NewReader(b, s.pass)
	if assert.Nil(s.T(), err) {
		assert.Equal(s.T(), s.vheader.Magic, result.header.Magic)
		assert.Equal(s.T(), s.vheader.Version, result.header.Version)
	}
}
