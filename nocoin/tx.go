package nocoin

import (
	"log"
)

// func TxFromString(txStr string) *Tx {}

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

	tx := &Tx { vin: vin, vout: vout, amount: amount, addr: addr }
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

	// amount to transfer
	amount int

	// address to transfer amount to OR address of contract
	addr string
}

// Coin base is first transaction in a block used to pay a reward
// there is only 1 reward of the entire market cap for the first mined
// block in nocoin int64 9223372036854775807
func (tx *Tx) IsCoinBase() bool {
	// TODO:
	return false
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

func (tx *Tx) ValidateTx() bool {
	if tx.IsCoinBase() && latestBlockHeight != 0 {
		log.Printf("coin base transaction only valid in block height: 0\n")
		return false
	}
	// TODO:
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