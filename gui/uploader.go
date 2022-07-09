package gui

import (
	"bytes"
	"fmt"
	"io"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	gutil "github.com/pilinsin/lontan/gui/util"
	store "github.com/pilinsin/lontan/store"
	ipfs "github.com/pilinsin/p2p-verse/ipfs"
)

var exts = map[string]struct{}{
	"text":  {},
	"image": {},
	"pdf":   {},
	"video": {},
	"audio": {},
}

func sliceToMap(slc []string) map[string]struct{} {
	mp := make(map[string]struct{}, len(slc))
	for _, elem := range slc {
		mp[elem] = struct{}{}
	}
	return mp
}

type iText interface {
	SetText(string)
}

func uplpadDialog(w fyne.Window, lable iText, is ipfs.Ipfs, ext string, ub *uploadBtn) func() {
	return func() {
		onSelected := func(rc fyne.URIReadCloser, err error) {
			if rc == nil || err != nil {
				lable.SetText("no file is selected")
				return
			}
			if _, ok := exts[ext]; !ok {
				lable.SetText("invalid file is selected")
				return
			}

			var r io.Reader
			if ext == "pdf" {
				pdf, err := store.EncodePdf(rc)
				if err != nil {
					lable.SetText("invalid pdf is selected")
					return
				}
				r = pdf
			} else if ext == "video" {
				v, err := store.EncodeVideo(rc, is)
				if err != nil {
					lable.SetText("invalid video is selected")
					return
				}
				r = v
			} else if ext == "audio" {
				a, err := store.EncodeAudio(rc, is)
				if err != nil {
					lable.SetText("invalid audio is selected")
					return
				}
				r = a
			} else if ext == "image" {
				img, err := store.EncodeImage(rc)
				if err != nil {
					lable.SetText("invalid image is selected")
					return
				}
				r = img
			}

			ub.td = store.NewTypedData(ext, r)
			lable.SetText(ext + " added")
		}
		dialog.ShowFileOpen(onSelected, w)
	}
}

//dialog
func NewUploadPage(w fyne.Window, st store.IDocumentStore) fyne.CanvasObject {
	noteLabel := widget.NewLabel("upload file")

	ui := widget.NewEntry()
	ui.SetPlaceHolder("user identity")
	name := widget.NewEntry()
	name.SetPlaceHolder("document name: <pid/username/docname>")
	title := widget.NewEntry()
	title.SetPlaceHolder("title")
	description := widget.NewMultiLineEntry()
	description.SetPlaceHolder("description")
	tags := gutil.NewRemovableEntryForm()

	dataObjs := container.NewVBox()

	txtBtn := newTextUploadButton(dataObjs)
	imgBtn := newDataUploadButton(w, dataObjs, "image", st.Ipfs())
	pdfBtn := newDataUploadButton(w, dataObjs, "pdf", st.Ipfs())
	//vdBtn := newDataUploadButton(w, dataObjs, "video", st.Ipfs())
	adBtn := newDataUploadButton(w, dataObjs, "audio", st.Ipfs())
	btns := container.NewHBox(txtBtn, imgBtn, pdfBtn, adBtn) //, vdBtn)

	uploadBtn := widget.NewButtonWithIcon("", theme.UploadIcon(), func() {
		noteLabel.SetText("processing...")

		tds := make([]*store.TypedData, 0)
		docTypes := make([]string, 0)
		for _, obj := range dataObjs.Objects {
			tdExtractor, ok := extractorFromRemoveBtn(obj)
			if !ok {
				continue
			}

			td := tdExtractor.TypedData()
			if td == nil || td.Data() == nil {
				continue
			}
			tds = append(tds, td)
			docTypes = append(docTypes, td.Type())
		}
		if len(docTypes) == 0 {
			noteLabel.SetText("no valid data")
			return
		}

		uid := &store.UserIdentity{}
		if err := uid.FromString(ui.Text); err != nil {
			uid = nil
		}
		st.SetUserIdentity(uid)

		docInfo := store.NewDocumentInfo(title.Text, description.Text, sliceToMap(docTypes), sliceToMap(tags.Texts()), time.Now().UTC())
		if err := st.Put(name.Text, docInfo, tds...); err != nil {
			noteLabel.SetText(fmt.Sprintln("upload error", err))
		} else {
			noteLabel.SetText("uploaded")
		}
	})

	upBtnLabel := container.NewBorder(nil, nil, uploadBtn, nil, noteLabel)
	page := container.NewVBox(ui, name, title, description, tags.Render(), btns, dataObjs, upBtnLabel)
	return container.NewMax(container.NewVScroll(page))
}

func withRemoveBtn(objs *fyne.Container, obj fyne.CanvasObject) fyne.CanvasObject {
	rmBtn := &widget.Button{
		Text: "",
		Icon: theme.ContentClearIcon(),
	}
	withRmObj := container.NewBorder(container.NewBorder(nil, nil, nil, rmBtn), nil, nil, nil, obj)
	rmBtn.OnTapped = func() { objs.Remove(withRmObj) }
	rmBtn.ExtendBaseWidget(rmBtn)

	return withRmObj
}
func extractorFromRemoveBtn(rmbObj fyne.CanvasObject) (iTypedDataExtractor, bool) {
	ct, ok := rmbObj.(*fyne.Container)
	if !ok {
		return nil, false
	}

	obj := ct.Objects[0]
	tdExtractor, ok := obj.(iTypedDataExtractor)
	return tdExtractor, ok
}

func extToIcon(ext string) fyne.Resource {
	switch ext {
	case "text":
		return theme.DocumentCreateIcon()
	case "image":
		return theme.MediaPhotoIcon()
	case "pdf":
		return theme.DocumentIcon()
	case "video":
		return theme.MediaVideoIcon()
	case "audio":
		return theme.MediaMusicIcon()
	default:
		return nil
	}
}

type iTypedDataExtractor interface {
	fyne.CanvasObject
	TypedData() *store.TypedData
}

type multiEntry struct {
	*widget.Entry
	/*
		isRT bool
		e *widget.Entry
		rt *widget.RichText
	*/
}

func NewMultiEntry() iTypedDataExtractor {
	me := &multiEntry{
		Entry: &widget.Entry{
			MultiLine: true,
			Wrapping:  fyne.TextTruncate,
		},
	}
	me.ExtendBaseWidget(me)
	me.SetPlaceHolder("input markdown text")
	return me
	/*
		rt := widget.NewRichTextFromMarkdown("")
		rt.Wrapping = fyne.TextWrapWord
		rt.Hide()

		return &multiEntry{
			e: e,
			rt: rt,
			isRT: false,
		}
	*/
}

/*
func (me *multiEntry) Render() fyne.CanvasObject{
	eBtn := widget.NewButtonWithIcon("", theme.DocumentCreateIcon(), func(){
		if me.isRT{
			me.isRT = false
			me.rt.Hide()
			me.e.Show()
		}
	})
	rtBtn := widget.NewButtonWithIcon("", theme.VisibilityIcon(), func(){
		if !me.isRT{
			me.isRT = true
			me.rt.ParseMarkdown(me.e.Text)
			me.e.Hide()
			me.rt.Show()
		}
	})

	btns := container.NewHBox(eBtn, rtBtn)
	ert := container.NewVBox(me.e, me.rt)
	return container.NewBorder(container.NewBorder(nil,nil,btns,nil),nil,nil,nil, ert)
}
*/
func (me *multiEntry) TypedData() *store.TypedData {
	if me.Text == "" {
		return nil
	} else {
		return store.NewTypedData("text", bytes.NewBufferString(me.Text))
	}
}

func newTextUploadButton(objs *fyne.Container) fyne.CanvasObject {
	return widget.NewButtonWithIcon("", extToIcon("text"), func() {
		objs.Add(withRemoveBtn(objs, NewMultiEntry()))
	})
}

type uploadBtn struct {
	*widget.Button
	td *store.TypedData
}

func NewUploadButton(w fyne.Window, ext string, is ipfs.Ipfs) iTypedDataExtractor {
	ub := &uploadBtn{}

	btn := &widget.Button{
		Text: "add " + ext + " file",
		Icon: theme.UploadIcon(),
	}
	btn.OnTapped = uplpadDialog(w, btn, is, ext, ub)

	ub.Button = btn
	ub.ExtendBaseWidget(ub)
	return ub
}
func (ub *uploadBtn) TypedData() *store.TypedData {
	return ub.td
}

func newDataUploadButton(w fyne.Window, objs *fyne.Container, ext string, is ipfs.Ipfs) fyne.CanvasObject {
	return widget.NewButtonWithIcon("", extToIcon(ext), func() {
		objs.Add(withRemoveBtn(objs, NewUploadButton(w, ext, is)))
	})
}
