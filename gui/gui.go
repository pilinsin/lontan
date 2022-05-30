package gui

import (
	"context"
	"encoding/base64"
	"golang.org/x/crypto/argon2"
	"path/filepath"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	i2p "github.com/pilinsin/go-libp2p-i2p"
	store "github.com/pilinsin/lontan/store"
	pv "github.com/pilinsin/p2p-verse"
)

func pageToTabItem(title string, page fyne.CanvasObject) *container.TabItem {
	return container.NewTabItem(title, page)
}

func storeHash(title, stAddr string) string {
	b := argon2.IDKey([]byte(stAddr), []byte(title), 1, 64*1024, 4, 64)
	return base64.URLEncoding.EncodeToString(b)
}

type GUI struct {
	rt     *i2p.I2pRouter
	stores map[string]store.IDocumentStore
	bs     map[string]pv.IBootstrap
	w      fyne.Window
	size   fyne.Size
	tabs   *container.AppTabs
	page   *fyne.Container
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
func (gui *GUI) addPageToTabs(title string, page fyne.CanvasObject) {
	withRmvPage := gui.withRemove(page)
	withRmvTab := pageToTabItem(title, withRmvPage)
	gui.tabs.Append(withRmvTab)
	gui.tabs.Select(withRmvTab)
	gui.page.Refresh()
}

func (gui *GUI) loadSearchPage(bAddr, stAddr string) (string, fyne.CanvasObject) {
	addrs := strings.Split(stAddr, "/")
	if len(addrs) != 2 {
		return "", nil
	}
	title, rawStAddr := addrs[0], addrs[1]

	storesKey := storeHash(title, rawStAddr)
	st, ok := gui.stores[storesKey]
	if !ok {
		baseDir := store.BaseDir(filepath.Join("stores", storesKey))
		var err error
		st, err = store.LoadDocumentStore(context.Background(), bAddr+"/"+stAddr, baseDir)
		if err != nil {
			return "", nil
		}
		gui.stores[storesKey] = st
	}

	return title, gui.NewSearchPage(gui.w, title, st)
}

func (gui *GUI) loadPageForm() fyne.CanvasObject {
	bAddrEntry := widget.NewEntry()
	bAddrEntry.SetPlaceHolder("Bootstraps Address")
	stAddrEntry := widget.NewEntry()
	stAddrEntry.SetPlaceHolder("Store Address")
	addrEntry := container.NewGridWithColumns(2, bAddrEntry, stAddrEntry)

	onTapped := func() {
		title, loadPage := gui.loadSearchPage(bAddrEntry.Text, stAddrEntry.Text)
		bAddrEntry.SetText("")
		stAddrEntry.SetText("")
		if loadPage == nil {
			return
		}
		//page.SetMinSize(fyne,NewSize(101.1,201.2))
		gui.addPageToTabs(title, loadPage)
	}
	loadBtn := widget.NewButtonWithIcon("", theme.NavigateNextIcon(), onTapped)

	return container.NewBorder(nil, nil, nil, loadBtn, addrEntry)
}
func (gui *GUI) defaultPage(note *widget.Label) *container.TabItem {
	newForm := gui.loadPageForm()
	setup := gui.NewSetupPage()
	return pageToTabItem("top page", container.NewBorder(newForm, note, nil, nil, setup))
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

func (gui *GUI) Close() {
	for _, st := range gui.stores {
		st.Close()
		st = nil
	}
	for _, b := range gui.bs {
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
