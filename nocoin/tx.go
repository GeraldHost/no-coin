package nocoin

// Take a byte array and generates a TX structure
// shape of broadcast transaction is:
// <sig><sender:pubkey<recv:addr><vin>
func GenerateTxFromBytes(bytes []byte) *Tx {}

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
// block in nocoin
func (tx *Tx) isCoinBase() bool {
  return tx.in == 0 
}

func (tx *Tx) getHash() string {
  return 
}

func (tx *Tx) AddToMemPool() {
  hash := tx.getHash()
  txPool[hash] = tx
}

func (tx *Tx) RemoveFromMemPool() {
  hash := tx.getHash()
  if _, ok := txPool[hash]; ok {
    delete(txPool, hash);
  }
}

func (tx *Tx) ValidateTx() bool {
  if tx.isCoinBase() && latestBlockHeight !== 0 {
    log.Printf("Coin base transaction only valid in block height: 0\n")
    return
  }
  // check transaction inputs
  // checkout outputs
  // validate sigs and pub keys
  return false
}

