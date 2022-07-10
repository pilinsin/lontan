package gui

import (
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	gutil "github.com/pilinsin/lontan/gui/util"
	store "github.com/pilinsin/lontan/store"
	ipfs "github.com/pilinsin/p2p-verse/ipfs"
)

func loadMedia(gui *GUI, tp, cid string, is ipfs.Ipfs) (fyne.CanvasObject, gutil.Closer) {
	switch tp {
	case "text":
		return LoadText(cid, is)
	case "image":
		return LoadImage(gui, cid, is)
	case "pdf":
		return LoadPdf(gui, cid, is)
	case "video":
		return LoadVideo(cid, is)
	case "audio":
		return LoadAudio(cid, is)
	default:
		return errorLabel("invalid cid"), nil
	}
}

func docTypesToIcons(docTypes []string) fyne.CanvasObject {
	icons := make([]fyne.CanvasObject, len(docTypes))
	for idx, ext := range docTypes {
		icons[idx] = widget.NewButtonWithIcon("", extToIcon(ext), nil)
	}
	return container.NewHBox(icons...)
}

func tagsLabel(tags []string) fyne.CanvasObject {
	lbl := &widget.Label{
		Text:     strings.Join(tags, ", "),
		Wrapping: fyne.TextWrapWord,
	}
	lbl.ExtendBaseWidget(lbl)
	return lbl
}
func descriptionLabel(text string) fyne.CanvasObject {
	lbl := &widget.Label{
		Text:     text,
		Wrapping: fyne.TextWrapWord,
	}
	lbl.ExtendBaseWidget(lbl)
	return lbl
}

func NewViewerPage(gui *GUI, nmDoc *store.NamedDocument, st store.IDocumentStore) (fyne.CanvasObject, gutil.Closer) {
	if nmDoc == nil {
		return container.NewCenter(widget.NewLabel("no document")), nil
	}

	medias := make([]fyne.CanvasObject, len(nmDoc.Cids))
	closers := make([]gutil.Closer, 0)
	for idx, cid := range nmDoc.Cids {
		media, closer := loadMedia(gui, cid.Type, cid.Cid, st.Ipfs())
		medias[idx] = media
		if closer != nil {
			closers = append(closers, closer)
		}
	}
	closer := func() error {
		var err error
		for _, closer := range closers {
			if closeErr := closer(); closeErr != nil {
				err = closeErr
			}
		}
		return err
	}

	name := descriptionLabel(nmDoc.Name)
	title := descriptionLabel(nmDoc.Title)
	tm := descriptionLabel(nmDoc.Time.String())
	dTypes := docTypesToIcons(nmDoc.DocTypes)
	tags := tagsLabel(nmDoc.Tags)
	description := descriptionLabel(nmDoc.Description)

	objs := make([]fyne.CanvasObject, 0)
	objs = append(objs, title, name, tm, dTypes, tags, description)
	objs = append(objs, medias...)
	page := container.NewVBox(objs...)
	return container.NewMax(container.NewVScroll(page)), closer
}
