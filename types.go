package wallets

import (
)

type TypeWallet struct {
  Code         string
  Name         string
  Description  string
  Coins        []string
}

type TypesWallet struct {
  mapTypes   map[uint32]TypeWallet
}

func NewTypesWallet() (*TypesWallet) {
  return &TypesWallet{
      mapTypes  : map[uint32]TypeWallet{
                                        /* https://goethereumbook.org/hd-wallet/ */
                                        TypeWalletHD: TypeWallet{
                                                Code: "hd",
                                                Name: "HD Wallet",
                                                Description: "Coins: BTC, LTC, DOGE, DASH, ETH, ETC, BCH, QTUM, USDT, IOST, USDC, TRX, BNB(Binance Chain), FIL",
                                                Coins: []string{"BTC", "ETH", "USDT", "ECOS", "EVER"},
                                              },
                                      },
  }
}


func (t *TypesWallet) GetCodes() ([]string) {
  keys := make([]string, len(t.mapTypes))

  i := 0
  for _, v := range t.mapTypes {
    keys[i] = v.Name
    i++
  }
  return keys
}

func (t *TypesWallet) FindIdByName(name string) (uint32, bool) {
  for k, v := range t.mapTypes {
    if name == v.Name {
      return k, true
    }
  }
  return 0, false
}

func (t *TypesWallet) GetName(id uint32) (string) {
  l, ok := t.mapTypes[id]
  if ok {
    return l.Name
  }
  return ""
}

func (t *TypesWallet) Get(id uint32) (*TypeWallet) {
  l, ok := t.mapTypes[id]
  if ok {
    return &l
  }
  return nil
}
