package store

import (
	"time"

	proto "google.golang.org/protobuf/proto"

	pb "github.com/pilinsin/lontan/store/pb"
)

func mapToSlice(mp map[string]struct{}) []string {
	slice := make([]string, len(mp))
	idx := 0
	for k := range mp {
		slice[idx] = k
		idx++
	}
	return slice
}

type DocumentInfo struct {
	Title       string
	Time        time.Time
	DocTypes    []string
	Tags        []string
	Description string
}

func NewDocumentInfo(title, description string, docTypesMap, tagsMap map[string]struct{}, t time.Time) *DocumentInfo {
	docTypes := mapToSlice(docTypesMap)
	tags := mapToSlice(tagsMap)
	return &DocumentInfo{title, t, docTypes, tags, description}
}

type typedCid struct {
	Type string
	Cid  string
}

func (tc *typedCid) encode() *pb.TypedCid {
	return &pb.TypedCid{
		Type: tc.Type,
		Cid:  tc.Cid,
	}
}
func (tc *typedCid) decode(pbtc *pb.TypedCid) {
	tc.Type = pbtc.GetType()
	tc.Cid = pbtc.GetCid()
}
func encodeTypedCids(tcs []typedCid) []*pb.TypedCid {
	pbtcs := make([]*pb.TypedCid, len(tcs))
	for idx, tc := range tcs {
		pbtcs[idx] = tc.encode()
	}
	return pbtcs
}
func decodeTypedCids(pbtcs []*pb.TypedCid) []typedCid {
	tcs := make([]typedCid, len(pbtcs))
	for idx, pbtc := range pbtcs {
		tcs[idx].decode(pbtc)
	}
	return tcs
}

type Document struct {
	*DocumentInfo
	Cids []typedCid
}

func newEmptyDocument() *Document {
	cids := make([]typedCid, 0)
	return &Document{
		DocumentInfo: &DocumentInfo{},
		Cids:         cids,
	}
}
func newDocument(di *DocumentInfo, cids ...typedCid) *Document {
	return &Document{di, cids}
}
func (d *Document) Marshal() []byte {
	mt, _ := d.Time.MarshalBinary()
	mui := &pb.Document{
		Cids:   encodeTypedCids(d.Cids),
		Title:  d.Title,
		Time:   mt,
		Types:  d.DocTypes,
		Tags:   d.Tags,
		Dscrpt: d.Description,
	}
	m, _ := proto.Marshal(mui)
	return m
}
func (d *Document) Unmarshal(m []byte) error {
	md := &pb.Document{}
	if err := proto.Unmarshal(m, md); err != nil {
		return err
	}
	t := time.Time{}
	if err := t.UnmarshalBinary(md.GetTime()); err != nil {
		return err
	}

	d.Cids = decodeTypedCids(md.GetCids())
	d.Title = md.GetTitle()
	d.Time = t
	d.DocTypes = md.GetTypes()
	d.Tags = md.GetTags()
	d.Description = md.GetDscrpt()
	return nil
}

type NamedDocument struct {
	*Document
	Name string
}
