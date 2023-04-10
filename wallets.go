package wallets

import (
  "path/filepath"
  "github.com/golang/glog"
)

type Wallets struct {
  Wallets    []IWallet             `yaml:"wallets,omitempty"`
  Types      *TypesWallet          `yaml:"-"`
}

func NewWallets() *Wallets {
  return &Wallets{
       Wallets: make([]IWallet, 0),
       Types: NewTypesWallet(),
    }
}

func (ws *Wallets) Count() int {
  return len(ws.Wallets)
}

func (ws *Wallets) Get(i int) IWallet {
  return ws.Wallets[i] 
}

func (ws *Wallets) Remove(w IWallet) {
  
  //ws.Wallets = append(ws.Wallets[:i], ws.Wallets[i+1:]...)
}

func (ws *Wallets) Add(w IWallet) {
  ws.Wallets = append(ws.Wallets, w)
}

func (ws *Wallets) Load(scanPath string, password string) bool {
  files, err := filepath.Glob(scanPath)
  if err != nil {
    glog.Errorf("ERR: scanPath(%s)  #%v ", scanPath, err)
    return false
  }
  for _, filename := range files {
    if glog.V(2) {
      glog.Infof("LOG: Loading file: '%s'", filename)
    }
    w := newWallet()
    if !w.Load(filename, password) {
      continue
    }
    nw := NewWallet(w.GetType())
    if !nw.Load(filename, password) {
      continue
    }
    ws.Add(nw)
  }

  return true
}

