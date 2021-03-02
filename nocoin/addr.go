package nocoin

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"os"
)

var curve elliptic.Curve = elliptic.P256()

var addrLength int = 64

// Address are going to basically just be hashes of public
// keys. When we send money to a hash we are actually just
// creating an output that can be spent by somebody who has
// a valid public key that hashes into the same address hash
// used in the previous output.
//
// We use ecdsa which is just a digital signature algo. That allows
// use to sign and verify requests with public and private keys.
// We can generate pub/priv keys by calling addr.generate(). The keys
// are stored locally on disk in public.pem and private.pem
//
// Public keys are stored as PEM encoded but they are sent in a
// transaction compressed as a hex string.
type Addr struct {
	pem    string
	pemPub string
}

func (addr *Addr) generate() {
	if _, err := os.Stat("private.pem"); err == nil {
		fmt.Println("you are about to overwrite private.pem you probably don't want to do that")
		return
	}

	privateKey, _ := ecdsa.GenerateKey(curve, rand.Reader)
	publicKey := &privateKey.PublicKey

	addr.pem, addr.pemPub = encodePem(privateKey, publicKey)
	savePem(addr.pem, addr.pemPub)
}

// Convert pem encoded public key to hex string
func (addr *Addr) PubKeyToHexStr() string {
	_, publicKey := decodePem(addr.pem, addr.pemPub)
	bytes := elliptic.MarshalCompressed(curve, publicKey.X, publicKey.Y)
	return hex.EncodeToString(bytes)
}

func (addr *Addr) Get() string {
	return addr.PubKeyHash()
}

// Convert hex encoded public key to SHA256 hash
func (addr *Addr) PubKeyHash() string {
	pubKeyStr := addr.PubKeyToHexStr()
	hash := sha256.Sum256([]byte(pubKeyStr))
	return fmt.Sprintf("%x", hash)
}

func (addr *Addr) LoadFromFile() bool {
	pub_data, err := ioutil.ReadFile("public.pem")
	if err != nil {
		fmt.Println("unable to read file public.pem")
		return false
	}

	priv_data, err := ioutil.ReadFile("private.pem")
	if err != nil {
		fmt.Println("unable to read file public.pem")
		return false
	}

	addr.pem = string(priv_data)
	addr.pemPub = string(pub_data)
	return true
}

func (addr *Addr) Sign(hash [sha256.Size]byte) ([]byte, error) {
	privateKey, _ := decodePem(addr.pem, addr.pemPub)
	sig, err := ecdsa.SignASN1(rand.Reader, privateKey, hash[:])
	return sig, err
}

func verifyPublicKey(publicKey *ecdsa.PublicKey, hash, sig []byte) bool {
	valid := ecdsa.VerifyASN1(publicKey, hash[:], sig)
	return valid
}

// Convert hex string encoded public key to ecdsa public key
func hexStrToPubKey(str string) *ecdsa.PublicKey {
	bytes, _ := hex.DecodeString(str)
	X, Y := elliptic.UnmarshalCompressed(curve, bytes)
	return &ecdsa.PublicKey{curve, X, Y}
}

func savePem(pem, pemPub string) {
	// Save private key
	f_priv, err := os.Create("private.pem")
	if err != nil {
		fmt.Println("unable to create private.pem")
		return
	}
	defer f_priv.Close()
	f_priv.WriteString(pem)
	f_priv.Sync()

	// Save public key
	f_pub, err := os.Create("public.pem")
	if err != nil {
		fmt.Println("unable to create public.pem")
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
