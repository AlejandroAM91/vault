package vault

import (
	"bytes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/binary"
	"errors"
	"io"
)

// Reader implements an encrypted vault file reader.
type Reader struct {
	buffer bytes.Buffer
	header *Header
	key    []byte
	r      io.Reader
	s      cipher.Stream
}

// NewReader returns a new Reader reading an encrypted vault file from r.
func NewReader(r io.Reader, pass []byte) (*Reader, error) {
	reader := Reader{r: r}
	if err := reader.init(pass); err != nil {
		return nil, err
	}
	return &reader, nil
}

func (r *Reader) Read(dst []byte) (int, error) {
	end := false
	for !end && r.buffer.Len() < len(dst) {
		err := r.decryptBlock()
		if err == io.EOF {
			end = true
		} else if err != nil {
			return 0, err
		}
	}

	n, err := r.buffer.Read(dst)
	if err != nil {
		return 0, err
	}

	if end && n < len(dst) {
		return n, io.EOF
	}

	return n, nil
}

func (r *Reader) decryptBlock() error {
	ibuf := make([]byte, lenSize+blockSize+macSize)
	bbuf := ibuf[:lenSize+blockSize]

	// Reads block from reader
	n, ferr := r.r.Read(ibuf)
	if n > 0 {
		// Decrypt block
		r.s.XORKeyStream(bbuf, bbuf)

		// Calculates block MAC
		mac := hmac.New(sha256.New, r.key)
		if _, err := mac.Write(bbuf); err != nil {
			return err
		}

		// Checks block MAC
		if !hmac.Equal(ibuf[lenSize+blockSize:], mac.Sum(nil)) {
			return errors.New("MAC doesn´t match")
		}

		// Writes block to buffer
		l := int(binary.BigEndian.Uint16(bbuf[:lenSize]))
		if _, err := r.buffer.Write(bbuf[lenSize : lenSize+l]); err != nil {
			return err
		}
	}
	return ferr
}

func (r *Reader) init(pass []byte) error {
	var err error
	// Reads and checks the header
	r.header = new(Header)
	if err = r.header.Read(r.r); err != nil {
		return err
	}

	if !r.header.Check() {
		return errors.New("Error on header validation")
	}

	// Generates the key
	r.key = generateKey(pass, r.header.Salt[:])

	// Creates the decryption stream
	if r.s, err = createStream(r.key, r.header.Iv[:]); err != nil {
		return err
	}

	return nil
}
