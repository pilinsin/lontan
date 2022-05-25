package store

import (
	"context"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	query "github.com/ipfs/go-datastore/query"

	i2p "github.com/pilinsin/go-libp2p-i2p"
	pv "github.com/pilinsin/p2p-verse"
	crdt "github.com/pilinsin/p2p-verse/crdt"
	ipfs "github.com/pilinsin/p2p-verse/ipfs"
)

type TypedData struct {
	tp   string
	data io.Reader
}

func NewTypedData(tp string, data io.Reader) *TypedData {
	return &TypedData{tp, data}
}
func (td *TypedData) Type() string    { return td.tp }
func (td *TypedData) Data() io.Reader { return td.data }

type IDocumentStore interface {
	Close()
	Ipfs() ipfs.Ipfs
	SetUserIdentity(*UserIdentity)
	Address() string
	Put(string, *DocumentInfo, ...*TypedData) error
	Get(string) (*NamedDocument, error)
	Query(...query.Query) (<-chan *NamedDocument, error) //time, tag, etc...
}

type documentStore struct {
	ctx       context.Context
	closer    func()
	dirCloser func()
	addr      string
	userName  string
	is        ipfs.Ipfs
	ss        crdt.ISignatureStore
}

func NewDocumentStore(ctx context.Context, title, bAddr, baseDir string) (IDocumentStore, error) {
	bootstraps := pv.AddrInfosFromString(bAddr)
	save := false
	dirCloser := func() { os.Remove(baseDir) }

	ipfsDir := filepath.Join(baseDir, "ipfs")
	is, err := ipfs.NewIpfsStore(i2p.NewI2pHost, ipfsDir, "ipfs_kw", save, false, bootstraps...)
	if err != nil {
		return nil, err
	}

	storeDir := filepath.Join(baseDir, "store")
	v := crdt.NewVerse(i2p.NewI2pHost, storeDir, save, false, bootstraps...)
	st, err := v.NewStore(pv.RandString(8), "signature")
	if err != nil {
		is.Close()
		return nil, err
	}
	ss := st.(crdt.ISignatureStore)
	ctx, cancel := context.WithCancel(context.Background())
	autoSync(ctx, ss)

	addr := bAddr + "/" + title + "/" + ss.Address()
	return &documentStore{ctx, cancel, dirCloser, addr, "Anonymous", is, ss}, nil
}
func LoadDocumentStore(ctx context.Context, addr, baseDir string) (IDocumentStore, error) {
	ui := parseUserIdentity(nil)
	bAddr, sAddr, err := parseAddr(addr)
	if err != nil {
		return nil, err
	}
	bootstraps := pv.AddrInfosFromString(bAddr)
	save := true

	ipfsDir := filepath.Join(baseDir, "ipfs")
	is, err := ipfs.NewIpfsStore(i2p.NewI2pHost, ipfsDir, "ipfs_kw", save, false, bootstraps...)
	if err != nil {
		return nil, err
	}

	storeDir := filepath.Join(baseDir, "store")
	v := crdt.NewVerse(i2p.NewI2pHost, storeDir, save, false, bootstraps...)
	opt := &crdt.StoreOpts{Pub: ui.verfKey, Priv: ui.signKey}
	st, err := v.LoadStore(ctx, sAddr, "signature", opt)
	if err != nil {
		is.Close()
		return nil, err
	}
	ss := st.(crdt.ISignatureStore)
	ctx, cancel := context.WithCancel(context.Background())
	autoSync(ctx, ss)

	return &documentStore{ctx, cancel, func() {}, addr, ui.userName, is, ss}, nil
}

func parseAddr(addr string) (string, string, error) {
	addrs := strings.Split(strings.TrimPrefix(addr, "/"), "/")
	if len(addrs) != 3 {
		return "", "", errors.New("invalid addr")
	}
	return addrs[0], addrs[2], nil
}

func parseUserIdentity(ui *UserIdentity) *UserIdentity {
	if ui == nil {
		return &UserIdentity{"Anonymous", nil, nil}
	} else {
		invalidName := ui.userName == ""
		invalidVerf := ui.verfKey == nil
		invalidSign := ui.signKey == nil
		if invalidName || invalidVerf || invalidSign {
			return &UserIdentity{"Anonymous", nil, nil}
		}
	}

	return ui
}

func autoSync(ctx context.Context, ss crdt.IStore) {
	ticker := time.NewTicker(time.Second * 30)
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				ss.Sync()
			}
		}
	}()
}

func (ds *documentStore) Close() {
	ds.closer()
	ds.is.Close()
	ds.ss.Close()

	time.Sleep(time.Second)
	ds.dirCloser()
}

func (ds *documentStore) Ipfs() ipfs.Ipfs { return ds.is }

func (ds *documentStore) SetUserIdentity(ui *UserIdentity) {
	ui = parseUserIdentity(ui)
	ds.userName = ui.userName
	ds.ss.ResetKeyPair(ui.signKey, ui.verfKey)
}
func (ds *documentStore) Address() string { return ds.addr }

func (ds *documentStore) Put(docName string, docInfo *DocumentInfo, data ...*TypedData) error {
	cids := make([]typedCid, 0)
	for _, td := range data {
		cid, err := ds.is.AddReader(td.data)
		if err == nil {
			cids = append(cids, typedCid{td.tp, cid})
		}
	}
	if len(cids) == 0 {
		return errors.New("no valid data")
	}

	doc := newDocument(docInfo, cids...)
	return ds.ss.Put(ds.userName+"/"+docName, doc.Marshal())
}

func (ds *documentStore) Get(key string) (*NamedDocument, error) {
	m, err := ds.ss.Get(key)
	if err != nil {
		return nil, err
	}
	doc := newEmptyDocument()
	if err := doc.Unmarshal(m); err != nil {
		return nil, err
	}

	return &NamedDocument{doc, key}, nil
}

func (ds *documentStore) Query(qs ...query.Query) (<-chan *NamedDocument, error) {
	var q query.Query
	if len(qs) > 0 {
		q = qs[0]
	} else {
		q = query.Query{
			Orders: []query.Order{TimeOrder{true}},
		}
	}
	q.Filters = append([]query.Filter{documentFilter{}}, q.Filters...)
	if q.KeysOnly {
		q.KeysOnly = false
	}

	rs, err := ds.ss.Query()
	if err != nil {
		return nil, err
	}
	rs = query.NaiveQueryApply(q, rs)

	ch := make(chan *NamedDocument, 10)
	go func() {
		defer close(ch)
		for res := range rs.Next() {
			doc := newEmptyDocument()
			if err := doc.Unmarshal(res.Value); err != nil {
				continue
			}
			ch <- &NamedDocument{doc, res.Key}
		}
	}()
	return ch, nil
}

type documentFilter struct{}

func (f documentFilter) Filter(e query.Entry) bool {
	// e.Key: pid/username/docname
	keys := strings.Split(strings.TrimPrefix(e.Key, "/"), "/")
	if len(keys) != 3 {
		return false
	}

	d := newEmptyDocument()
	err := d.Unmarshal(e.Value)
	return err == nil
}
