package nocoin

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"io/ioutil"
	"log"
	"os"
)

var curve elliptic.Curve = elliptic.P256()

type Addr struct {
	pem    string
	pemPub string
}

func (addr *Addr) generate() {
	if _, err := os.Stat("private.pem"); err == nil {
		log.Print("You are about to overwrite private.pem. You probably don't want to do that.")
		return
	}

	privateKey, _ := ecdsa.GenerateKey(curve, rand.Reader)
	publicKey := &privateKey.PublicKey

	addr.pem, addr.pemPub = encodePem(privateKey, publicKey)
	savePem(addr.pem, addr.pemPub)
}

func (addr *Addr) pubKeyToHexStr() string {
	_, publicKey := decodePem(addr.pem, addr.pemPub)
	bytes := elliptic.MarshalCompressed(curve, publicKey.X, publicKey.Y)
	return hex.EncodeToString(bytes)
}

// func (addr *Addr) pubKeyHash() string {}

func (addr *Addr) loadFromFile() {
	pub_data, err := ioutil.ReadFile("public.pem")
	if err != nil {
		log.Print("Unable to read file public.pem")
		return
	}

	priv_data, err := ioutil.ReadFile("private.pem")
	if err != nil {
		log.Print("Unable to read file public.pem")
		return
	}

	addr.pem = string(priv_data)
	addr.pemPub = string(pub_data)
}

func (addr *Addr) sign(hash [sha256.Size]byte) ([]byte, error) {
	privateKey, _ := decodePem(addr.pem, addr.pemPub)
	sig, err := ecdsa.SignASN1(rand.Reader, privateKey, hash[:])
	return sig, err
}

func verifyPublicKey(publicKey *ecdsa.PublicKey, hash, sig []byte) bool {
    valid := ecdsa.VerifyASN1(publicKey, hash[:], sig)
    return valid
}

func hexStrToPubKey(str string) *ecdsa.PublicKey {
	bytes, _ := hex.DecodeString(str)
	X, Y := elliptic.UnmarshalCompressed(curve, bytes)
	return &ecdsa.PublicKey{curve, X, Y}
}

func savePem(pem, pemPub string) {
	// Save private key
	f_priv, err := os.Create("private.pem")
	if err != nil {
		log.Print("Unable to create private.pem")
		return
	}
	defer f_priv.Close()
	f_priv.WriteString(pem)
	f_priv.Sync()

	// Save public key
	f_pub, err := os.Create("public.pem")
	if err != nil {
		log.Print("Unable to create public.pem")
		return
	}
	defer f_pub.Close()
	f_pub.WriteString(pemPub)
	f_pub.Sync()
}

func encodePem(privateKey *ecdsa.PrivateKey, publicKey *ecdsa.PublicKey) (string, string) {
	x509Encoded, _ := x509.MarshalECPrivateKey(privateKey)
	pemEncoded := pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: x509Encoded})

	x509EncodedPub, _ := x509.MarshalPKIXPublicKey(publicKey)
	pemEncodedPub := pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: x509EncodedPub})

	return string(pemEncoded), string(pemEncodedPub)
}

func decodePem(pemEncoded string, pemEncodedPub string) (*ecdsa.PrivateKey, *ecdsa.PublicKey) {
	block, _ := pem.Decode([]byte(pemEncoded))
	x509Encoded := block.Bytes
	privateKey, _ := x509.ParseECPrivateKey(x509Encoded)

	blockPub, _ := pem.Decode([]byte(pemEncodedPub))
	x509EncodedPub := blockPub.Bytes
	genericPublicKey, _ := x509.ParsePKIXPublicKey(x509EncodedPub)
	publicKey := genericPublicKey.(*ecdsa.PublicKey)

	return privateKey, publicKey
}
