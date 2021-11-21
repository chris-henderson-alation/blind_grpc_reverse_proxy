package kdf

import (
	"crypto/rand"
	"io"

	"golang.org/x/crypto/argon2"
)

const ARGON2_SALT_SIZE = 16
const ARGON_TIME = 3
const ARGON_64_MB = 64 * 1024

type KeyReader struct {
	saltSrc io.Reader
	secret  []byte
}

func NewKeyReader(secret []byte, saltSrc io.Reader) *KeyReader {
	if saltSrc == nil {
		saltSrc = rand.Reader
	}
	return &KeyReader{saltSrc, secret}
}

func (k *KeyReader) Read(p []byte) (n int, err error) {
	salt := make([]byte, ARGON2_SALT_SIZE)
	_, err = io.ReadFull(k.saltSrc, salt)
	if err != nil {
		return
	}
	n = copy(p, argon2.IDKey(k.secret, salt, ARGON_TIME, ARGON_64_MB, 1, uint32(len(p))))
	return
}
