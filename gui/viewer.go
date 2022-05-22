package gui

import(
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"fyne.io/fyne/v2/container"

	ipfs "github.com/pilinsin/p2p-verse/ipfs"
	store "github.com/pilinsin/lontan/store"
)

func loadMedia(tp, cid string, is ipfs.Ipfs) fyne.CanvasObject{
	switch tp {
	case "text":
		return LoadText(cid, is)
	case "image":
		return LoadImage(cid, is)
	case "pdf":
		return LoadPdf(cid, is)
	case "video":
		return LoadVideo(cid, is)
	case "audio":
		return LoadAudio(cid, is)
	default:
		return errorLabel("invalid cid")
	}
}

func docTypesToIcons(docTypes []string) fyne.CanvasObject{
	icons := make([]fyne.CanvasObject, len(docTypes))
	for idx, ext := range docTypes{
		icons[idx] = widget.NewButtonWithIcon("", extToIcon(ext), nil)
	}
	return container.NewHBox(icons...)
}


func tagsLabel(tags []string) fyne.CanvasObject{
	lbl := &widget.Label{
		Text: strings.Join(tags, ", "),
		Wrapping: fyne.TextWrapWord,
	}
	lbl.ExtendBaseWidget(lbl)
	return lbl
}
func descriptionLabel(text string) fyne.CanvasObject{
	lbl := &widget.Label{
		Text: text,
		Wrapping: fyne.TextWrapWord,
	}
	lbl.ExtendBaseWidget(lbl)
	return lbl
}

func NewViewerPage(nmDoc *store.NamedDocument, st store.IDocumentStore) fyne.CanvasObject{
	if nmDoc == nil{
		return container.NewCenter(widget.NewLabel("no document"))
	}
	
	medias := make([]fyne.CanvasObject, len(nmDoc.Cids))
	for idx, cid := range nmDoc.Cids{
		medias[idx] = loadMedia(cid.Type, cid.Cid, st.Ipfs())
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
	return container.NewMax(container.NewVScroll(page))
}