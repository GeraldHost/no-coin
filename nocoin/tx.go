package nocoin

import (
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
)

var txIdLength int = 64

func TxPartsFromReader(r *bytes.Buffer) []*TxPart {
	txParts := make([]*TxPart, 0)
	count, _ := BytesToInt(r.Next(2))
	for i := 0; i < int(count); i++ {
		amount := VarIntFromReader(r)
		addr := string(r.Next(addrLength))
		txPart := &TxPart{amount: int(amount), addr: addr}
		txParts = append(txParts, txPart)
	}
	return txParts
}

func TxFromString(txStr string) (string, *Tx) {
	r := bytes.NewBuffer([]byte(txStr))
	sigLen := VarIntFromReader(r)
	sig := r.Next(int(sigLen))
	txHash := r.Next(txIdLength)
	pubKeyLen := VarIntFromReader(r)
	pubKey := r.Next(int(pubKeyLen))
	vin := TxPartsFromReader(r)
	vout := TxPartsFromReader(r)
	return string(sig), &Tx{vin: vin, vout: vout, pubKeyStr: string(pubKey), id: string(txHash)}
}

// Return vin and vout for transaction
// <vin> <amount><address>
// <vout> <amount><address>
func NewTxTransfer(amount int, addr string) *Tx {
	pubKeyStr := myAddr.PubKeyToHexStr()
	myAddrAddr := myAddr.Get()
	vin := make([]*TxPart, 0)
	vout := make([]*TxPart, 0)
	// Build inputs
	utxos, sum := FindInUtxoPoolSumValue(myAddrAddr, amount)
	for _, utxo := range utxos {
		vin = append(vin, &TxPart{amount: utxo.amount, addr: utxo.addr})
	}
	// append transfer to vout
	vout = append(vout, &TxPart{amount: amount, addr: addr})
	// check if change is required
	if sum > amount {
		// The sum of utxos is greater than the amount so we need some change
		change := sum - amount
		vout = append(vout, &TxPart{amount: change, addr: myAddrAddr})
	}

	tx := &Tx{vin: vin, vout: vout, pubKeyStr: pubKeyStr}
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
	// sha256 hash of the tx string
	id string

	// TX value input
	vin []*TxPart

	// TX value output (must match input in value)
	vout []*TxPart

	// Payload to be sent in if this TX is a function call
	payload []byte

	// Return value from function call
	fnRet []byte

	// public key of who is sending the transaction
	pubKeyStr string
}

// Return the transaction hash which can be used to check validty when
// we create our merkle tree.
func (tx *Tx) Hash() string {
	txStr := tx.String()
	return Sha256(txStr)
}

// All transactions are kept in a memory pool until they are ready
// to be added to a block.
func (tx *Tx) AddToMemPool() {
	// TODO: we have to remove the inputs from the uxto pool
	// and put the outputs in the uxto pool
	hash := tx.Hash()
	txPool[hash] = tx
}

// Once a transaction is added to a block we pull it
// out of the memory pull
func (tx *Tx) RemoveFromMemPool() {
	hash := tx.Hash()
	if _, ok := txPool[hash]; ok {
		delete(txPool, hash)
	}
}

// Validate transaction, check utxo pool, check pub keys, check input
// output values are equal, validate signature
func (tx *Tx) Validate(sig string) (bool, error) {
	// validate the sender actually owns the vin credits
	// and that the vin credits are actually in the pool
	senderAddr := Sha256(tx.pubKeyStr)
	for _, txPart := range tx.vin {
		if senderAddr != txPart.addr {
			return false, errors.New("tx inout addr does not match sender, sender can't send credits sender does not own")
		}
		_, err := FindOneInUtxoPool(txPart.addr, txPart.amount)
		if err != nil {
			return false, errors.New("tx input is not in utxo pool")
		}
	}
	pubKey := hexStrToPubKey(tx.pubKeyStr)
	sigB, _ := hex.DecodeString(sig)
	if validSig := verifyPublicKey(pubKey, []byte(tx.id), sigB); !validSig {
		return false, errors.New("public key not valid")
	}
	return true, nil
}

func (tx *Tx) String() string {
	vinBytes := make([]byte, 0)
	for _, txPart := range tx.vin {
		vinBytes = append(vinBytes, []byte(txPart.String())...)
	}
	voutBytes := make([]byte, 0)
	for _, txPart := range tx.vout {
		voutBytes = append(voutBytes, []byte(txPart.String())...)
	}

	return EncodeVarInt(len(tx.pubKeyStr)) + tx.pubKeyStr + EncodeVarInt(len(tx.vin)) + string(vinBytes) + EncodeVarInt(len(tx.vout)) + string(voutBytes)
}

func (tx *Tx) SignTx() string {
	txStr := tx.String()
	tx.id = Sha256(txStr)
	sig, err := myAddr.Sign([]byte(tx.id))
	if err != nil {
		fmt.Println("failed to sign tx")
	}
	sigStr := fmt.Sprintf("%X", sig)
	return EncodeVarInt(len(sigStr)) + sigStr + tx.id + txStr
}

type TxPart struct {
	amount int
	addr   string
}

func (txP *TxPart) String() string {
	return EncodeVarInt(txP.amount) + txP.addr
}
