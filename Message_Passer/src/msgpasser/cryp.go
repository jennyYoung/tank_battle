package msgpasser

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"encoding/pem"
	"hash"
	"log"
	"bytes"
)

type CryptoTool struct {
	privateKey *rsa.PrivateKey
	publicKey *rsa.PublicKey
	sha1 hash.Hash
}

func (ct *CryptoTool) Init() {
	ct.sha1 = sha1.New()
}

func (ct *CryptoTool) GenerateKey() {
	bitSize := 512
	var err error
	ct.privateKey, err = rsa.GenerateKey(rand.Reader, bitSize)
	if err != nil {
		log.Fatalln("Generate key fail")
	}
	ct.publicKey = &ct.privateKey.PublicKey
}

func (ct *CryptoTool) Enc(in []byte) []byte {
	out, err := rsa.EncryptOAEP(ct.sha1, rand.Reader, ct.publicKey, in, nil)
	if err != nil {
		log.Fatalln("EncryptOAEP fail", err)
	}
	return out
}

func (ct *CryptoTool) Dec(in []byte) []byte {
	if ct.privateKey == nil {
		log.Fatalln("No private key?!")
	}
	out, err := rsa.DecryptOAEP(ct.sha1, rand.Reader, ct.privateKey, in, nil)
	if err != nil {
		log.Fatalln("DecryptOAEP fail", err)
	}
	return out
}

func (ct *CryptoTool) MarshalPublicKey() []byte {
	stream, err := x509.MarshalPKIXPublicKey(ct.publicKey)
	if err != nil {
		log.Fatalln("Marshal public key error")
	}
	block := &pem.Block{Type: "PUBLIC KEY", Bytes: stream}
	var b bytes.Buffer
	err = pem.Encode(&b, block)
	if err != nil {
		log.Fatalln("encode on block error")
	}
	return b.Bytes()
}

func (ct *CryptoTool) ParsePublicKey(b []byte) {
	block, _ := pem.Decode(b)
	if block == nil {
		log.Fatalln("parse public key error")
	}
	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		log.Fatalln("parse2 public key error")
	}
	ct.publicKey = pubInterface.(*rsa.PublicKey)
	if ct.publicKey == nil {
		log.Fatalln("parse3 public key error")
	}
}
