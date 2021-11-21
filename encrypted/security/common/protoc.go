package common

import (
	"net"
	"time"

	"github.com/pkg/errors"
)

type Request byte

const (
	GetPublicCertificate Request = 1 << iota
	GetIssuance
)

func (r Request) SendInitiate(conn net.Conn) error {
	conn.SetDeadline(time.Now().Add(time.Second))
	defer conn.SetDeadline(time.Time{})
	_, err := conn.Write([]byte{byte(r)})
	return err
}

func RecvInitiate(conn net.Conn) (Request, error) {
	conn.SetDeadline(time.Now().Add(time.Second))
	defer conn.SetDeadline(time.Time{})
	r := make([]byte, 1)
	n, err := conn.Read(r)
	if err != nil {
		return 0, err
	}
	if n != 1 {
		return 0, errors.Errorf("failed to read request type byte from %v", conn.RemoteAddr())
	}
	switch Request(r[0]) {
	case GetPublicCertificate:
		return GetPublicCertificate, nil
	case GetIssuance:
		return GetIssuance, nil
	default:
		return 0, errors.Errorf("unknown request type %v from %v", r, conn.RemoteAddr())
	}
}
