package tls_psk

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"io"
	"net"
	"time"

	"github.com/pkg/errors"

	"github.com/Alation/alation_connector_manager/security/common"
	"github.com/Alation/alation_connector_manager/security/common/kdf"
)

// The existence of this package is an unfortunate side effect of https://github.com/golang/go/issues/6379
// remaining open and essentially in a Will Not Do state. What we would really like is a standard library
// implementation of a TLS-PSK (TLS Pre-Shared Secret) (https://tools.ietf.org/html/rfc4279), which can
// bootstrap a TLS connection from a pre-shared secret rather than x509 certificates.
//
// In lieu of such convenience, the following procedures use the Argon2 key derivation function
// (https://password-hashing.net/argon2-specs.pdf) (https://en.wikipedia.org/wiki/Argon2) to derive a
// NIST P-256 curve that can then be used to upgrade the raw TCP connection to a standard TLS connection.
//
// The algorithm is as follows:
//	1. The server and the client independently generate, and exchange, cryptographically secure nonces.
// 		1.a. Demanding that both entities generate nonces protects against poor entropy from either of the actors.
// 	2. Each actor creates a salt that is the concatenation of the serverNonce || clientNonce.
//	3. The resulting salt is fed into the Argon2ID function along with the issuance secret
//	4. A new key, that is the length of the a NIST P-256 curve, is derived.
//	5. This new key is used to generate the same public/private key pair on the both the server and the client.
//	6. The server and the client engage a mutual TLS handshake, wherein both demand that the other's public certificate
//			was signed by their own private key.
//	7. The connection is now upgraded to a normal TLS connection and secure communications can confidently proceed.

func UpgradeServerConnection(conn net.Conn, secret []byte) (*tls.Conn, error) {
	priv, err := ecdsa.GenerateKey(elliptic.P256(), NewServerKeyReader(secret, conn))
	if err != nil {
		return nil, err
	}
	cert, err := common.NewSelfSignedCertificate(priv)
	if err != nil {
		return nil, err
	}
	pool := x509.NewCertPool()
	pool.AddCert(cert)
	config := &tls.Config{
		Certificates: []tls.Certificate{{
			Certificate: [][]byte{cert.Raw},
			PrivateKey:  priv,
			Leaf:        cert,
		}},
		// Lower versions of TLS are insecure, however
		// TLS1.3 is not supported by all networking stacks.
		// An incompatibility here will manifest as a
		// "Error: tls: oversized record received with length 29045"
		// or something similar. This is not something that WE can
		// deal with as every layer (OS, ISP, switches) needs
		// to be compatible. So set both the minimum and the maximum
		// to TLS1.2 for a happy medium of security and compatiblity.
		MinVersion: tls.VersionTLS12,
		MaxVersion: tls.VersionTLS12,
		// ClientAuth enables checking of a client provided
		// certificate in order to authenticate them. This is
		// because normally a client is an entirely anonymous
		// entity to a given server (from a TLS perspective).
		ClientAuth:               tls.RequireAnyClientCert,
		ClientCAs:                pool,
		PreferServerCipherSuites: true,
		VerifyPeerCertificate:    verifyFunc(pool),
	}
	server := tls.Server(conn, config)
	return server, server.Handshake()
}

func UpgradeClientConnection(conn net.Conn, secret []byte) (*tls.Conn, error) {
	priv, err := ecdsa.GenerateKey(elliptic.P256(), NewClientKeyReader(secret, conn))
	if err != nil {
		return nil, err
	}
	cert, err := common.NewSelfSignedCertificate(priv)
	if err != nil {
		return nil, err
	}
	roots := x509.NewCertPool()
	roots.AddCert(cert)
	config := &tls.Config{
		Certificates: []tls.Certificate{{
			Certificate: [][]byte{cert.Raw},
			PrivateKey:  priv,
			Leaf:        cert,
		}},
		// Lower versions of TLS are insecure, however
		// TLS1.3 is not supported by all networking stacks.
		// An incompatibility here will manifest as a
		// "Error: tls: oversized record received with length 29045"
		// or something similar. This is not something that WE can
		// deal with as every layer (OS, ISP, switches) needs
		// to be compatible. So set both the minimum and the maximum
		// to TLS1.2 for a happy medium of security and compatiblity.
		MinVersion: tls.VersionTLS12,
		MaxVersion: tls.VersionTLS12,
		// If you leave InsecureSkipVerify to true then it will not only
		// check the signature of the provided certificate, but it will
		// also demand that the DNS entry lines up. These are just
		// ephemeral, session level, certificates and we cannot reasonably
		// attach any DNS semantics to them.
		//
		// As such, we disable that procedure and instead provide our own
		// verification function that WILL do signature verification using
		// the provided root certificates, however it will skip DNS
		// verification.
		InsecureSkipVerify:    true,
		VerifyPeerCertificate: verifyFunc(roots),
		RootCAs:               roots,
	}
	client := tls.Client(conn, config)
	return client, client.Handshake()
}

// The differences between a ServerKeyReader and a ClientKeyReader are:
//
//	1. A server writes and then reads. A client reads and then writes.
//	2. A sever generates ourNonce || theirNonce. A client generates theirNone || ourNonce
//
// #1 is to ensure a deadlock does not occur on IO. #2 is to ensure that both the client
// and the server generate the same salt for the argon2id function.
type ServerKeyReader struct {
	conn io.ReadWriter
}

type ClientKeyReader struct {
	conn io.ReadWriter
}

func NewServerKeyReader(secret []byte, conn io.ReadWriter) io.Reader {
	return kdf.NewKeyReader(secret, &ServerKeyReader{conn: conn})
}

func NewClientKeyReader(secret []byte, conn io.ReadWriter) io.Reader {
	return kdf.NewKeyReader(secret, &ClientKeyReader{conn: conn})
}

// 1. Generate
// 2. Write
// 3. Read
// 4. Return ours || theirs
func (s *ServerKeyReader) Read(p []byte) (n int, err error) {
	ourSalt := make([]byte, len(p))
	_, err = rand.Read(ourSalt)
	if err != nil {
		return 0, err
	}
	if _, err := s.conn.Write(ourSalt); err != nil {
		return 0, err
	}
	theirSalt := make([]byte, len(p))
	if _, err := io.ReadFull(s.conn, theirSalt); err != nil {
		return 0, err
	}
	n = copy(p, append(ourSalt, theirSalt...))
	return
}

// 1. Generate
// 2. Read
// 3. Write
// 4. Return theirs || ours
func (c *ClientKeyReader) Read(p []byte) (n int, err error) {
	ourSalt := make([]byte, len(p))
	_, err = rand.Read(ourSalt)
	if err != nil {
		return 0, err
	}
	theirSalt := make([]byte, len(p))
	if _, err := io.ReadFull(c.conn, theirSalt); err != nil {
		return 0, err
	}
	if _, err := c.conn.Write(ourSalt); err != nil {
		return 0, err
	}
	n = copy(p, append(theirSalt, ourSalt...))
	return
}

// verifyFunc returns a function that satisfies the tls.Config.VerifyPeerCertificate.
// The returned function verifies that any given certificate was signed by at least one
// root certificate from within `roots`. Other verification options, such as DNS verification,
// are skipped.
func verifyFunc(roots *x509.CertPool) func(certs [][]byte, _ [][]*x509.Certificate) error {
	return func(certs [][]byte, _ [][]*x509.Certificate) error {
		cert, err := x509.ParseCertificate(certs[0])
		if err != nil {
			return err
		}
		_, err = cert.Verify(x509.VerifyOptions{
			Roots:       roots,
			CurrentTime: time.Now(),
		})
		// Unfortunately, the x509 library is not typing their errors, so if we
		// REALLY wanted to differentiate all of the different things that can go
		// wrong in certificate verification then we would have to start doing string
		// inspection, which is just a terribly flaky thing to do as those strings
		// can be changed at any time.
		//
		// Almost always an error at this point is a bad password:
		//
		// Example: 'x509: certificate signed by unknown authority (possibly because of "x509: ECDSA verification failure" while trying to verify candidate authority certificate "serial:54928905515025366959284855478945917488")'
		//
		// Anything else is likely caused by a bad software deployment (say, perhaps, if we change how this works
		// and the certificate authority is upgraded but the client is left on an old version).
		return errors.Wrap(err, "Incorrect TLS shared secret. For debugging purposes only, the following is the raw error")
	}
}
