package nocoin

import (
	"fmt"
	"sync"
)

var maxBlockSize int = 10000000

func NewBlock() *Block {
	return &Block{height: latestBlockHeight + 1, size: 0}
}

type Block struct {
	sync.Mutex

	// block height
	height int

	// transactions
	txs []*Tx

	// block size
	size int
}

// Should there be a min size for the the number of transactions in a block?
// I guess while we are mining the first block we may have 1 transaction but
// as we mine tx should come in. So as long as the block stay below a certain size
// then we should be good
// QUESTION: Should there be a limit to the size of the transactions?
func (block *Block) CollectTx() {
	for _, tx := range txPool {
		// check we haven't reach the block size
		if block.size >= maxBlockSize {
			return
		}
		// validate again as Uxto pool may have changed since this TX was last validated
		_, err := tx.Validate()
		if err != nil {
			fmt.Println("tx not valid not adding to block")
			continue
		}
		block.PutTx(tx)
	}
}

func (block *Block) PutTx(tx *Tx) {
	// remove Tx from uxto pool because we are
	block.Lock()
	defer block.Unlock()
	block.txs = append(block.txs, tx)
	block.size += len(tx.String()) * 4
}

// Create hash of block header
func (block *Block) Hash() {

}

func (block *Block) Header() {

}

func (block *Block) Mine() {
	// check if block height is 0
	// add coin base transaction to block
	// to give first miner (centralized vendor) the entire market cap
	// this token generation will only even happen once
}
