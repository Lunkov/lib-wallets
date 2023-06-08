package wallets

import (
  "os"
  "sync"
  "path/filepath"
  "github.com/Lunkov/go-hdwallet"
  
  "github.com/golang/glog"
)

type WalletStorage struct {
  Wallet   IWallet
  Filename string
}

type Wallets struct {
  Wallets    []WalletStorage       `yaml:"wallets,omitempty"`
  mu           sync.RWMutex        `yaml:"-"`
  Types      *TypesWallet          `yaml:"-"`
}

func NewWallets() *Wallets {
  return &Wallets{
       Wallets: make([]WalletStorage, 0),
       Types: NewTypesWallet(),
    }
}

func (ws *Wallets) Count() int {
  ws.mu.RLock()
  sz := len(ws.Wallets)
  ws.mu.RUnlock()
  return sz
}

func (ws *Wallets) Get(i int) IWallet {
  ws.mu.RLock()
  defer ws.mu.RUnlock()
  defer func() {
                if r := recover(); r != nil {
                }
        }()
  wallet := ws.Wallets[i] 
  return wallet.Wallet
}

func (ws *Wallets) Remove(w IWallet) {
  ws.mu.Lock()
  defer ws.mu.Unlock()
  addr := w.GetAddress(hdwallet.ECOS)
  for i := 0; i < len(ws.Wallets); i ++ {
    if ws.Wallets[i].Wallet.GetAddress(hdwallet.ECOS) == addr {
      _, err := os.Stat(ws.Wallets[i].Filename)
      if err == nil {
        os.Remove(ws.Wallets[i].Filename)
      }
      ws.Wallets = append(ws.Wallets[:i], ws.Wallets[i+1:]...)
      break
    }
  }
}

func (ws *Wallets) Add(w IWallet) {
  wallet := WalletStorage{Wallet: w}
  ws.mu.Lock()
  ws.Wallets = append(ws.Wallets, wallet)
  ws.mu.Unlock()
}

func (ws *Wallets) GetList() []string {
  ws.mu.RLock()
  defer ws.mu.RUnlock()
  res := make([]string, 0)
  for _, w := range ws.Wallets {
    res = append(res, w.Wallet.GetName() + " (" + w.Wallet.GetAddress(hdwallet.ECOS) + ")")
  } 
  return res
}

func (ws *Wallets) GetWalletByName(name string) (IWallet, bool) {
  ws.mu.RLock()
  defer ws.mu.RUnlock()
  for _, w := range ws.Wallets {
    if name == w.Wallet.GetName() + " (" + w.Wallet.GetAddress(hdwallet.ECOS) + ")" {
      return w.Wallet, true
    }
  } 
  return nil, false
}

func (ws *Wallets) Load(scanPath string, password string) bool {
  ws.mu.Lock()
  defer ws.mu.Unlock()
  files, err := filepath.Glob(scanPath)
  if err != nil {
    if glog.V(2) {
      glog.Errorf("ERR: scanPath(%s)  #%v ", scanPath, err)
    }
    return false
  }
  for _, filename := range files {
    if glog.V(2) {
      glog.Infof("LOG: Loading file: '%s'", filename)
    }
    w := NewEmptyWallet()
    if !w.Load(filename, password) {
      continue
    }
    nw := NewWallet(w.Type)
    if !nw.Load(filename, password) {
      continue
    }
    wallet := WalletStorage{Wallet: nw, Filename: filename}
    ws.Wallets = append(ws.Wallets, wallet)
  }
  return true
}

