package nocoin

import (
  "testing"
)

func TestTxValidation(b *testing.T) {
	SetupAddr()
	SetupVendor()

    tx := NewTxTransfer(20, "D80C9BF910F144738EF983724BC04BD6BD3F17C5C83ED57BEDEE1B1B9278E811")
    signedTx := tx.SignTx()

    sig, tx := TxFromString(signedTx)
    // validate TX
    if valid := tx.Validate(sig); valid {
      tx.AddToMemPool()
    }
}
