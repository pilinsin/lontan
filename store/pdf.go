package store

import (
	"bytes"
	"errors"
	"io"
	"os"
	"strconv"

	bimg "github.com/h2non/bimg"
	pdfapi "github.com/pdfcpu/pdfcpu/pkg/api"
	proto "google.golang.org/protobuf/proto"

	pb "github.com/pilinsin/lontan/store/pb"
)

// convert to webp
func marshalPdfTopPage(pdfPath, outName string) ([]byte, error) {
	buf, err := bimg.Read(pdfPath)
	if err != nil {
		return nil, err
	}
	img, err := bimg.NewImage(buf).Convert(bimg.WEBP)
	if err != nil {
		return nil, err
	}

	pbImage := &pb.Image{
		Name: outName,
		Data: img,
	}
	return proto.Marshal(pbImage)
}

func encodePdf(filename, tmpname string) ([][]byte, error) {
	n, err := pdfapi.PageCountFile(tmpname)
	if err != nil {
		return nil, err
	}
	if n == 0 {
		return nil, errors.New("invalid pdf: pageCount == 0")
	}

	mImgs := make([][]byte, 0)
	idx := 0
	for {
		outName := filename + "_" + strconv.Itoa(idx) + ".webp"
		m, err := marshalPdfTopPage(tmpname, outName)
		if err != nil {
			return nil, err
		}
		mImgs = append(mImgs, m)

		n, err := pdfapi.PageCountFile(tmpname)
		if err != nil {
			return nil, err
		}
		if n == 1 {
			break
		}

		if err := pdfapi.RemovePagesFile(tmpname, "", []string{"1"}, nil); err != nil {
			return nil, err
		}
		idx++
	}

	return mImgs, nil
}

func EncodePdf(name string, r io.Reader) (io.Reader, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}
	tmpname := "pdf_temp.pdf"
	if err := os.WriteFile(tmpname, data, 0666); err != nil {
		return nil, err
	}
	defer os.Remove(tmpname)

	mImgs, err := encodePdf(name, tmpname)
	if err != nil {
		return nil, err
	}

	pbPdf := &pb.Pdf{
		Images: mImgs,
	}
	m, err := proto.Marshal(pbPdf)
	if err != nil {
		return nil, err
	}

	return bytes.NewBuffer(m), nil
}
