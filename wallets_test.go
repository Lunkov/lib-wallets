package wallets

import (
  "testing"
  "github.com/stretchr/testify/assert"
)

func TestWallets(t *testing.T) {
  w1 := NewWallet(TypeWalletHD)
  w1.Create(&map[string]string{"mnemonic": "chase oil pigeon elegant ketchup whip frozen beauty unknown brass amount slender pony pottery attitude flavor rifle primary beach sign glue oven crazy lottery"})

  w2 := NewWallet(TypeWalletHD)
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
}
