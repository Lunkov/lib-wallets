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
  mapTypes   map[string]TypeWallet
}

func NewTypesWallet() (*TypesWallet) {
  return &TypesWallet{
      mapTypes  : map[string]TypeWallet{
                                        /* https://goethereumbook.org/hd-wallet/ */
                                        "hd": TypeWallet{
                                                Code: "hd",
                                                Name: "HD Wallet",
                                                Description: "Coins: BTC, LTC, DOGE, DASH, ETH, ETC, BCH, QTUM, USDT, IOST, USDC, TRX, BNB(Binance Chain), FIL",
                                                Coins: []string{"BTC", "ETH", "ETC", "USDT", "ECOS"},
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

func (t *TypesWallet) FindCodeByName(code string) (string, bool) {
  for k, v := range t.mapTypes {
    if code == v.Name {
      return k, true
    }
  }
  return "", false
}

func (t *TypesWallet) GetName(code string) (string) {
  l, ok := t.mapTypes[code]
  if ok {
    return l.Name
  }
  return ""
}

func (t *TypesWallet) Get(code string) (*TypeWallet) {
  l, ok := t.mapTypes[code]
  if ok {
    return &l
  }
  return nil
}
