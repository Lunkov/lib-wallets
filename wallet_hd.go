package wallets

import (
  "bytes"
  "sync"
  "encoding/gob"
  "context"
  "github.com/golang/glog"

  "github.com/Lunkov/go-hdwallet"
  "github.com/Lunkov/lib-messages"
  
  "github.com/ethereum/go-ethereum/common"
  "github.com/ethereum/go-ethereum/ethclient"
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
  return w
}

func (w *WalletHD) Create(prop *map[string]string) bool {
  mnemonic, ok := (*prop)["mnemonic"]
  if !ok {
    return false
  }
  seed, err := hdwallet.NewSeed(mnemonic, "", hdwallet.English)
  if err != nil {
    glog.Errorf("ERR: Wallet.Create: %v", err)
    return false
  }
  w.Mnemonic = mnemonic
  w.Master, _ = hdwallet.NewKey(false, hdwallet.Seed(seed))
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

  seed, err := hdwallet.NewSeed(w.Mnemonic, "", hdwallet.English)
  if err != nil {
    glog.Errorf("ERR: Wallet.Create: %v", err)
    return false
  }
  w.Master, _ = hdwallet.NewKey(false, hdwallet.Seed(seed))

  return true
}

func (w *WalletHD) GetBalance() ([]*messages.Balance) {
  var wg sync.WaitGroup
  b := messages.NewBalances()
  // b["btc"] = "10.00"
  // b["ecos"] = "10000.00"
  wg.Add(1)
  go func() {
    defer wg.Done()
    client, err := ethclient.Dial("https://mainnet.infura.io")
    if err != nil {
      glog.Errorf("ERR: Wallet.ethclient.Dial: %v", err)
      return
    }
    b1 := messages.NewBalance()
    b1.Address = w.GetAddress("ETH")
    b1.Coin = "ETH"
    account := common.HexToAddress(b1.Address)
    balance, errb := client.BalanceAt(context.Background(), account, nil)
    if errb != nil {
      glog.Errorf("ERR: Wallet.ethclient.BalanceAt: %v", errb)
      return
    }
    b1.Balance = balance.Uint64()
    b.Add(b1)
  } ()
  wg.Wait()
  return b.GetBalanses()
}

func (w *WalletHD) GetAddress(coin string) string {
  var address string
  switch coin {
    case "USDT":
         wallet, _ := w.Master.GetWallet(hdwallet.CoinType(hdwallet.USDT), hdwallet.AddressIndex(1))
         address, _ = wallet.GetAddress()
         break
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
    case "ECOS":
         wallet, _ := w.Master.GetWallet(hdwallet.CoinType(hdwallet.ECOS))
         address, _ = wallet.GetAddress()
         break
    case "EVER":
         wallet, _ := w.Master.GetWallet(hdwallet.CoinType(hdwallet.EVER))
         address, _ = wallet.GetAddress()
         break
  }
  return address
}
