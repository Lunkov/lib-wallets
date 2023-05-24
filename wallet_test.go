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
  
  assert.Equal(t, "0x5f7ae710cED588D42E863E9b55C7c51e56869963", addr1)
  assert.Equal(t, "0xfa242EE498857ec3C06a2E5E9e37b090807B467a", addr2)
  
}
