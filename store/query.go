package store

import(
	"time"
	"strings"
	query "github.com/ipfs/go-datastore/query"
)

func sliceToMap(slc []string) map[string]struct{}{
	mp := make(map[string]struct{}, len(slc))
	for _, elem := range slc{
		mp[elem] = struct{}{}
	}
	return mp
}

//a<b:-1, a==b:0, a>b:1
//x0 < x1 < x2 < ... < xn
type TimeOrder struct{
	FrontNew bool
}
func (o TimeOrder) Compare(a, b query.Entry) int {
	da := newEmptyDocument()
	if err := da.Unmarshal(a.Value); err != nil{return 1}
	db := newEmptyDocument()
	if err := db.Unmarshal(b.Value); err != nil{return -1}

	if da.Time.Equal(db.Time){
		return 0
	}
	if da.Time.After(db.Time){
		if o.FrontNew{
			return -1
		}else{
			return 1
		}
	}else{
		if o.FrontNew{
			return 1
		}else{
			return -1
		}
	}
}


type CidsFilter struct{
	Cids []string
}
func (f CidsFilter) Filter(e query.Entry) bool{
	d := newEmptyDocument()
	if err := d.Unmarshal(e.Value); err != nil{return false}
	
	docCids := make([]string, len(d.Cids))
	for idx, tc := range d.Cids{
		docCids[idx] = tc.Cid
	}
	cidMap := sliceToMap(docCids)
	for _, cid := range f.Cids{
		if _, ok := cidMap[cid]; !ok{return false}
	}
	return true
}

type TitleFilter struct{
	Title string
}
func (f TitleFilter) Filter(e query.Entry) bool{
	d := newEmptyDocument()
	if err := d.Unmarshal(e.Value); err != nil{return false}
	return strings.Contains(d.Title, f.Title)
}

type DocTypesFilter struct{
	DocTypes []string
}
func (f DocTypesFilter) Filter(e query.Entry) bool{
	d := newEmptyDocument()
	if err := d.Unmarshal(e.Value); err != nil{return false}

	docTypeMap := sliceToMap(d.DocTypes)
	for _, docType := range f.DocTypes{
		if _, ok := docTypeMap[docType]; !ok{return false}
	}
	return true
}

type TimeFilter struct{
	Begin, End time.Time
}
func (f TimeFilter) Filter(e query.Entry) bool{
	d := newEmptyDocument()
	if err := d.Unmarshal(e.Value); err != nil{return false}

	return (f.Begin.Before(d.Time) || f.Begin.Equal(d.Time)) && f.End.After(d.Time)
}


type TagsFilter struct{
	Tags []string
}
func (f TagsFilter) Filter(e query.Entry) bool{
	d := newEmptyDocument()
	if err := d.Unmarshal(e.Value); err != nil{return false}

	tagsMap := sliceToMap(d.Tags)
	for _, fTag := range f.Tags{
		if _, ok := tagsMap[fTag]; !ok{return false}
	}
	return true
}