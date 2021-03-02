package nocoin

func AddToUtxoPool(utxo *Utxo) {
	if _, ok := utxoPool[utxo.addr]; !ok {
		utxoPool[utxo.addr] = make([]*Utxo, 0)
	}
	utxoPool[utxo.addr] = append(utxoPool[utxo.addr], utxo)
}

func RemoveFromUtxoPool() {}

func FindInUtxoPool(addr string) []*Utxo {
	return utxoPool[addr]
}

func FindInUtxoPoolSumValue(addr string, amount int) ([]*Utxo, int) {
	utxos := FindInUtxoPool(addr)
	sum := 0
	result := make([]*Utxo, 0)
	for _, utxo := range utxos {
		if sum >= amount {
			break
		}
		sum += utxo.amount
		result = append(result, utxo)
	}
	return result, sum
}

type Utxo struct {
	addr   string
	amount int
}
