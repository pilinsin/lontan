package gui

import (
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	query "github.com/ipfs/go-datastore/query"

	store "github.com/pilinsin/lontan/store"
	crdt "github.com/pilinsin/p2p-verse/crdt"
)

var mode = []string{
	"key (pid/username/docname)",
	"title",
	"cid",
	"document type",
	"tag",
}
var order = []string{
	"newer",
	"older",
}

func (gui *GUI) NewSearchPage(w fyne.Window, title string, st store.IDocumentStore) fyne.CanvasObject {
	uploadBtn := widget.NewButtonWithIcon("", theme.UploadIcon(), func() {
		gui.addPageToTabs(title+"_upload", NewUploadPage(w, st))
	})

	modeSelector := widget.NewSelect(mode, nil)
	searchEntry := widget.NewEntry()
	searchEntry.SetPlaceHolder("search text")
	orderBtn := widget.NewSelect(order, nil)
	orderBtn.Selected = order[0]

	docs := container.NewVBox()
	var ndocs <-chan *store.NamedDocument
	newViewPageButton := func(ndoc *store.NamedDocument, st store.IDocumentStore) fyne.CanvasObject {
		hline := widget.NewRichTextFromMarkdown("-----")
		btn := widget.NewButtonWithIcon("", theme.NavigateNextIcon(), func() {
			vpage, closer := NewViewerPage(gui, ndoc, st)
			gui.addPageToTabs(title+"_view_"+ndoc.Title, vpage, closer)
		})
		return container.NewBorder(hline, nil, btn, nil, newDocumentCard(ndoc))
	}
	resetDocs := func() {
		for _, obj := range docs.Objects {
			docs.Remove(obj)
		}
	}
	loadDocs := func() {
		N := 10
		for i := 0; i < N; i++ {
			ndoc, ok := <-ndocs
			if ok {
				docs.Add(newViewPageButton(ndoc, st))
			}
		}
	}
	searchBtn := widget.NewButtonWithIcon("", theme.SearchIcon(), func() {
		es := strings.Fields(searchEntry.Text)
		qf := modeToQueryFunc(modeSelector.Selected)
		q := qf(es...)
		q.Orders = []query.Order{store.TimeOrder{FrontNew: orderBtn.Selected == order[0]}}
		var err error
		ndocs, err = st.Query(q)
		if err != nil {
			searchEntry.SetText("")
			return
		}
		resetDocs()
		loadDocs()
	})
	moreBtn := widget.NewButtonWithIcon("", theme.MoveDownIcon(), loadDocs)

	orderSearch := container.NewHBox(orderBtn, searchBtn)
	searchObj := container.NewBorder(nil, nil, modeSelector, orderSearch, searchEntry)
	upObj := container.NewBorder(nil, nil, uploadBtn, nil)

	searchBar := container.NewBorder(upObj, nil, nil, nil, searchObj)
	moreObj := container.NewCenter(moreBtn)
	docsObj := container.NewMax(container.NewVScroll(docs))

	return container.NewBorder(searchBar, moreObj, nil, nil, docsObj)
}

type queryFunc func(strs ...string) query.Query

func modeToQueryFunc(mode string) queryFunc {
	return func(strs ...string) query.Query {
		switch mode {
		case "key (pid/username/docname)":
			fs := make([]query.Filter, len(strs))
			for idx, str := range strs {
				if strings.Contains(str, "/") {
					fs[idx] = crdt.KeyMatchFilter{Key: str}
				} else {
					fs[idx] = crdt.KeyExistFilter{Key: str}
				}
			}
			return query.Query{Filters: fs}
		case "title":
			fs := make([]query.Filter, len(strs))
			for idx, str := range strs {
				fs[idx] = store.TitleFilter{Title: str}
			}
			return query.Query{Filters: fs}
		case "cid":
			return query.Query{Filters: []query.Filter{store.CidsFilter{Cids: strs}}}
		case "document type":
			return query.Query{Filters: []query.Filter{store.DocTypesFilter{DocTypes: strs}}}
		case "tag":
			return query.Query{Filters: []query.Filter{store.TagsFilter{Tags: strs}}}
		default:
			return query.Query{}
		}
	}
}

func newDocumentCard(ndoc *store.NamedDocument) fyne.CanvasObject {
	nm := &widget.Label{
		Text:     ndoc.Name,
		Wrapping: fyne.TextTruncate,
	}
	nm.ExtendBaseWidget(nm)

	ttl := &widget.Label{
		Text:     ndoc.Title,
		Wrapping: fyne.TextTruncate,
	}
	ttl.ExtendBaseWidget(ttl)

	tm := &widget.Label{
		Text:     ndoc.Time.String(),
		Wrapping: fyne.TextTruncate,
	}
	tm.ExtendBaseWidget(tm)

	nDesc := 200
	desc := &widget.Label{
		Text:     extractDescription(ndoc.Description, nDesc),
		Wrapping: fyne.TextTruncate,
	}
	desc.ExtendBaseWidget(desc)

	tps := docTypesToIcons(ndoc.DocTypes)

	return container.NewVBox(ttl, desc, tps, tm, nm)
}
func extractDescription(desc string, n int) string {
	if len(desc) <= n {
		return desc
	} else {
		return desc[:n]
	}
}
