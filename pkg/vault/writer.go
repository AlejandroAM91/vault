package vault

import (
	"bytes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/binary"
	"io"
)

// Writer implements an encrypted vault file writer.
type Writer struct {
	buffer bytes.Buffer
	header *Header
	key    []byte
	s      cipher.Stream
	w      io.Writer
}

// NewWriter returns a new Writer writing an encrypted vault file to w.
func NewWriter(w io.Writer, pass []byte) (*Writer, error) {
	writer := Writer{w: w}
	if err := writer.init(pass); err != nil {
		return nil, err
	}
	return &writer, nil
}

// Close write the data if necessary
//
// This function do not close the underlying writer
func (w *Writer) Close() error {
	for w.buffer.Len() > 0 {
		if err := w.encryptBlock(); err != nil {
			return err
		}
	}
	return nil
}

func (w *Writer) Write(src []byte) (int, error) {
	n, err := w.buffer.Write(src)
	if err != nil {
		return 0, err
	}

	for w.buffer.Len() >= blockSize {
		err := w.encryptBlock()
		if err != nil {
			return 0, err
		}
	}
	return n, nil
}

func (w *Writer) init(pass []byte) error {
	var err error
	// Creates the header
	if w.header, err = NewHeader(); err != nil {
		return err
	}

	// Generates the key
	w.key = generateKey(pass, w.header.Salt[:])

	// Creates the encryption stream
	if w.s, err = createStream(w.key, w.header.Iv[:]); err != nil {
		return err
	}

	// Writes the header into the writer
	if err = w.header.Write(w.w); err != nil {
		return err
	}
	return nil
}

func (w *Writer) encryptBlock() error {
	obuf := make([]byte, lenSize+blockSize+macSize)
	bbuf := obuf[:lenSize+blockSize]

	// Reads block from buffer
	n, berr := w.buffer.Read(bbuf[lenSize:])
	if n > 0 {
		binary.BigEndian.PutUint16(bbuf[:lenSize], uint16(n))

		// Calculates len and block MAC
		mac := hmac.New(sha256.New, w.key)
		if _, err := mac.Write(bbuf); err != nil {
			return err
		}
		copy(obuf[lenSize+blockSize:], mac.Sum(nil))

		// Encrypt len and block
		w.s.XORKeyStream(bbuf, bbuf)

		// Writes encrypted block to writer
		if _, err := w.w.Write(obuf); err != nil {
			return err
		}
	}
	return berr
}
