package main

import (
	"bytes"
	"crypto/ecdsa"
	"testing"

	"github.com/Alation/alation_connector_manager/security/common"
)

func TestAasdas(t *testing.T) {
	key, err := common.GetPrivateKey("/home/chris/alation/blind_grpc_reverse_proxy/encrypted/agent/private_key.pem", []byte("password@1234"))
	if err != nil {
		panic(err)
	}
	private := key.(*ecdsa.PrivateKey)
	pubkey, err := common.GetPublicCert("/home/chris/alation/blind_grpc_reverse_proxy/encrypted/servicer/public_certificate.pem")
	if err != nil {
		panic(err)
	}
	public := pubkey.PublicKey.(*ecdsa.PublicKey)
	myPrivateKey = common.Marshal(private)
	theirPublicKey = common.Marshal(public)

	plaintext := []byte("naughty secrets")
	ciphertext, nonce := common.Encrypt(plaintext, &theirPublicKey, &myPrivateKey)
	t.Log(ciphertext)
	ptext, err := common.Decrypt(ciphertext, nonce, &theirPublicKey, &myPrivateKey)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(ptext))
	t.Log(bytes.Equal(ptext, plaintext))
	b, err := marshal(makeBall())
	if err != nil {
		t.Fatal(err)
	}
	c, err := common.Decrypt(b.Ciphertext, b.Nonce, &theirPublicKey, &myPrivateKey)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(c)
}
