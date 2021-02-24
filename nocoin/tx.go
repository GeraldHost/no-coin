package nocoin

type Tx struct {
  // TODO: not acctually sure what the types for these are yet
  vin int
  vout int
  
  // TODO: not acctually sure what the types for these are yet
  sin int
  sout int
	
  // Coin base is first transaction in a block used to pay a reward
  // there is only 1 reward of the entire market cap for the first mined
  // block in nocoin
  isCoinBase bool
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
  // check transaction inputs
  // checkout outputs
  // validate sigs and pub keys
  return false
}

