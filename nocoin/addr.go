package nocoin 

import (
    "os"
	"log"
    "io/ioutil"
	"crypto/rand"
	"crypto/ecdsa"
	"crypto/x509"
    "crypto/elliptic"
	"encoding/pem"
)

type Addr struct {
    pem string
    pemPub string
}

func (addr *Addr) generate() {
    if _, err := os.Stat("private.pem"); err == nil {
        log.Print("You are about to overwrite private.pem. You probably don't want to do that.")
        return
    }

    privateKey, _ := ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
    publicKey := &privateKey.PublicKey

    addr.pem, addr.pemPub = encodePem(privateKey, publicKey)
    savePem(addr.pem, addr.pemPub)
}

func (addr *Addr) load() {
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

// func (addr *Addr) sign() {}

// func (addr *Addr) verify() {}

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