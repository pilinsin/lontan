package gui

import (
	"context"
	"strings"
	"encoding/base64"
	"golang.org/x/crypto/argon2"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	i2p "github.com/pilinsin/go-libp2p-i2p"
	pv "github.com/pilinsin/p2p-verse"
	store "github.com/pilinsin/lontan/store"
)

func pageToTabItem(title string, page fyne.CanvasObject) *container.TabItem {
	return container.NewTabItem(title, page)
}

func storeHash(title, stAddr string) string{
	b := argon2.IDKey([]byte(stAddr), []byte(title), 1, 64*1024, 4, 64)
	return base64.URLEncoding.EncodeToString(b)
}

type GUI struct {
	rt *i2p.I2pRouter
	stores map[string]store.IDocumentStore
	bs map[string]pv.IBootstrap
	w    fyne.Window
	size fyne.Size
	tabs *container.AppTabs
	page *fyne.Container
}

func New(title string, width, height float32) *GUI {
	rt := i2p.NewI2pRouter()
	stores := make(map[string]store.IDocumentStore)
	bs := make(map[string]pv.IBootstrap)
	size := fyne.NewSize(width, height)
	a := app.New()
	a.Settings().SetTheme(theme.LightTheme())
	win := a.NewWindow(title)
	win.Resize(size)
	tabs := container.NewAppTabs()
	page := container.NewMax()
	return &GUI{rt, stores, bs, win, size, tabs, page}
}

func (gui *GUI) withRemove(page fyne.CanvasObject) fyne.CanvasObject {
	rmvBtn := widget.NewButtonWithIcon("", theme.ContentClearIcon(), func() {
		gui.tabs.Remove(gui.tabs.Selected())
	})
	return container.NewBorder(container.NewBorder(nil, nil, nil, rmvBtn), nil, nil, nil, page)
}
func (gui *GUI) addPageToTabs(title string, page fyne.CanvasObject){
	withRmvPage := gui.withRemove(page)
	withRmvTab := pageToTabItem(title, withRmvPage)
	gui.tabs.Append(withRmvTab)
	gui.tabs.Select(withRmvTab)
	gui.page.Refresh()
}

func (gui *GUI) loadSearchPage(addr, uiStr string) (string, fyne.CanvasObject){
	addrs := strings.Split(strings.TrimPrefix(addr, "/"), "/")
	if len(addrs) != 3{
		return "", nil
	}
	title := addrs[0]
	stAddr := strings.Join(addrs[1:], "/")

	ui := &store.UserIdentity{}
	if err := ui.FromString(uiStr); err != nil{ui = nil}

	storesKey := storeHash(title, addrs[2])
	st, ok := gui.stores[storesKey]
	if !ok{
		var err error
		st, err = store.LoadDocumentStore(context.Background(), stAddr, ui)
		if err != nil{return "", nil}
		gui.stores[storesKey] = st
	}

	return gui.NewSearchPage(gui.w, title, st, ui)
}

func (gui *GUI) loadPageForm() fyne.CanvasObject {
	addrEntry := widget.NewEntry()
	addrEntry.PlaceHolder = "Store Address"
	uiEntry := widget.NewEntry()
	uiEntry.PlaceHolder = "User Identity Address"

	onTapped := func() {
		title, loadPage := gui.loadSearchPage(addrEntry.Text, uiEntry.Text)
		addrEntry.SetText("")
		uiEntry.SetText("")
		if loadPage == nil {
			return
		}
		//page.SetMinSize(fyne,NewSize(101.1,201.2))
		gui.addPageToTabs(title, loadPage)
	}
	loadBtn := widget.NewButtonWithIcon("", theme.NavigateNextIcon(), onTapped)

	entries := container.NewVBox(addrEntry, uiEntry)
	return container.NewBorder(nil, nil, nil, loadBtn, entries)
}
func (gui *GUI) defaultPage(note *widget.Label) *container.TabItem {
	newForm := gui.loadPageForm()
	setup := gui.NewSetupPage()
	return pageToTabItem("top page", container.NewBorder(newForm,note,nil,nil, setup))
}

func (gui *GUI) initErrorPage() {
	for _, obj := range gui.page.Objects {
		gui.page.Remove(obj)
	}
	failed := widget.NewLabel("i2p router failed to start. please try again later.")
	gui.page.Add(container.NewCenter(failed))
	gui.page.Refresh()
}
func (gui *GUI) i2pStart(i2pNote *widget.Label) {
	go func() {
		if err := gui.rt.Start(); err == nil {
			i2pNote.SetText("i2p router on")
		} else {
			gui.initErrorPage()
		}
	}()
}

func (gui *GUI) Close(){
	for _, st := range gui.stores{
		st.Close()
		st = nil
	}
	for _, b := range gui.bs{
		b.Close()
		b = nil
	}
	gui.rt.Stop()
}

func (gui *GUI) Run() {
	i2pNote := widget.NewLabel("i2p router setup...")
	gui.i2pStart(i2pNote)

	gui.tabs.Append(gui.defaultPage(i2pNote))
	gui.page.Add(gui.tabs)
	gui.w.SetContent(gui.page)
	gui.w.SetOnClosed(gui.Close)
	gui.w.ShowAndRun()
}
