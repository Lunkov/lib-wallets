package wallets

import (
  "errors"
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

func (w *WalletHD) Create(prop *map[string]string) error {
  mnemonic, ok := (*prop)["mnemonic"]
  if !ok {
    return errors.New("Mnemonic phrase does not exists")
  }
  seed, err := hdwallet.NewSeed(mnemonic, "", hdwallet.English)
  if err != nil {
    return err
  }
  w.Mnemonic = mnemonic
  w.Master, err = hdwallet.NewKey(false, hdwallet.Seed(seed))
  return err
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
  
func (w *WalletHD) Serialize() ([]byte, error) {
  we := w.Export()
  var buff bytes.Buffer
  encoder := gob.NewEncoder(&buff)
  encoder.Encode(we)
  return buff.Bytes(), nil
}

func (w *WalletHD) Import(we WalletExport) error {
  w.Name = we.Name
  w.Type = we.Type
  w.Mnemonic = we.Secret

  seed, err := hdwallet.NewSeed(w.Mnemonic, "", hdwallet.English)
  if err != nil {
    return err
  }
  w.Master, err = hdwallet.NewKey(false, hdwallet.Seed(seed))
  if err != nil {
    return err
  }

  return nil
}

func (w *WalletHD) Deserialize(buffer []byte) error {
  var we WalletExport
  buf := bytes.NewBuffer(buffer)
  decoder := gob.NewDecoder(buf)
  err := decoder.Decode(&we)
  if err != nil {
    return err
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

func (w *WalletHD) Save2Folder(pathname string, password string) error {
  cf := cipher.NewCFile()
  filename := pathname + string(os.PathSeparator) + calcMD5Hash(w.GetAddress(hdwallet.ECOS)) + ".wallet"
  buf, _ := w.Serialize()
  return cf.SaveFilePwd(filename, password, buf)
}

func (w *WalletHD) SaveFile(filename string, password string) error {
  cf := cipher.NewCFile()
  buf, _ := w.Serialize()
  return cf.SaveFilePwd(filename, password, buf)
}

func (w *WalletHD) LoadFile(filename string, password string) error {
  cf := cipher.NewCFile()
  buf, err := cf.LoadFilePwd(filename, password)
  if err != nil {
    return err
  }
  return w.Deserialize(buf)
}
