package nocoin

import (
	"log"
	"strings"
	"strconv"
)

var marketCap int64 = 9223372036854775807

// Take a byte array and generates a TX structure
// shape of broadcast transaction is:
// <sig><sender:pubkey><recv:addr><vin>
// func GenerateTxFromBytes(bytes []byte) *Tx {}

// Generate a transfer TX from a string input
// Example:
// <amount> <recv:addr>
// eg: 20 D80C9BF910F144738EF983724BC04BD6BD3F17C5C83ED57BEDEE1B1B9278E811
func GenerateTransferTxFromString(input string) *Tx {
	parts := strings.Split(" ", input)
	if len(parts) > 2 {
		log.Print("Too many inputs. Expect format <amount> <addr>")
		return
	}
	amount, err := strconv.Atoi(parts[0])
	if err != nil {
		log.Print("First input must be a number")
		return
	}
	addr := parts[1]

	tx := &Tx{}
	tx.BuildTransfer(amount, addr)
	return tx
}

// There are three types of transaction:
// - TX transfer
// - TX function deploy 
// - TX function call
//
// shape of transaction is something like:
// TX: 		<sig><sender:pubkey><recv:addr><vin><payload><vout><fnRet>
// VIN:		[]<sender:addr><amount>
// VOUT:	[]<recv:addr><amount>
// PAYLOAD:	[]<args>
// fnRet:	<value> (format::json)
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

// Build the transfer TX
// Based on an amount and addr we dip into the UXTO pool and find transactions
// that we are able to spend which will form our VIN. Then we will contruct
// out VOUT based on the address we are sending to and if we need to send back 
// some change
func (tx *Tx) BuildTransfer(amount int, addr string) {}