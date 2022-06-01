package gui

import (
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type pdfPlayer struct {
	pdfObjs []fyne.CanvasObject
	idx     int
}

func NewPdfPlayer(pdfObjs ...fyne.CanvasObject) *pdfPlayer {
	if len(pdfObjs) == 0 {
		pdfObjs = []fyne.CanvasObject{errorLabel("no page")}
	}
	for idx := range pdfObjs {
		pdfObjs[idx].Hide()
	}
	return &pdfPlayer{pdfObjs, 0}
}

func (pdf *pdfPlayer) Render() fyne.CanvasObject {
	pdf.pdfObjs[pdf.idx].Show()

	var page *fyne.Container

	idxLabel := widget.NewLabel(pdf.pageIndex())
	upBtn := widget.NewButtonWithIcon("", theme.MenuDropUpIcon(), func() {
		if pdf.idx > 0 {
			pdf.pdfObjs[pdf.idx].Hide()
			pdf.idx--
			pdf.pdfObjs[pdf.idx].Show()
			page.Objects[0] = pdf.pdfObjs[pdf.idx]
			page.Refresh()
			idxLabel.SetText(pdf.pageIndex())
		}
	})
	dnBtn := widget.NewButtonWithIcon("", theme.MenuDropDownIcon(), func() {
		if pdf.idx < len(pdf.pdfObjs)-1 {
			pdf.pdfObjs[pdf.idx].Hide()
			pdf.idx++
			pdf.pdfObjs[pdf.idx].Show()
			page.Objects[0] = pdf.pdfObjs[pdf.idx]
			page.Refresh()
			idxLabel.SetText(pdf.pageIndex())
		}
	})
	menuBar := container.NewVBox(upBtn, idxLabel, dnBtn)

	page = container.NewBorder(nil, nil, menuBar, nil, pdf.pdfObjs[pdf.idx])
	return page
}
func (pdf *pdfPlayer) RenderTile() fyne.CanvasObject {
	for _, obj := range pdf.pdfObjs {
		obj.Show()
	}

	page := container.NewVBox(pdf.pdfObjs...)
	return container.NewMax(container.NewVScroll(page))
}
func (pdf *pdfPlayer) pageIndex() string {
	return strconv.Itoa(pdf.idx+1) + "/" + strconv.Itoa(len(pdf.pdfObjs))
}
