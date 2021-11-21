package common

import (
	"crypto/x509"
	"reflect"
	"testing"
)

const chainLeaf = `-----BEGIN CERTIFICATE-----
MIICxDCCAmmgAwIBAgIRAJEDWJB4gRKIHdaDqDYr3hQwCgYIKoZIzj0EAwIwgaEx
CzAJBgNVBAYTAlVTMQswCQYDVQQIEwJDQTEVMBMGA1UEBxMMUmVkd29vZCBDaXR5
MTAwLgYDVQQJEyczIExhZ29vbiBEcml2ZSwgU3VpdGUgMzAwLCBSZWR3b29kIENp
dHkxDjAMBgNVBBETBTk0MDY1MRYwFAYDVQQKEw1BbGF0aW9uLCBJbmMuMRQwEgYD
VQQLEwtFbmdpbmVlcmluZzAeFw0xOTEwMjgyMzM2NDdaFw0yMDEwMjcyMzM2NDda
MIHIMQswCQYDVQQGEwJVUzELMAkGA1UECBMCQ0ExFTATBgNVBAcTDFJlZHdvb2Qg
Q2l0eTEwMC4GA1UECRMnMyBMYWdvb24gRHJpdmUsIFN1aXRlIDMwMCwgUmVkd29v
ZCBDaXR5MQ4wDAYDVQQREwU5NDA2NTEWMBQGA1UEChMNQWxhdGlvbiwgSW5jLjEU
MBIGA1UECxMLRW5naW5lZXJpbmcxJTAjBgNVBAMTHGNocmlzLWhlbmRlcnNvbi1D
MDJXUTE0Rkc4V0wwWTATBgcqhkjOPQIBBggqhkjOPQMBBwNCAATfX6DcrK6fQHNn
1xoa5sSC2bdGnmuXlLwCW8lyCI/YU4yKqZiCFhDxaNwJayvAeUH6+FOlKDxvGHc2
b/EOI/g+o1kwVzAOBgNVHQ8BAf8EBAMCB4AwHQYDVR0lBBYwFAYIKwYBBQUHAwEG
CCsGAQUFBwMCMCYGA1UdEQQfMB2BG2NocmlzLmhlbmRlcnNvbkBhbGF0aW9uLmNv
bTAKBggqhkjOPQQDAgNJADBGAiEAh/oYpFCDrZC2CNJQUehTnRwyH6gBeaGfsSdk
AzEgp88CIQCNVrjXOEzuiyB1qunYXnCJJMtSaWDpCzLwOWIIR1Czhg==
-----END CERTIFICATE-----
`
const chainRoot = `-----BEGIN CERTIFICATE-----
MIICZDCCAgugAwIBAgIQeJOePjnDxmTSRAmB8efSWTAKBggqhkjOPQQDAjCBoTEL
MAkGA1UEBhMCVVMxCzAJBgNVBAgTAkNBMRUwEwYDVQQHEwxSZWR3b29kIENpdHkx
MDAuBgNVBAkTJzMgTGFnb29uIERyaXZlLCBTdWl0ZSAzMDAsIFJlZHdvb2QgQ2l0
eTEOMAwGA1UEERMFOTQwNjUxFjAUBgNVBAoTDUFsYXRpb24sIEluYy4xFDASBgNV
BAsTC0VuZ2luZWVyaW5nMB4XDTE5MTAyODIzMjk0NFoXDTIwMTAyNzIzMjk0NFow
gaExCzAJBgNVBAYTAlVTMQswCQYDVQQIEwJDQTEVMBMGA1UEBxMMUmVkd29vZCBD
aXR5MTAwLgYDVQQJEyczIExhZ29vbiBEcml2ZSwgU3VpdGUgMzAwLCBSZWR3b29k
IENpdHkxDjAMBgNVBBETBTk0MDY1MRYwFAYDVQQKEw1BbGF0aW9uLCBJbmMuMRQw
EgYDVQQLEwtFbmdpbmVlcmluZzBZMBMGByqGSM49AgEGCCqGSM49AwEHA0IABNbL
HNd5oyXq4c4VKf5px7YaU8ZxKgSwkSMB1mwcLrM2BcpnPxZEL/lE6iOvbbXz22C+
Fg/YOxCEbLNy2YsBFX2jIzAhMA4GA1UdDwEB/wQEAwICBDAPBgNVHRMBAf8EBTAD
AQH/MAoGCCqGSM49BAMCA0cAMEQCIFszCyAMCPgEZTTCIOIlQVRtOVt4u7xdaM4N
2eLyvw4XAiAaZIWSIHbuV/sYjTI2nRH9sXcGFuen4UPN9c9y7mBdig==
-----END CERTIFICATE-----
`

func TestChainToSlice(t *testing.T) {
	c := PEMChainToSlice(append([]byte(chainLeaf), []byte(chainRoot)...))
	if len(c) != 2 {
		t.Fatalf("wanted 2 certs in the chain, got %d", len(c))
	}
	if !reflect.DeepEqual(c[0], []byte(chainLeaf)) {
		t.Errorf("expected the chainLeaf, got '%v'", string(c[0]))
	}
	if !reflect.DeepEqual(c[1], []byte(chainRoot)) {
		t.Errorf("expected the chainRoot, got %v", string(c[1]))
	}
}

func TestSingeCertInChain(t *testing.T) {
	c := PEMChainToSlice([]byte(chainLeaf))
	if len(c) != 1 {
		t.Fatalf("wanted 1 cert in the chain, got %d", len(c))
	}
}

func TestEmptyChain(t *testing.T) {
	c := PEMChainToSlice([]byte{})
	if len(c) != 0 {
		t.Fatalf("expected an empty slice, got %v", c)
	}
}

func TestNilChain(t *testing.T) {
	c := PEMChainToSlice(nil)
	if len(c) != 0 {
		t.Fatalf("expected an empty slice, got %v", c)
	}
}

func TestPrivKeyPem(t *testing.T) {
	priv, err := NewKeyPair()
	if err != nil {
		t.Fatal(err)
	}
	_, err = x509.MarshalPKCS8PrivateKey(priv)
	if err != nil {
		t.Fatal(err)
	}
}

func TestEncDecPEM(t *testing.T) {
	private, err := NewKeyPair()
	if err != nil {
		t.Fatal(err)
	}
	password := []byte("secret")
	p, err := EncryptPrivateKey(private, password)
	if err != nil {
		t.Fatal(err)
	}
	got, err := DecryptPrivateKey(p, []byte("secret"))
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(private, got) {
		t.Fatalf("the decrypted result of the private did not match the original. Got %X, want %X", got, private)
	}
}
