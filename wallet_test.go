package wallets

import (
  "testing"
  "github.com/stretchr/testify/assert"
  
  "github.com/Lunkov/go-hdwallet"
)

func TestWalletHD(t *testing.T) {
  w1 := NewWallet(TypeWalletHD)
  w1.Create(&map[string]string{"mnemonic": "chase oil pigeon elegant ketchup whip frozen beauty unknown brass amount slender pony pottery attitude flavor rifle primary beach sign glue oven crazy lottery"})

  w2 := NewWallet(TypeWalletHD)
  w2.Create(&map[string]string{"mnemonic": "fall farm prepare palm sign city analyst liquid orange naive hire lawn marble object old cradle exchange visa caught base robot online undo possible"})

  addr1 := w1.GetAddress(hdwallet.ECOS)
  addr2 := w2.GetAddress(hdwallet.ECOS)
  
  assert.Equal(t, "0x5fAD534AadacBe64E43944CAEfAC04B087B75F9D", addr1)
  assert.Equal(t, "0x8E348b74e2f52f1c97ADBf0aA42b4c3FC7961fA6", addr2)
  
}
