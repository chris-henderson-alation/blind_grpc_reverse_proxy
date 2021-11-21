# Certificate Authority

## Mutual TLS
The architecture working group for security agreed that a scheme of [mutual TLS](https://en.wikipedia.org/wiki/Mutual_authentication) was an attractive and common practice among microservice implementations. The details of TLS, mutual TLS, and public key infrastructures could easily overwhelm this document, however I would like to offer a brief summary of the technique (for a delightfully in-depth review of TLS 1.3, please see [The New Illustrated TLS Connection](https://tls13.ulfheim.net/)).

### How does TLS work from day-to-day? (Unidirectional Authentication)
Common consumer TLS features uni-directional authentication. That is, the TLS certificate of amazon.com serves as proof that a client is in fact speaking with Amazon Inc. A client will attempt to ascertain this by (at minimum):

1. Performing chain-of-trust resolution in order to verify that the certificate presented to the client was signed by a certificate authority that the client trusts apriori.
2. Presenting the remote server with a cryptographic challenge using the public key encoded within the provided certificate in order to prove that the server also holds the complementing private key.

Clients may also wish to verify the certificate’s [Extended Key Usage](https://tools.ietf.org/html/rfc5280#section-4.2.1.12) for authorization to respond as a TLS server (`id-kp-serverAuth`), demand that the certificate’s [Subject Alternative Name](https://tools.ietf.org/html/rfc5280#section-4.2.1.6) is equivalent to the DNS name used to resolve the host, or request a response from a [Certificate Revocation List](https://tools.ietf.org/html/rfc5280#section-5) (CRL) responder or [Online Certificate Status Protocol](https://tools.ietf.org/html/rfc6960) (OCSP) responder.

After these checks are completed, the connection has been successfully established and the website has proven that they are, in fact, Amazon.

### Bidirectional Authentication
However, Amazon has no proof that the client is anyone but an anonymous client on the internet. In order to conduct a counter-authentication, the client must now go through Amazon’s user/password backed authentication mechanisms. This works well for consumer-to-business interactions. However, it is unwieldy in machine-to-machine microservice architectures.

In mutual TLS, the above authentication flow (in this particular case, amazon.com) is conducted - not only for the server, but for the client as well. That is, the client must now offer up a TLS certificate and answer a challenge that the server finds satisfactory. In this way, you can build small, autonomous, [public key infrastructures](https://en.wikipedia.org/wiki/Public_key_infrastructure).

## Default Bundled CA

The following section details an optional, ready-out-of-the-box, certificate authority and cluster issuance protocol that comes with each installation of ACM.

As of this moment, enabling TLS enables the certificate authority. A flag remains to be added that can turn off the CA while leaving TLS active. That is, an operator can turn off the CA and install their own appropriate certificates. There is no real difficulty in doing so, it is merely a piece of backlog work that hasn’t been bumped up yet.

### Certificate Signing over TLS-PSK
The signing of a public certificate broadly consists of two steps - a challenge and a certificate signing ceremony.

#### Background
In common day-to-day issuances (where an individual or business wishes to receive a certificate) the challenge often consists of human-to-human interactions wherein the identity and ownership of the requested domain are verified (documentation, email and phone correspondence, etc).

However, we are disinterested in manual steps and would prefer automated protocols. We can take, for example, the Automatic Certificate Management Environment (ACME, [RFC 8555](https://tools.ietf.org/html/rfc8555)) protocol which powers the ever popular [Let’s Encrypt](https://letsencrypt.org/) certificate authority. We must take a variant on it, however, as ACME does not support the issuance of a certificate without [domain validation](https://en.wikipedia.org/wiki/Domain-validated_certificate). Not only are our components mostly unexposed to the internet, they may often have no DNS entry associated with them at all - instead being only tied to an IP.

As such, we branch ACME slightly with a challenge that is based solely on proof that a client holds a secret that has been shared out-of-band (you will often see this in distributed systems, for example RabbitMQ’s [Erlang cookie](https://www.rabbitmq.com/clustering.html#erlang-cookie)). If the client holds the same secret as the certificate authority, then they may engage mutual TLS with a pre-shared key ([TLS-PSK](https://en.wikipedia.org/wiki/TLS-PSK)). That is, a fully functional, bi-directionally authenticated, TLS connection may be established if both the client and the certificate authority hold the same secret. Once the mutually authenticated, and encrypted, channel has been established the certificate signing ceremony may take place as normal.

#### Certificate Signing Algorithm
The concrete algorithm is as follows:

```
1. The server and the client independently generate, and exchange, cryptographically secure 128 bit nonces.
    a. Demanding that both entities generate nonces protects against poor entropy from either of the actors 
        (for example, if the server is the only generator of the nonce, and they are malicious, then they 
        may produce predictable nonces (say, all zeroes) in order to trick clients into gradually leaking secrets).
2. Each actor creates a salt that is the concatenation of the serverNonce || clientNonce.
    1. For effective security, we only need a 128bit nonce. However, constructing a 256bit nonce protects
        us from one of the actors attempting to introduce predictability.
3. The resulting salt is fed into the Argon2ID key derivation function along with the issuance secret.
    a. Cost parameters for Argon2ID are 3 rounds, a maximum memory limit of 64MB, and a single CPU core.
4. A new, ephemeral, key that is the length of the [NIST P-256 ECC](https://csrc.nist.gov/csrc/media/events/workshop-on-elliptic-curve-cryptography-standards/documents/papers/session6-adalier-mehmet.pdf) curve is derived.
5. This new key is used to generate the same public/private key pair on both the server and the client.
6. The server and the client then engage in a mutual TLS handshake, wherein both demand that the other's 
    public certificate was signed by their own private key (this can be done as each entity will generate 
    the same key via step 4).
    a. If certificate verification fails, then the connection is immediately terminated and an “incorrect secret” 
        posture is assumed.
7. The TCP connection is now upgraded to a normal, bi-directionally authenticated, TLS connection and secure 
    communications can confidently proceed.
```

- For more information regarding the chosen NIST curve, please see [NIST P-256 ECC](https://csrc.nist.gov/csrc/media/events/workshop-on-elliptic-curve-cryptography-standards/documents/papers/session6-adalier-mehmet.pdf).

We do not use this TLS-PSK technique for day-to-day operations as we are using Nginx for our reverse proxy and it is not aware of using TLS-PSK in the way (it is most typically used for session resumption on low power devices that cannot effectively execute asymmetric algorithms).

#### Certificate Template
All generated certificates are of the following template.

#### Common Name/Subject Alternative Name
The common name (as well as the subject alternative name) for each generated certificate is derived from the hostname in the [agent] address setting within the agent configuration. If no such field is configured, then the result of the hostname system call is used.

#### Extended Key Usage
Each certificate is authorized for the following usages:
- TLS Web Server Authentication ([OID 1.3.6.1.5.5.7.3.1](http://www.oid-info.com/cgi-bin/display?oid=1.3.6.1.5.5.7.3.1&action=display))
- TLS Web Client Authentication ([OID 1.3.6.1.5.5.7.3.2](http://www.oid-info.com/cgi-bin/display?oid=1.3.6.1.5.5.7.3.2&action=display))

#### Algorithms
- Key Algorithm: Elliptic Curve Public Key Cryptography ([OID 1.2.840.10045.2.1](http://www.oid-info.com/cgi-bin/display?oid=1.2.840.10045.2.1&action=display))
- Key Curve: NIST P-256 ([OID 1.2.840.10045.3.1.7](http://www.oid-info.com/cgi-bin/display?oid=1.2.840.10045.3.1.7&action=display))
- Signature: ECDSA-SHA256 ([OID 1.2.840.10045.4.3.2](http://www.oid-info.com/cgi-bin/display?oid=1.2.840.10045.4.3.2&action=display))

#### Serial Numbers
Serial numbers are not monotonically increasing. Rather, they are randomly generated 128bit UUIDs generated via `/dev/urandom`.

The following is an example certificate and it’s raw [ASN.1 decoding](https://lapo.it/asn1js/#MIICBTCCAaugAwIBAgIQX232Sg8KC_hD8QTJ2-65mjAKBggqhkjOPQQDAjAAMB4XDTIwMDUyMTAwMDA0MVoXDTIxMDUyMTAwMDA0MVowgboxCzAJBgNVBAYTAlVTMQswCQYDVQQIEwJDQTEVMBMGA1UEBxMMUmVkd29vZCBDaXR5MTAwLgYDVQQJEyczIExhZ29vbiBEcml2ZSwgU3VpdGUgMzAwLCBSZWR3b29kIENpdHkxDjAMBgNVBBETBTk0MDY1MRYwFAYDVQQKEw1BbGF0aW9uLCBJbmMuMRQwEgYDVQQLEwtFbmdpbmVlcmluZzEXMBUGA1UEAxMOZGFyay1jaGFyaXphcmQwWTATBgcqhkjOPQIBBggqhkjOPQMBBwNCAAQVWOyffe_4vGtnHoSo6m0juLUxlpslfUx4986ID0KFD1X3b_ttqGPV27Jhq8osyXhtNyBbKpCXmKb_o22jbcFlo0wwSjAOBgNVHQ8BAf8EBAMCB4AwHQYDVR0lBBYwFAYIKwYBBQUHAwEGCCsGAQUFBwMCMBkGA1UdEQQSMBCCDmRhcmstY2hhcml6YXJkMAoGCCqGSM49BAMCA0gAMEUCIQCxh7dUGmlAjMQI09Ohjznay8HtYfkzmAh0yblQmapdJAIgVsueyY6eWdFGboRG1Yk3Hu3wtnezI2urjw_iEVAE39M).

```
-----BEGIN CERTIFICATE-----
MIICBTCCAaugAwIBAgIQX232Sg8KC/hD8QTJ2+65mjAKBggqhkjOPQQDAjAAMB4X
DTIwMDUyMTAwMDA0MVoXDTIxMDUyMTAwMDA0MVowgboxCzAJBgNVBAYTAlVTMQsw
CQYDVQQIEwJDQTEVMBMGA1UEBxMMUmVkd29vZCBDaXR5MTAwLgYDVQQJEyczIExh
Z29vbiBEcml2ZSwgU3VpdGUgMzAwLCBSZWR3b29kIENpdHkxDjAMBgNVBBETBTk0
MDY1MRYwFAYDVQQKEw1BbGF0aW9uLCBJbmMuMRQwEgYDVQQLEwtFbmdpbmVlcmlu
ZzEXMBUGA1UEAxMOZGFyay1jaGFyaXphcmQwWTATBgcqhkjOPQIBBggqhkjOPQMB
BwNCAAQVWOyffe/4vGtnHoSo6m0juLUxlpslfUx4986ID0KFD1X3b/ttqGPV27Jh
q8osyXhtNyBbKpCXmKb/o22jbcFlo0wwSjAOBgNVHQ8BAf8EBAMCB4AwHQYDVR0l
BBYwFAYIKwYBBQUHAwEGCCsGAQUFBwMCMBkGA1UdEQQSMBCCDmRhcmstY2hhcml6
YXJkMAoGCCqGSM49BAMCA0gAMEUCIQCxh7dUGmlAjMQI09Ohjznay8HtYfkzmAh0
yblQmapdJAIgVsueyY6eWdFGboRG1Yk3Hu3wtnezI2urjw/iEVAE39M=
-----END CERTIFICATE-----
```

```
Certificate:
    Data:
        Version: 3 (0x2)
        Serial Number:
            5f:6d:f6:4a:0f:0a:0b:f8:43:f1:04:c9:db:ee:b9:9a
   Signature Algorithm: ecdsa-with-SHA256
        Issuer: 
        Validity
            Not Before: May 21 00:00:41 2020 GMT
            Not After : May 21 00:00:41 2021 GMT
        Subject: C=US, ST=CA, L=Redwood City/street=3 Lagoon Drive, Suite 300, Redwood City/postalCode=94065, O=Alation, Inc., OU=Engineering, CN=dark-charizard
        Subject Public Key Info:
            Public Key Algorithm: id-ecPublicKey
                Public-Key: (256 bit)
                pub: 
                    04:15:58:ec:9f:7d:ef:f8:bc:6b:67:1e:84:a8:ea:
                    6d:23:b8:b5:31:96:9b:25:7d:4c:78:f7:ce:88:0f:
                    42:85:0f:55:f7:6f:fb:6d:a8:63:d5:db:b2:61:ab:
                    ca:2c:c9:78:6d:37:20:5b:2a:90:97:98:a6:ff:a3:
                    6d:a3:6d:c1:65
                ASN1 OID: prime256v1
                NIST CURVE: P-256
        X509v3 extensions:
            X509v3 Key Usage: critical
                Digital Signature
            X509v3 Extended Key Usage: 
                TLS Web Server Authentication, TLS Web Client Authentication
            X509v3 Subject Alternative Name: 
                DNS:dark-charizard
    Signature Algorithm: ecdsa-with-SHA256
         30:45:02:21:00:b1:87:b7:54:1a:69:40:8c:c4:08:d3:d3:a1:
         8f:39:da:cb:c1:ed:61:f9:33:98:08:74:c9:b9:50:99:aa:5d:
         24:02:20:56:cb:9e:c9:8e:9e:59:d1:46:6e:84:46:d5:89:37:
         1e:ed:f0:b6:77:b3:23:6b:ab:8f:0f:e2:11:50:04:df:d3
```

#### Certificate Expiration
By default, certificates should be valid for relatively short periods of time (90 days or so). In order to ease the pain that this can cause in production (surprise outages due to missed expirations), the ACM agent runs a coroutine that is set to fire off several days before the certificate is set to expire and requests a new issuance.

#### Private Key Rotation Propagation
A change of the issuance secret triggers a re-key of the certificate authority. All ACM agents periodically (and upon bootup) request a copy of the public certificate of the CA to ensure that their own certificate still verifies. If the ACM agent fails to verify itself, then a re-key of the CA is assumed, thus triggering a re-key of the agent itself.

#### Private Key Encryption
The private key for both the certificate authority as well as all ACM agents are stored as an encrypted [PKCS#8](https://en.wikipedia.org/wiki/PKCS_8) where the encrypting key is the result of feeding the pre-shared issuance secret through a key derivation function ([KDF](https://en.wikipedia.org/wiki/Key_derivation_function)) in the presence of a cryptographically secure [nonce](https://en.wikipedia.org/wiki/Cryptographic_nonce) that is unique among all agents. The encryption algorithm is AES-256-CBC.

The following is a concrete example of this file output.

```
-----BEGIN PRIVATE KEY-----
Proc-Type: 4,ENCRYPTED
DEK-Info: AES-256-CBC,071a7fc9c6d5ae6b0f9cac5bf29195db

P+3nkHNT/cX+3sK0Gu3lp5AdCi8LzU9HupLEdBahpg7PK9VUbGauPLdZCOR4AoId
N5LTAk+24pTpTMHUyfXnT3EHS9o0tZQMAGpYLfsTONrKRXdHb1H/6RLp4E5/P1+h
tbVo71iegzObtDf+jqJe+8eOxk3OHp84EvHtkWIZaWLQtn2EpnlAQETBKg91AObe
-----END PRIVATE KEY-----
```
