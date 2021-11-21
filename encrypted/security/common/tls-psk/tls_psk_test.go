package tls_psk

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"errors"
	"fmt"
	"net"
	"reflect"
	"sync"
	"testing"
	"time"

	"github.com/Alation/alation_connector_manager/security/common"
)

func TestIntegration(t *testing.T) {
	server, client := net.Pipe()
	secret := []byte("super secret squirrel")
	message := []byte("Hello world!")
	wg := sync.WaitGroup{}
	wg.Add(2)
	errs := make(chan error, 2)
	go func() {
		defer wg.Done()
		conn, err := UpgradeServerConnection(server, secret)
		if err != nil {
			errs <- err
			return
		}
		got := make([]byte, len(message))
		_, err = conn.Read(got)
		if err != nil {
			errs <- err
			return
		}
		if !bytes.Equal(got, message) {
			// It honestly more useful to leave the represented as a byte slice, as debugging this with
			// string encodings can hide things from you.
			errs <- errors.New(fmt.Sprintf(`Unexpected result. Got "%v", wanted "%v"`, got, message))
		}
	}()
	go func() {
		defer wg.Done()
		conn, err := UpgradeClientConnection(client, secret)
		if err != nil {
			errs <- err
			return
		}
		_, err = conn.Write(message)
		if err != nil {
			errs <- err
			return
		}
	}()
	wg.Wait()
	close(errs)
	for err := range errs {
		t.Error(err)
	}
}

func TestIntegrationAuthFail(t *testing.T) {
	server, client := net.Pipe()
	// These in-memory pipes will hang indefinitely on a
	// bad handshake, so we're going to cause a timeout instead.
	server.SetDeadline(time.Now().Add(1))
	client.SetDeadline(time.Now().Add(1))
	serverSecret := []byte("super secret squirrel")
	clientSecret := []byte("morocco mole")
	wg := sync.WaitGroup{}
	wg.Add(2)
	errs := make(chan error, 2)
	go func() {
		defer wg.Done()
		_, err := UpgradeServerConnection(server, serverSecret)
		if err == nil {
			errs <- errors.New("server did not lock out an untrusted client")
		}
	}()
	go func() {
		defer wg.Done()
		_, err := UpgradeClientConnection(client, clientSecret)
		if err == nil {
			errs <- errors.New("client did not lock out an untrusted server")
		}
	}()
	wg.Wait()
	close(errs)
	for err := range errs {
		t.Error(err)
	}
}

func TestIntegrationClosedClient(t *testing.T) {
	server, client := net.Pipe()
	client.Close()
	secret := []byte("super secret squirrel")
	_, err := UpgradeServerConnection(server, secret)
	if err == nil {
		t.Fatal("server did not lock out an untrusted client")
	}
}

func TestIntegrationClosedServer(t *testing.T) {
	server, client := net.Pipe()
	server.Close()
	secret := []byte("super secret squirrel")
	_, err := UpgradeClientConnection(client, secret)
	if err == nil {
		t.Fatal("client did not lock out an untrusted server")
	}
}

func TestKeyGeneration(t *testing.T) {
	server, client := net.Pipe()
	secret := []byte("super secret squirrel")
	serverKey := make([]byte, elliptic.P256().Params().BitSize/8+8)
	clientKey := make([]byte, elliptic.P256().Params().BitSize/8+8)
	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		defer wg.Done()
		if _, err := NewServerKeyReader(secret, server).Read(serverKey); err != nil {
			t.Error(err)
		}
	}()
	go func() {
		defer wg.Done()
		if _, err := NewClientKeyReader(secret, client).Read(clientKey); err != nil {
			t.Error(err)
		}
	}()
	wg.Wait()
	if !reflect.DeepEqual(serverKey, clientKey) {
		t.Fatalf("serverKey and clientKey differed\nServer: %X\nClient: %X", serverKey, clientKey)
	}
}

func TestKeyECCGeneration(t *testing.T) {
	server, client := net.Pipe()
	secret := []byte("super secret squirrel")
	var serverKey *ecdsa.PrivateKey
	var clientKey *ecdsa.PrivateKey
	wg := sync.WaitGroup{}
	wg.Add(2)
	fail := make(chan error, 2)
	go func() {
		defer wg.Done()
		var err error
		serverKey, err = ecdsa.GenerateKey(elliptic.P256(), NewServerKeyReader(secret, server))
		fail <- err
	}()
	go func() {
		defer wg.Done()
		var err error
		clientKey, err = ecdsa.GenerateKey(elliptic.P256(), NewClientKeyReader(secret, client))
		fail <- err
	}()
	wg.Wait()
	close(fail)
	var failed bool
	for err := range fail {
		if err != nil {
			failed = true
			t.Error(err)
		}
	}
	if failed {
		return
	}
	if !reflect.DeepEqual(serverKey, clientKey) {
		t.Fatalf("serverKey and clientKey differed\nServer: %X\nClient: %X", serverKey, clientKey)
	}
}
func TestKeyECCMismatch(t *testing.T) {
	server, client := net.Pipe()
	serverSecret := []byte("super secret squirrel")
	clientSecret := []byte("morocco mole")
	var serverKey *ecdsa.PrivateKey
	var clientKey *ecdsa.PrivateKey
	wg := sync.WaitGroup{}
	wg.Add(2)
	fail := make(chan error, 2)
	go func() {
		defer wg.Done()
		var err error
		serverKey, err = ecdsa.GenerateKey(elliptic.P256(), NewServerKeyReader(serverSecret, server))
		fail <- err
	}()
	go func() {
		defer wg.Done()
		var err error
		clientKey, err = ecdsa.GenerateKey(elliptic.P256(), NewClientKeyReader(clientSecret, client))
		fail <- err
	}()
	wg.Wait()
	close(fail)
	var failed bool
	for err := range fail {
		if err != nil {
			failed = true
			t.Error(err)
		}
	}
	if failed {
		return
	}
	if reflect.DeepEqual(serverKey, clientKey) {
		t.Fatal("serverKey and clientKey were the same")
	}
}

func TestSelfSigned(t *testing.T) {
	server, client := net.Pipe()
	secret := []byte("super secret squirrel")
	var serverKey *ecdsa.PrivateKey
	var clientKey *ecdsa.PrivateKey
	wg := sync.WaitGroup{}
	wg.Add(2)
	fail := make(chan error, 2)
	go func() {
		defer wg.Done()
		var err error
		serverKey, err = ecdsa.GenerateKey(elliptic.P256(), NewServerKeyReader(secret, server))
		fail <- err
	}()
	go func() {
		defer wg.Done()
		var err error
		clientKey, err = ecdsa.GenerateKey(elliptic.P256(), NewClientKeyReader(secret, client))
		fail <- err
	}()
	wg.Wait()
	close(fail)
	var failed bool
	for err := range fail {
		if err != nil {
			failed = true
			t.Error(err)
		}
	}
	if failed {
		return
	}
	if !reflect.DeepEqual(serverKey, clientKey) {
		t.Fatalf("serverKey and clientKey differed\nServer: %X\nClient: %X", serverKey, clientKey)
	}
	serverCert, err := common.NewSelfSignedCertificate(serverKey)
	if err != nil {
		t.Fatal(err)
	}
	clientCert, err := common.NewSelfSignedCertificate(clientKey)
	if err != nil {
		t.Fatal(err)
	}
	serverRootPool := x509.NewCertPool()
	serverRootPool.AddCert(serverCert)
	clientRootPool := x509.NewCertPool()
	clientRootPool.AddCert(clientCert)
	// Make sure that the verifying a client cert works from the server perspective
	_, err = clientCert.Verify(x509.VerifyOptions{
		Roots: serverRootPool,
	})
	if err != nil {
		t.Fatal(err)
	}
	// Make sure that the verifying a server cert works from the client perspective
	_, err = serverCert.Verify(x509.VerifyOptions{
		Roots: clientRootPool,
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestSelfSignedMismatch(t *testing.T) {
	// Ensure that when the client and the server disagree on the ECC key generated that everything fails.
	serverKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	clientKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	serverCert, err := common.NewSelfSignedCertificate(serverKey)
	if err != nil {
		t.Fatal(err)
	}
	clientCert, err := common.NewSelfSignedCertificate(clientKey)
	if err != nil {
		t.Fatal(err)
	}
	serverRootPool := x509.NewCertPool()
	serverRootPool.AddCert(serverCert)
	clientRootPool := x509.NewCertPool()
	clientRootPool.AddCert(clientCert)
	// Make sure that the verifying a client cert fails from the server perspective
	_, err = clientCert.Verify(x509.VerifyOptions{
		Roots: serverRootPool,
	})
	if err == nil {
		t.Fatal("expected a verification failure")
	}
	// Make sure that the verifying a server cert fails from the client perspective
	_, err = serverCert.Verify(x509.VerifyOptions{
		Roots: clientRootPool,
	})
	if err == nil {
		t.Fatal("expected a verification failure")
	}
}
