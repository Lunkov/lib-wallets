package wallets

import (
  "testing"
  "github.com/stretchr/testify/assert"
  
  "github.com/Lunkov/go-hdwallet"
)

func TestWallets(t *testing.T) {
  w1 := NewWallet(TypeWalletHD)
  w1.SetName("Wallet #1")
  w1.Create(&map[string]string{"mnemonic": "chase oil pigeon elegant ketchup whip frozen beauty unknown brass amount slender pony pottery attitude flavor rifle primary beach sign glue oven crazy lottery"})

  w2 := NewWallet(TypeWalletHD)
  w2.SetName("Wallet #2")
  w2.Create(&map[string]string{"mnemonic": "fall farm prepare palm sign city analyst liquid orange naive hire lawn marble object old cradle exchange visa caught base robot online undo possible"})

  wa := NewWallets()
  wa.Add(w1)
  assert.Equal(t, 1, wa.Count())
  wa.Add(w2)
  assert.Equal(t, 2, wa.Count())
  
  
  wa.Remove(w2)
  assert.Equal(t, 1, wa.Count())
  assert.Equal(t, w1, wa.Get(0))
  assert.Equal(t, nil, wa.Get(1))
  
  list := wa.GetList()
  assert.Equal(t, []string{"Wallet #1 (0x5f7ae710cED588D42E863E9b55C7c51e56869963)"}, list)
  
  
  wa.Add(w2)

  wr, ok := wa.FindByName("name")
  assert.False(t, ok)
  
  wr, ok = wa.FindByName("Wallet #2")
  assert.True(t, ok)
  assert.Equal(t, w2, wr)

  wr, ok = wa.FindByAddress("0x00000")
  assert.False(t, ok)
  
  wr, ok = wa.FindByAddress(w1.GetAddress(hdwallet.ECOS))
  assert.True(t, ok)
  assert.Equal(t, w1, wr)
}
