package wallets

import (
  "bytes"
  "crypto/ecdsa"
  "crypto/md5"
  "encoding/hex"
  "encoding/gob"
  
  "github.com/golang/glog"

  "github.com/Lunkov/lib-cipher"
)

type IWallet interface {
  SetName(name string)
  GetName() string
  SetType(t string)
  GetType() string
  SetPath(p string)
  GetPath() string
  
  Create(prop *map[string]string) bool

  Load(filename string, password string) bool
  Save(pathname string, password string) bool
  
  Export() []byte
  Import(buffer []byte) bool
  
  GetAddress(coin string) string

  GetECDSAPrivateKey() *ecdsa.PrivateKey
  GetECDSAPublicKey()  *ecdsa.PublicKey
}

func NewWallet(t string) IWallet {
  switch t {
    case "HD":
         w := newWalletHD()
         w.SetType(t)
         return w
  }
  return nil
}


type WalletExport struct {
  Name          string   `yaml:"name"`
  Type          string   `yaml:"type"`
  Public        string   `yaml:"public"`
  Secret        string   `yaml:"secret"`
}

type EmptyWallet struct {
  Name            string   `yaml:"name"`
  Type            string   `yaml:"type"`
}

func NewEmptyWallet() *EmptyWallet {
  return &EmptyWallet{}
}

func (w *EmptyWallet) Load(filename string, password string) bool {
  cf := cipher.NewCFile()
  buf, ok := cf.LoadFilePwd(filename, password)
  if !ok {
    return ok
  }
  return w.Deserialize(buf)
}

func (w *EmptyWallet) Deserialize(buffer []byte) bool {
  buf := bytes.NewBuffer(buffer)
  decoder := gob.NewDecoder(buf)
  err := decoder.Decode(w)
  if err != nil {
    glog.Errorf("ERR: gob.NewDecoder: %v", err)
    return false
  }
  return true
}

func calcMD5Hash(text string) string {
  hasher := md5.New()
  hasher.Write([]byte(text))
  return hex.EncodeToString(hasher.Sum(nil))
}
