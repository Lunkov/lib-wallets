package wallets

import (
  "bytes"
  "crypto/ecdsa"
  "crypto/md5"
  "encoding/hex"
  "encoding/gob"

  "github.com/Lunkov/lib-cipher"
)

const (
  TypeWalletHD     = 1
)


type IWallet interface {
  SetName(name string)
  GetName() string
  SetType(t uint32)
  GetType() uint32
  SetPath(p string)
  GetPath() string
  
  Create(prop *map[string]string) error

  LoadFile(filename string, password string) error
  SaveFile(filename string, password string) error
  Save2Folder(pathname string, password string) error
  
  Export() WalletExport
  Import(e WalletExport) error

  Serialize() ([]byte, error)
  Deserialize(buf []byte) error

  GetAddress(coin uint32) string

  GetECDSAPrivateKey() *ecdsa.PrivateKey
  GetECDSAPublicKey()  *ecdsa.PublicKey
}

func NewWallet(t uint32) IWallet {
  switch t {
    case TypeWalletHD:
         w := newWalletHD()
         w.SetType(t)
         return w
  }
  return nil
}


type WalletExport struct {
  Name          string   `yaml:"name"`
  Type          uint32   `yaml:"type"`
  Public        string   `yaml:"public"`
  Secret        string   `yaml:"secret"`
}

type EmptyWallet struct {
  Name            string   `yaml:"name"`
  Type            uint32   `yaml:"type"`
}

func NewEmptyWallet() *EmptyWallet {
  return &EmptyWallet{}
}

func (w *EmptyWallet) Load(filename string, password string) error {
  cf := cipher.NewCFile()
  buf, err := cf.LoadFilePwd(filename, password)
  if err != nil {
    return err
  }
  return w.Deserialize(buf)
}

func (w *EmptyWallet) Deserialize(buffer []byte) error {
  buf := bytes.NewBuffer(buffer)
  decoder := gob.NewDecoder(buf)
  err := decoder.Decode(w)
  return err
}

func calcMD5Hash(text string) string {
  hasher := md5.New()
  hasher.Write([]byte(text))
  return hex.EncodeToString(hasher.Sum(nil))
}
