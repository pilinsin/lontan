package gui

import (
	"bytes"
	"errors"
	"io"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	gutil "github.com/pilinsin/lontan/gui/util"
	pb "github.com/pilinsin/lontan/store/pb"
	ipfs "github.com/pilinsin/p2p-verse/ipfs"
	proto "google.golang.org/protobuf/proto"
)

func errorLabel(text string) fyne.CanvasObject {
	label := &widget.Label{
		Text:     text,
		Wrapping: fyne.TextWrapWord,
	}
	label.ExtendBaseWidget(label)

	return container.NewCenter(label)
}

func loadImage(r io.Reader) (fyne.CanvasObject, error) {
	buf := &bytes.Buffer{}
	if _, err := buf.ReadFrom(r); err != nil {
		return nil, err
	}
	pbImage := &pb.Image{}
	if err := proto.Unmarshal(buf.Bytes(), pbImage); err != nil {
		return nil, err
	}

	name := pbImage.GetName()
	data := pbImage.GetData()
	if data == nil {
		return nil, errors.New("invalid image")
	}

	res := &fyne.StaticResource{
		StaticName:    name,
		StaticContent: data,
	}
	imgCanvas := canvas.NewImageFromResource(res)
	imgCanvas.FillMode = canvas.ImageFillContain

	return imgCanvas, nil
}

func withZoom(obj fyne.CanvasObject) fyne.CanvasObject {
	baseSize := obj.Size()

	var page *fyne.Container
	zoomInbtn := widget.NewButtonWithIcon("", theme.ZoomInIcon(), func() {
		if obj.Size().Height < baseSize.Height*2 {
			width := obj.Size().Width
			height := obj.Size().Height
			obj.Resize(fyne.NewSize(width+50, height+50))

			grid := container.NewGridWrap(obj.Size(), obj)
			page.Objects[0] = container.NewScroll(grid)
			page.Refresh()
		}
	})

	zoomOutbtn := widget.NewButtonWithIcon("", theme.ZoomOutIcon(), func() {
		if obj.Size().Height > baseSize.Height {
			width := obj.Size().Width
			height := obj.Size().Height
			obj.Resize(fyne.NewSize(width-50, height-50))

			grid := container.NewGridWrap(obj.Size(), obj)
			page.Objects[0] = container.NewScroll(grid)
			page.Refresh()
		}
	})
	zoomBtns := container.NewHBox(zoomInbtn, zoomOutbtn)

	page = container.NewBorder(container.NewBorder(nil, nil, zoomBtns, nil), nil, nil, nil, obj)
	return page
}

func LoadImage(gui *GUI, cid string, is ipfs.Ipfs) (fyne.CanvasObject, gutil.Closer) {
	r, err := is.GetReader(cid)
	if err != nil {
		return errorLabel("load image error (ipfs)"), nil
	}

	img, err := loadImage(r)
	if err != nil {
		return errorLabel("load image error"), nil
	}
	imgCanvas := container.NewGridWrap(fyne.NewSize(400, 400), img)
	zoomBtn := widget.NewButtonWithIcon("", theme.ViewFullScreenIcon(), func() {
		name := img.(*canvas.Image).Resource.Name()
		gui.addPageToTabs(name, withZoom(img))
	})
	return container.NewBorder(container.NewBorder(nil, nil, zoomBtn, nil), nil, nil, nil, imgCanvas), nil
}

func LoadText(cid string, is ipfs.Ipfs) (fyne.CanvasObject, gutil.Closer) {
	r, err := is.GetReader(cid)
	if err != nil {
		return errorLabel("load text error (ipfs)"), nil
	}
	buf := &bytes.Buffer{}
	_, err = buf.ReadFrom(r)
	if err != nil {
		return errorLabel("load text error"), nil
	}

	rt := widget.NewRichTextFromMarkdown(buf.String())
	rt.Wrapping = fyne.TextWrapWord
	return rt, nil
}

func LoadAudio(cid string, is ipfs.Ipfs) (fyne.CanvasObject, gutil.Closer) {
	ap, err := NewAudioPlayer(cid, is)
	if err != nil {
		return errorLabel("load audio error"), nil
	}
	return ap.Render()
}
func LoadVideo(cid string, is ipfs.Ipfs) (fyne.CanvasObject, gutil.Closer) {
	vp, err := NewVideoPlayer(cid, is)
	if err != nil {
		return errorLabel("load video error"), nil
	}
	return vp.Render()
}

func LoadPdf(gui *GUI, cid string, is ipfs.Ipfs) (fyne.CanvasObject, gutil.Closer) {
	m, err := is.Get(cid)
	if err != nil {
		return errorLabel("load pdf error (ipfs)"), nil
	}

	pbPdf := &pb.Pdf{}
	if err := proto.Unmarshal(m, pbPdf); err != nil {
		return errorLabel("load pdf error"), nil
	}

	mImgs := pbPdf.GetImages()
	imgCanvases := make([]fyne.CanvasObject, len(mImgs))
	zoomImgs := make([]fyne.CanvasObject, len(mImgs))
	for idx, mImg := range mImgs {
		imgCanvas, err := loadImage(bytes.NewBuffer(mImg))
		if err != nil {
			return errorLabel("load pdf error"), nil
		}
		imgCanvases[idx] = container.NewGridWrap(fyne.NewSize(400, 400), imgCanvas)
		zoomImgs[idx] = withZoom(imgCanvas)
	}
	if len(imgCanvases) == 0 {
		return errorLabel("load pdf error"), nil
	}

	player := NewPdfPlayer(imgCanvases...)
	zoomBtn := widget.NewButtonWithIcon("", theme.ViewFullScreenIcon(), func() {
		zoomPlayer := NewPdfPlayer(zoomImgs...)
		name := zoomImgs[0].(*fyne.Container).Objects[0].(*canvas.Image).Resource.Name()
		gui.addPageToTabs(name, zoomPlayer.Render())
	})
	return container.NewBorder(container.NewBorder(nil, nil, zoomBtn, nil), nil, nil, nil, player.Render()), nil
}
