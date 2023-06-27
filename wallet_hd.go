package wallets

import (
  "os"
  "bytes"
  "encoding/gob"
  "crypto/ecdsa"

  "github.com/Lunkov/go-hdwallet"
  
  "github.com/Lunkov/lib-cipher"
)

// https://dev.to/nheindev/building-a-blockchain-in-go-pt-v-wallets-12na


type WalletHD struct {
  Name            string          `yaml:"name"`
  Type            uint32          `yaml:"type"`
  Path            string          `yaml:"path"`
  Loaded          bool            `yaml:"-"`
  Pass            string          `yaml:"-"`
  
  Mnemonic        string          `yaml:"-"`
  Master         *hdwallet.Key    `yaml:"-"`
}

func newWalletHD() IWallet {
  return &WalletHD{Type: TypeWalletHD}
}

func (w *WalletHD) SetName(name string) { w.Name = name }
func (w *WalletHD) GetName() string     { return w.Name }

func (w *WalletHD) SetType(t uint32) { w.Type = t }
func (w *WalletHD) GetType() uint32  { return w.Type }

func (w *WalletHD) SetPath(p string) { w.Path = p + string(os.PathSeparator) + calcMD5Hash(w.GetAddress(hdwallet.ECOS))  }
func (w *WalletHD) GetPath() string  { return w.Path }

func (w *WalletHD) Create(prop *map[string]string) bool {
  mnemonic, ok := (*prop)["mnemonic"]
  if !ok {
    return false
  }
  seed, err := hdwallet.NewSeed(mnemonic, "", hdwallet.English)
  if err != nil {
    return false
  }
  w.Mnemonic = mnemonic
  w.Master, _ = hdwallet.NewKey(false, hdwallet.Seed(seed))
  return w.Master != nil
}

func (w *WalletHD) GetECDSAPrivateKey() *ecdsa.PrivateKey {
  if w.Master == nil {
    return nil
  }
  return w.Master.PrivateECDSA
}

func (w *WalletHD) GetECDSAPublicKey()  *ecdsa.PublicKey {
  if w.Master == nil {
    return nil
  }
  return w.Master.PublicECDSA
}

func (w *WalletHD) Export() WalletExport {
  return WalletExport{Name: w.Name, Type: w.Type, Secret: w.Mnemonic}
}
  
func (w *WalletHD) ExportBuf() []byte {
  we := w.Export()
  var buff bytes.Buffer
  encoder := gob.NewEncoder(&buff)
  encoder.Encode(we)
  return buff.Bytes()
}

func (w *WalletHD) Import(we WalletExport) bool {
  w.Name = we.Name
  w.Type = we.Type
  w.Mnemonic = we.Secret

  seed, err := hdwallet.NewSeed(w.Mnemonic, "", hdwallet.English)
  if err != nil {
    return false
  }
  w.Master, _ = hdwallet.NewKey(false, hdwallet.Seed(seed))

  return true
}

func (w *WalletHD) ImportBuf(buffer []byte) bool {
  var we WalletExport
  buf := bytes.NewBuffer(buffer)
  decoder := gob.NewDecoder(buf)
  err := decoder.Decode(&we)
  if err != nil {
    return false
  }
  return w.Import(we)
}

func (w *WalletHD) GetAddress(coin uint32) string {
  var address string
  switch coin {
    case hdwallet.USDT:
         wallet, _ := w.Master.GetWallet(hdwallet.CoinType(hdwallet.USDT), hdwallet.AddressIndex(1))
         address, _ = wallet.GetAddress()
         break
    case hdwallet.BTC:
         wallet, _ := w.Master.GetWallet(hdwallet.CoinType(hdwallet.BTC), hdwallet.AddressIndex(1))
         address, _ = wallet.GetAddress()
         break
    case hdwallet.ETH, hdwallet.ETC, hdwallet.ECOS, hdwallet.EVER:
         wallet, _ := w.Master.GetWallet(hdwallet.CoinType(coin))
         address, _ = wallet.GetAddress()
         break
  }
  return address
}

func (w *WalletHD) Save(pathname string, password string) bool {
  cf := cipher.NewCFile()
  filename := pathname + string(os.PathSeparator) + calcMD5Hash(w.GetAddress(hdwallet.ECOS)) + ".wallet"
  return cf.SaveFilePwd(filename, password, w.ExportBuf())
}

func (w *WalletHD) Load(filename string, password string) bool {
  cf := cipher.NewCFile()
  buf, ok := cf.LoadFilePwd(filename, password)
  if !ok {
    return ok
  }
  return w.ImportBuf(buf)
}
