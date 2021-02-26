package nocoin

import (
	"log"
)

var marketCap int64 = 9223372036854775807

// Take a byte array and generates a TX structure
// shape of broadcast transaction is:
// <sig><sender:pubkey><recv:addr><vin>
// func GenerateTxFromBytes(bytes []byte) *Tx {}

// shape of transaction is:
// <sig><sender:pubkey><recv:addr><vin><payload><vout><fnRet>
type Tx struct {
	// TX value input
	vin []byte

	// TX value output (must match input in value)
	vout []byte

	// Payload to be sent in if this TX is a function call
	payload []byte

	// Return value from function call
	fnRet []byte
}

// Coin base is first transaction in a block used to pay a reward
// there is only 1 reward of the entire market cap for the first mined
// block in nocoin int64 9223372036854775807
func (tx *Tx) IsCoinBase() bool {
	return false
}

// Return the transaction hash which can be used to check validty when
// we create our merkle tree. 
func (tx *Tx) GetHash() string {
	return "empty"
}

// All transactions are kept in a memory pool until they are ready
// to be added to a block. 
func (tx *Tx) AddToMemPool() {
	hash := tx.GetHash()
	txPool[hash] = tx
}

// Once a transaction is added to a block we pull it
// out of the memory pull
func (tx *Tx) RemoveFromMemPool() {
	hash := tx.GetHash()
	if _, ok := txPool[hash]; ok {
		delete(txPool, hash)
	}
}

func (tx *Tx) ValidateTx() bool {
	if tx.IsCoinBase() && latestBlockHeight != 0 {
		log.Printf("Coin base transaction only valid in block height: 0\n")
		return false
	}
	// check transaction inputs
	// checkout outputs
	// validate sigs and pub keys
	return false
}
