package gui

import (
	"path/filepath"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	peer "github.com/libp2p/go-libp2p-core/peer"

	i2p "github.com/pilinsin/go-libp2p-i2p"
	gutil "github.com/pilinsin/lontan/gui/util"
	store "github.com/pilinsin/lontan/store"
	pv "github.com/pilinsin/p2p-verse"
)

func (gui *GUI) NewSetupPage() fyne.CanvasObject {
	form := newBootstrapsForm()
	baddrsLabel := gutil.NewCopyButton("bootstrap list address")
	bFunc := newBootstrap(gui, baddrsLabel, form)
	addrsBtn := widget.NewButtonWithIcon("", theme.NavigateNextIcon(), bFunc)

	titleEntry := widget.NewEntry()
	titleEntry.SetPlaceHolder("document store title")
	storeLabel := gutil.NewCopyButton("document store address")
	stFunc := newStore(gui, titleEntry, baddrsLabel, storeLabel)
	storeBtn := widget.NewButtonWithIcon("", theme.NavigateNextIcon(), stFunc)

	userNameEntry := widget.NewEntry()
	userNameEntry.SetPlaceHolder("user name")
	uiLabel := gutil.NewCopyButton("user identity")
	uiBtn := widget.NewButtonWithIcon("", theme.NavigateNextIcon(), func() {
		kp := store.NewKeyPair()
		ui := store.NewUserIdentity(userNameEntry.Text, kp.Verify(), kp.Sign())
		uiLabel.SetText(ui.ToString())
	})

	hline2 := widget.NewRichTextFromMarkdown("-----")
	uiStr := container.NewBorder(nil, nil, uiBtn, nil, uiLabel.Render())
	userObj := container.NewVBox(hline2, userNameEntry, uiStr)

	hline := widget.NewRichTextFromMarkdown("-----")
	baddrs := container.NewBorder(nil, nil, addrsBtn, nil, baddrsLabel.Render())
	staddr := container.NewBorder(nil, nil, storeBtn, nil, storeLabel.Render())
	manObj := container.NewVBox(hline, form.Render(), baddrs, titleEntry, staddr)

	return container.NewGridWithColumns(1, userObj, manObj)
}

func newBootstrap(gui *GUI, lbl *gutil.CopyButton, form *bootstrapsForm) func() {
	return func() {
		go func() {
			lbl.SetText("processing...")

			bsKey := "setup"
			if _, exist := gui.bs[bsKey]; exist {
				gui.bs[bsKey].Close()
				gui.bs[bsKey] = nil
			}

			b, err := pv.NewBootstrap(i2p.NewI2pHost, form.AddrInfos()...)
			if err != nil {
				lbl.SetText("bootstrap list address")
				return
			}

			baddrs := append(b.ConnectedPeers(), b.AddrInfo())
			s := pv.AddrInfosToString(baddrs...)
			if s == "" {
				lbl.SetText("bootstrap list address")
			} else {
				gui.bs[bsKey] = b
				lbl.SetText(s)
			}
		}()
	}
}

func newStore(gui *GUI, te *widget.Entry, bLabel, stLabel *gutil.CopyButton) func() {
	return func() {
		go func() {
			if bLabel.GetText() == "bootstrap list address" {
				return
			}

			stLabel.SetText("processing...")
			storesKey := "setup"
			if _, exist := gui.stores[storesKey]; exist {
				gui.stores[storesKey].Close()
				gui.stores[storesKey] = nil
			}

			baseDir := store.BaseDir(filepath.Join("stores", storesKey))
			st, err := store.NewDocumentStore(te.Text, bLabel.GetText(), baseDir)
			if err != nil {
				stLabel.SetText("document store address")
			} else {
				gui.stores[storesKey] = st
				addrs := strings.Split(st.Address(), "/")
				addr := strings.Join(addrs[1:], "/")
				stLabel.SetText(addr)
			}
		}()
	}
}

func addrInfoMapToSlice(m map[string]peer.AddrInfo) []peer.AddrInfo {
	ais := make([]peer.AddrInfo, len(m))
	idx := 0
	for _, v := range m {
		ais[idx] = v
		idx++
	}
	return ais
}

type bootstrapsForm struct {
	*gutil.RemovableEntryForm
}

func newBootstrapsForm() *bootstrapsForm {
	ref := gutil.NewRemovableEntryForm()
	return &bootstrapsForm{ref}
}
func (bf *bootstrapsForm) AddrInfos() []peer.AddrInfo {
	txts := bf.Texts()
	aiMap := make(map[string]peer.AddrInfo)

	for _, txt := range txts {
		ai := pv.AddrInfoFromString(txt)
		if ai.ID != "" && len(ai.Addrs) > 0 {
			aiMap[txt] = ai
		} else {
			ais := pv.AddrInfosFromString(txt)
			for _, ai := range ais {
				if ai.ID == "" || len(ai.Addrs) == 0 {
					continue
				}
				s := pv.AddrInfoToString(ai)
				aiMap[s] = ai
			}
		}
	}

	return addrInfoMapToSlice(aiMap)
}
