package nocoin

type Block struct {
  height int
}

func (block *Block) Mine() {
  // check if block height is 0
  // add coin base transaction to block
  // to give first miner (centralized vendor) the entire market cap
  // this token generation will only even happen once
}
