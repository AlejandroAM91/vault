package vault

import (
	"bytes"
	"crypto/rand"
	"encoding/binary"
	"hash/crc64"
	"io"
)

const (
	// MagicNumber is the magic number to identify an encrypted vault file
	MagicNumber uint32 = 0x17927791
	// VersionMajor is the current encrypted vault file version major number accepted by this library
	VersionMajor uint16 = 1
	// VersionMinor is the current encrypted vault file version ninor number accepted by this library
	VersionMinor uint16 = 0
)

// Version represents encrypted vault file version
//
// Version number is based in semver
type Version struct {
	// Major indicates non backward compatible versions
	Major uint16
	// Minor indicates backward compatible versions
	Minor uint16
}

// Header represents encrypted vault file header
type Header struct {
	Magic    uint32
	Version  Version
	Salt     [32]byte
	Iv       [16]byte
	Checksum uint64
}

// NewHeader creates an empty file header
func NewHeader() (*Header, error) {
	header := Header{
		Magic: MagicNumber,
		Version: Version{
			Major: VersionMajor,
			Minor: VersionMinor,
		},
	}

	if _, err := io.ReadFull(rand.Reader, header.Salt[:]); err != nil {
		return nil, err
	}

	if _, err := io.ReadFull(rand.Reader, header.Iv[:]); err != nil {
		return nil, err
	}

	if err := header.UpdateChecksum(); err != nil {
		return nil, err
	}

	return &header, nil
}

// Check checks if the headers values are correct and compatible with the library
func (h Header) Check() bool {
	if h.Magic != MagicNumber {
		return false
	}

	if h.Version.Major != VersionMajor {
		return false
	}

	if checksum, err := h.checksum(); h.Checksum != checksum || err != nil {
		return false
	}

	return true
}

// Read reads file header from r
func (h *Header) Read(r io.Reader) error {
	return binary.Read(r, binary.BigEndian, h)
}

// Write writes file header into w
func (h Header) Write(w io.Writer) error {
	return binary.Write(w, binary.BigEndian, h)
}

// UpdateChecksum updates the checksum value
func (h *Header) UpdateChecksum() error {
	var err error
	h.Checksum, err = h.checksum()
	return err
}

func (h Header) checksum() (uint64, error) {
	b := new(bytes.Buffer)
	if err := binary.Write(b, binary.BigEndian, h.Magic); err != nil {
		return 0, err
	}

	if err := binary.Write(b, binary.BigEndian, h.Version); err != nil {
		return 0, err
	}

	return crc64.Checksum(b.Bytes(), crc64.MakeTable(crc64.ISO)), nil
}
