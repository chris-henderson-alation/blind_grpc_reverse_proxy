package common // import "github.com/Alation/alation_connector_manager/security/common"

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"time"

	"golang.org/x/crypto/nacl/box"

	"github.com/pkg/errors"
)

func Marshal(key interface{}) [32]byte {
	var b []byte
	switch k := key.(type) {
	case *ecdsa.PrivateKey:
		b = elliptic.MarshalCompressed(k.Curve, k.X, k.Y)
	case *ecdsa.PublicKey:
		b = elliptic.MarshalCompressed(k.Curve, k.X, k.Y)
	default:
		panic("unsupported key typ")
	}
	var m [32]byte
	copy(m[:], b)
	return m
}

func Encrypt(plaintext []byte, peerPublicKey, privateKey *[32]byte) ([]byte, []byte) {
	var nonce [24]byte
	rand.Read(nonce[:])
	return box.Seal([]byte{}, plaintext, &nonce, peerPublicKey, privateKey), nonce[:]
}

func Decrypt(ciphertext []byte, nonce []byte, peerPublicKey, privateKey *[32]byte) ([]byte, error) {
	var n [24]byte
	copy(n[:], nonce)
	plaintext, worked := box.Open([]byte{}, ciphertext, &n, peerPublicKey, privateKey)
	if !worked {
		return nil, errors.New("bad decryption")
	}
	return plaintext, nil
}

// GetPublicCertificate reads and deserializes to memory an x509 public certificate.
func GetPublicCert(certPath string) (*x509.Certificate, error) {
	c, err := ioutil.ReadFile(certPath)
	if err != nil {
		return nil, err
	}
	block, _ := pem.Decode(c)
	if block == nil {
		panic(fmt.Sprintf("nil private cert pem at %s", certPath))
	}
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, err
	}
	return cert, nil
}

// GetPrivateKey reads the provided file path and decrypts a PEM encoded PRIVATE KEY.
func GetPrivateKey(keyPath string, password []byte) (interface{}, error) {
	data, err := ioutil.ReadFile(keyPath)
	if err != nil {
		return nil, err
	}
	key, err := DecryptPrivateKey(data, password)
	if err != nil {
		return nil, err
	}
	return key, nil
}

// NewKeyPair returns a new ECDSA key configured for the NIST Curve P-256
// https://csrc.nist.gov/csrc/media/events/workshop-on-elliptic-curve-cryptography-standards/documents/papers/session6-adalier-mehmet.pdf
func NewKeyPair() (*ecdsa.PrivateKey, error) {
	return ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
}

// EncryptPrivateKey returns a PEM encoded document over the provided key and encrypted by "password".
func EncryptPrivateKey(key interface{}, password []byte) ([]byte, error) {
	keyBytes, err := x509.MarshalPKCS8PrivateKey(key)
	if err != nil {
		log.Printf("Unable to marshal private key: %v\n", err)
		return nil, err
	}
	return EncryptPEM(keyBytes, password, "PRIVATE KEY")
}

func DecryptPrivateKey(data, password []byte) (interface{}, error) {
	p, err := DecryptPEM(data, password)
	if err != nil {
		return nil, err
	}
	return x509.ParsePKCS8PrivateKey(p)
}

// EncryptPEM takes data, a password, and a PEM type and generates an encrypted, in-memory, PEM of
// that type over the provided data.
func EncryptPEM(data []byte, password []byte, blockType string) ([]byte, error) {
	block, err := x509.EncryptPEMBlock(rand.Reader, blockType, data, password, x509.PEMCipherAES256)
	if err != nil {
		return nil, err
	}
	return pem.EncodeToMemory(block), nil
}

func DecryptPEM(data, password []byte) ([]byte, error) {
	privPem, _ := pem.Decode(data)
	if privPem == nil {
		return nil, errors.New("The provided private key data does not appear to be a valid PEM format. See" +
			"https://en.wikipedia.org/wiki/Privacy-Enhanced_Mail for details on what a PEM file should look like.")
	}
	return x509.DecryptPEMBlock(privPem, password)
}

const END_CERTIFICATE = "-----END CERTIFICATE-----"

// PEMChainToSlice breaks up a flat certificate chain into a slice of individual PEM encoded certificates.
func PEMChainToSlice(chain []byte) [][]byte {
	count := bytes.Count(chain, []byte(END_CERTIFICATE))
	ret := make([][]byte, count)
	for i, val := range bytes.SplitAfterN(chain, []byte(END_CERTIFICATE), count) {
		ret[i] = append(bytes.TrimSpace(val), '\n')
	}
	return ret
}

// NewSelfSignedCertificate takes in a private key and generates a self signed certificate.
// This self signed certificate is intended for use in the certificate issuance protocol
// after the private key has been derived by a salted key derivation function.
//
// It can also be used to generate the certificate authority's trust anchor.
func NewSelfSignedCertificate(priv *ecdsa.PrivateKey) (*x509.Certificate, error) {
	// @TODO make this default much shorter and make it configurable.
	notBefore, notAfter := NewValidityRange(time.Hour * 24 * 365)
	template := x509.Certificate{
		SerialNumber:          NewSerial(),
		SignatureAlgorithm:    x509.ECDSAWithSHA256,
		NotBefore:             notBefore,
		NotAfter:              notAfter,
		KeyUsage:              x509.KeyUsageCertSign,
		BasicConstraintsValid: true,
		IsCA:                  true,
	}
	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, priv.Public(), priv)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to create a self singed certificate.")
	}
	return x509.ParseCertificate(derBytes)
}

// NewValidityRange returns a NotBefore and NotAfter as described by
// RC5280 Section 4.1.2.5 (https://tools.ietf.org/html/rfc5280#section-4.1.2.5)
func NewValidityRange(duration time.Duration) (time.Time, time.Time) {
	notBefore := time.Now()
	notAfter := notBefore.Add(duration)
	return notBefore, notAfter
}

// NewSerial generates a 128, cryptographically generated, UUID.
// It is intended for use to assign serial numbers during
// the certificate signing ceremony.
func NewSerial() *big.Int {
	serial := make([]byte, 16)
	_, err := rand.Read(serial)
	if err != nil {
		log.Fatalf("Failed to generate a serial number using crypto/rand: %s", err)
	}
	return big.NewInt(0).SetBytes(serial)
}
