package x509

import (
	"bytes"
	"crypto/x509"
	"encoding/gob"
	"encoding/pem"

	"github.com/dvyukov/go-fuzz/examples/fuzz"
)

func FuzzCRL(data []byte) int {
	list, err := x509.ParseCRL(data)
	if err != nil {
		if list != nil {
			panic("list is not nil on error")
		}
		return 0
	}
	return 1
}

func FuzzDERCRL(data []byte) int {
	list, err := x509.ParseDERCRL(data)
	if err != nil {
		if list != nil {
			panic("list is not nil on error")
		}
		return 0
	}
	return 1
}

func FuzzCertificate(data []byte) int {
	c, err := x509.ParseCertificate(data)
	if err != nil {
		if c != nil {
			panic("cert is not nil on error")
		}
		return 0
	}
	c.CheckSignature(x509.SHA1WithRSA, []byte("data"), []byte("01234567890123456789"))
	c.VerifyHostname("host.com")
	pool := x509.NewCertPool()
	pool.AddCert(c)
	c.Verify(x509.VerifyOptions{DNSName: "host.com", Intermediates: pool})
	return 1
}

func FuzzPEM(data []byte) int {
	var b pem.Block
	err := gob.NewDecoder(bytes.NewReader(data)).Decode(&b)
	if err != nil {
		return 0
	}
	b1, err := x509.DecryptPEMBlock(&b, []byte("pass"))
	if err != nil {
		return 0
	}
	b2, err := x509.EncryptPEMBlock(zeroReader(0), "msg", b1, []byte("pass1"), x509.PEMCipherDES)
	if err != nil {
		panic(err)
	}
	_, err = x509.DecryptPEMBlock(b2, []byte("pass"))
	if err == nil {
		panic("decoded with a wrong pass")
	}
	b3, err := x509.DecryptPEMBlock(b2, []byte("pass1"))
	if err != nil {
		panic(err)
	}
	if !bytes.Equal(b1, b3) {
		panic("data changed")
	}
	return 1
}

func FuzzPKIX(data []byte) int {
	key, err := x509.ParsePKIXPublicKey(data)
	if err != nil {
		return 0
	}
	data1, err := x509.MarshalPKIXPublicKey(key)
	if err != nil {
		panic(err)
	}
	key1, err := x509.ParsePKIXPublicKey(data1)
	if err != nil {
		panic(err)
	}
	if !fuzz.DeepEqual(key, key1) {
		panic("keys are not equal")
	}
	return 1
}

type zeroReader int

func (zeroReader) Read(data []byte) (int, error) {
	for i := range data {
		data[i] = byte(i)
	}
	return len(data), nil
}
