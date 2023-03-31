package wallets

import (
  "bytes"
  "encoding/gob"

  "github.com/golang/glog"
  "github.com/foxnut/go-hdwallet"
  "github.com/Lunkov/lib-messages"
)

// https://dev.to/nheindev/building-a-blockchain-in-go-pt-v-wallets-12na


type WalletHD struct {
  Wallet
  Mnemonic        string          `yaml:"-"`
  Pass            string          `yaml:"-"`
  Master         *hdwallet.Key    `yaml:"-"`
}

func newWalletHD() IWallet {
  w := &WalletHD{}
  w._export = w.__export
  w._import = w.__import
  //w._get_public_key = w.__get_public_key
  return w
}

func (w *WalletHD) Create(prop *map[string]string) bool {
  mnemonic, ok := (*prop)["mnemonic"]
  if !ok {
    return ok
  }
  w.Mnemonic = mnemonic
  w.Master, _ = hdwallet.NewKey(
    hdwallet.Mnemonic(mnemonic),
  )
  return w.Master != nil
}

func (w *WalletHD) __export() []byte {
  we := WalletExport{Name: w.Name, Type: w.Type, Secret: w.Mnemonic}
  var buff bytes.Buffer
  encoder := gob.NewEncoder(&buff)
  encoder.Encode(we)
  return buff.Bytes()
}

func (w *WalletHD) __import(buffer []byte) bool {
  var we WalletExport
  buf := bytes.NewBuffer(buffer)
  decoder := gob.NewDecoder(buf)
  err := decoder.Decode(&we)
  if err != nil {
    glog.Errorf("ERR: gob.NewDecoder: %v", err)
    return false
  }
  w.Name = we.Name
  w.Type = we.Type
  w.Mnemonic = we.Secret
  w.Master, _ = hdwallet.NewKey(
    hdwallet.Mnemonic(w.Mnemonic),
  )
  return true
}

func (w *WalletHD) GetBalance() ([]*messages.Balance) {
  b := messages.NewBalances()
  // b["btc"] = "10.00"
  // b["ecos"] = "10000.00"
  return b.GetBalanses()
}

func (w *WalletHD) GetAddress(coin string) string {
  var address string
  switch coin {
    case "BTC":
         wallet, _ := w.Master.GetWallet(hdwallet.CoinType(hdwallet.BTC), hdwallet.AddressIndex(1))
         address, _ = wallet.GetAddress()
         break
    case "ETH":
         wallet, _ := w.Master.GetWallet(hdwallet.CoinType(hdwallet.ETH))
         address, _ = wallet.GetAddress()
         break
    case "ETC":
         wallet, _ := w.Master.GetWallet(hdwallet.CoinType(hdwallet.ETC))
         address, _ = wallet.GetAddress()
         break
  }
  return address
}