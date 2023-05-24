package wallets

import (
  "github.com/Lunkov/go-hdwallet"
)

type TypeCoin struct {
  Code         uint32
  Name         string
  Symbol       string
  Description  string
}

type TypesCoin struct {
  mapTypes   map[uint32]TypeCoin
}

func NewTypesCoin() (*TypesCoin) {
  return &TypesCoin{
      mapTypes  : map[uint32]TypeCoin{
                                        hdwallet.ECOS: TypeCoin{
                                                Code: hdwallet.ECOS,
                                                Name: "ECOS",
                                                Symbol: "ecos",
                                                Description: "Coin for Economy Collaboration Space",
                                              },
                                      },
  }
}


func (t *TypesCoin) GetCodes() ([]string) {
  keys := make([]string, len(t.mapTypes))

  i := 0
  for _, v := range t.mapTypes {
    keys[i] = v.Name
    i++
  }
  return keys
}

func (t *TypesCoin) FindCodeByName(name string) (uint32, bool) {
  for k, v := range t.mapTypes {
    if name == v.Name {
      return k, true
    }
  }
  return 0, false
}

func (t *TypesCoin) GetName(code uint32) (string) {
  l, ok := t.mapTypes[code]
  if ok {
    return l.Name
  }
  return ""
}

func (t *TypesCoin) Get(code uint32) (*TypeCoin) {
  l, ok := t.mapTypes[code]
  if ok {
    return &l
  }
  return nil
}
