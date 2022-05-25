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
	for _, obj := range pdfObjs {
		obj.Hide()
	}
	return &pdfPlayer{pdfObjs, 0}
}
func (pdf *pdfPlayer) Render() fyne.CanvasObject {
	pdf.pdfObjs[pdf.idx].Show()

	screen := container.NewVBox(pdf.pdfObjs...)

	idxLabel := widget.NewLabel(pdf.pageIndex())
	upBtn := widget.NewButtonWithIcon("", theme.MenuDropUpIcon(), func() {
		if pdf.idx > 0 {
			pdf.pdfObjs[pdf.idx].Hide()
			pdf.idx--
			pdf.pdfObjs[pdf.idx].Show()
			idxLabel.SetText(pdf.pageIndex())
		}
	})
	dnBtn := widget.NewButtonWithIcon("", theme.MenuDropDownIcon(), func() {
		if pdf.idx < len(pdf.pdfObjs)-1 {
			pdf.pdfObjs[pdf.idx].Hide()
			pdf.idx++
			pdf.pdfObjs[pdf.idx].Show()
			idxLabel.SetText(pdf.pageIndex())
		}
	})
	menuBar := container.NewVBox(upBtn, idxLabel, dnBtn)

	return container.NewBorder(nil, nil, menuBar, nil, screen)
}
func (pdf *pdfPlayer) pageIndex() string {
	return strconv.Itoa(pdf.idx+1) + "/" + strconv.Itoa(len(pdf.pdfObjs))
}
