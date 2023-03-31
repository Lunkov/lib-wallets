package wallets

import (
  "bytes"
  "crypto/md5"
  "encoding/hex"
  "encoding/base64"
  "encoding/gob"
  
  "github.com/golang/glog"

  "github.com/Lunkov/lib-messages"
  "github.com/Lunkov/lib-cipher"
)

type IWallet interface {
  GetMD5Hash(text string) string

  SetName(name string)
  GetName() string
  SetType(t string)
  GetType() string
  SetPath(p string)
  GetPath() string
  
  Create(prop *map[string]string) bool

  Load(filename string, password string) bool
  Save(pathname string, password string) bool
  
  Export() string
  Import(str string) bool
  
  GetAddress(coin string) string
  GetBalance() ([]*messages.Balance)
}

type Wallet struct {
  Name            string   `yaml:"name"`
  Type            string   `yaml:"type"`
  Path            string   `yaml:"path"`
  Loaded          bool     `yaml:"-"`
  Pass            string   `yaml:"-"`
  
  _export func() []byte
  _import func(buffer []byte) bool
  //_get_address func(coin string) string
}

type WalletExport struct {
  Name          string   `yaml:"name"`
  Type          string   `yaml:"type"`
  Public        string   `yaml:"public"`
  Secret        string   `yaml:"secret"`
}

func NewWallet(t string) IWallet {
  switch t {
    case "hd":
         w := newWalletHD()
         w.SetType(t)
         return w
    default:
         w := newWallet()
         return w
  }
  return nil
}

func newWallet() IWallet {
  w := &Wallet{}
  w._export = w.__export
  w._import = w.__import
  return w
}

func (w *Wallet) SetName(name string) { w.Name = name }
func (w *Wallet) GetName() string     { return w.Name }

func (w *Wallet) SetType(t string) { w.Type = t }
func (w *Wallet) GetType() string  { return w.Type }

func (w *Wallet) SetPath(p string) { w.Path = p + "/" + w.GetMD5Hash(w.Name)  }
func (w *Wallet) GetPath() string  { return w.Path }

func (w *Wallet) Create(prop *map[string]string) bool {return false}

func (w *Wallet) __export() []byte {
  we := WalletExport{Name: w.Name, Type: w.Type}
  var buff bytes.Buffer
  encoder := gob.NewEncoder(&buff)
  encoder.Encode(we)
  return buff.Bytes()
}

func (w *Wallet) __import(buffer []byte) bool {
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
  return true
}

func (w *Wallet) Export() string {
  return base64.StdEncoding.EncodeToString(w._export())
}

func (w *Wallet) Import(str string) bool {
  encBytes, err := base64.StdEncoding.DecodeString(str)
  if err != nil {
    glog.Errorf("ERR: Import: %v", err)
    return false
  }
  return w._import(encBytes)
}

/*
func (w *Wallet) PublicKeyHash(coin string) ([]byte, bool) {
  hashedPublicKey := sha256.Sum256(w._get_public_key(coin))

  hasher := ripemd160.New()
  _, err := hasher.Write(hashedPublicKey[:])
  if err != nil {
    glog.Errorf("ERR: PublicKeyHash: %v", err)
    return nil, true
  }
  publicRipeMd := hasher.Sum(nil)

  return publicRipeMd, true
}*/

func (w *Wallet) GetMD5Hash(text string) string {
  hasher := md5.New()
  hasher.Write([]byte(text))
  return hex.EncodeToString(hasher.Sum(nil))
}

func (w *Wallet) Save(pathname string, password string) bool {
  cf := cipher.NewCFile()
  filename := pathname + "/" + w.GetMD5Hash(w.Name + w.Type) + ".wallet"
  return cf.SaveFilePwd(filename, password, w._export())
}

func (w *Wallet) Load(filename string, password string) bool {
  cf := cipher.NewCFile()
  buf, ok := cf.LoadFilePwd(filename, password)
  if !ok {
    return ok
  }
  return w._import(buf)
}

func (w *Wallet) GetBalance() ([]*messages.Balance) {
  b := messages.NewBalances()
  return b.GetBalanses()
}

func (w *Wallet) GetAddress(coin string) string {
  return ""
}
