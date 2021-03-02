package nocoin

import (
	"fmt"
	"bytes"
)

func TxPartsFromReader(r *bytes.Buffer) []*TxPart {
	txParts := make([]*TxPart, 0)
	count, _ := BytesToInt(r.Next(2))
	for i := 0; i < int(count); i++ {
		maybePrefix := r.Next(2)
		var amount int64
		if n, ok := varIntPrefixes[string(maybePrefix)]; ok {
			amount, _ = BytesToInt(r.Next(n))
		} else {
			amount, _ = BytesToInt(maybePrefix)
		}
		addr := string(r.Next(addrLength))

		txPart := &TxPart{ amount: int(amount), addr: addr }
		txParts = append(txParts, txPart)	
	}
	return txParts
}

func TxFromString(txStr string) *Tx {
	r := bytes.NewBuffer([]byte(txStr))
	vin := TxPartsFromReader(r)
	vout := TxPartsFromReader(r)
	return &Tx { vin: vin, vout: vout }
}

// Return vin and vout for transaction
// <vin> <amount><address>
// <vout> <amount><address>
// TODO move to tx.go
func NewTxTransfer(amount int, addr string) *Tx {
	myAddrAddr := myAddr.Get()
	vin := make([]*TxPart, 0)
	vout := make([]*TxPart, 0)
	// Build inputs
	utxos, sum := FindInUtxoPoolSumValue(myAddrAddr, amount)
	for _, utxo := range utxos {
		vin = append(vin, &TxPart{ amount: utxo.amount, addr: utxo.addr })
	}
	// append transfer to vout
	vout = append(vout, &TxPart{ amount: amount, addr: addr });
	// check if change is required
	if sum > amount {
		// The sum of utxos is greater than the amount so we need some change
		change := sum - amount
		vout = append(vout, &TxPart{ amount: change, addr: myAddrAddr })
	}

	tx := &Tx { vin: vin, vout: vout }
	return tx
}

// Take a byte array and generates a TX structure
// shape of broadcast transaction is:
// <sig><sender:pubkey><recv:addr><vin>
// func GenerateTxFromBytes(bytes []byte) *Tx {}

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
	vin []*TxPart

	// TX value output (must match input in value)
	vout []*TxPart

	// Payload to be sent in if this TX is a function call
	payload []byte

	// Return value from function call
	fnRet []byte
}

// Return the transaction hash which can be used to check validty when
// we create our merkle tree.
func (tx *Tx) Hash() string {
	// TODO:
	return "empty"
}

// All transactions are kept in a memory pool until they are ready
// to be added to a block.
func (tx *Tx) AddToMemPool() {
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

// TODO:
func (tx *Tx) Validate() bool {
	sum := func(txParts []*TxPart) int {
		s := 0
		for _, txP := range txParts {
			s += txP.amount
		}
		return s
	}
	vinSum := sum(tx.vin)
	voutSum := sum(tx.vout)
	fmt.Println(vinSum, voutSum, vinSum == voutSum)
	// check transaction inputs
	// checkout outputs
	// validate sigs and pub keys
	return false
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
	return EncodeVarInt(len(tx.vin)) + string(vinBytes) + EncodeVarInt(len(tx.vout)) + string(voutBytes)
}


type TxPart struct {
	amount int
	addr string
}

func (txP *TxPart) String() string {
	return EncodeVarInt(txP.amount) + txP.addr
}