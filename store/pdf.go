package store

import(
	"errors"
	"bytes"
	"io"
	"os"
	"strconv"
	"path/filepath"

	pdfapi "github.com/pdfcpu/pdfcpu/pkg/api"
	bimg "github.com/h2non/bimg"
	proto "google.golang.org/protobuf/proto"

	pb "github.com/pilinsin/lontan/store/pb"
)

func pdfFrontPageToImage(pdfPath, outPath string) error{
	buf, err := bimg.Read(pdfPath)
	if err != nil{return err}
	img, err := bimg.NewImage(buf).Convert(bimg.WEBP)
	if err != nil{return err}

	return bimg.Write(outPath, img)
}

func imageFileToPbImageMarshal(filename string) ([]byte, error){
	data, err := os.ReadFile(filename)
	if err != nil{return nil, err}

	pbImage := &pb.Image{
		Name: filename,
		Data: data,
	}
	return proto.Marshal(pbImage)
}


func encodePdf(filename, tmpname string) ([][]byte, error){
	outDir := "converted"
	if err := os.Mkdir(outDir, 0755); err != nil{return nil, err}
	defer os.RemoveAll(outDir)

	n, err := pdfapi.PageCountFile(tmpname)
	if err != nil{return nil, err}
	if n == 0{return nil, errors.New("invalid pdf: pageCount == 0")}

	mImgs := make([][]byte, 0)
	idx := 0
	for{
		outName := filepath.Join(outDir, filename+"_"+strconv.Itoa(idx)+".webp")
		if err := pdfFrontPageToImage(tmpname, outName); err != nil{
			return nil, err
		}
		m, err := imageFileToPbImageMarshal(outName)
		if err != nil{return nil, err}
		mImgs = append(mImgs, m)

		n, err := pdfapi.PageCountFile(tmpname)
		if err != nil{return nil, err}
		if n == 1{break}

		if err := pdfapi.RemovePagesFile(tmpname, "", []string{"1"}, nil); err != nil{
			return nil, err
		}
		idx++
	}
	
	return mImgs, nil
}

func EncodePdf(name string, r io.Reader) (io.Reader, error){
	data, err := io.ReadAll(r)
	if err != nil{return nil, err}
	tmpname := "pdf_temp.pdf"
	if err := os.WriteFile(tmpname, data, 0666); err != nil{
		return nil, err
	}
	defer os.Remove(tmpname)

	mImgs, err := encodePdf(name, tmpname)
	if err != nil{return nil, err}

	pbPdf := &pb.Pdf{
		Images: mImgs,
	}
	m, err := proto.Marshal(pbPdf)
	if err != nil{return nil, err}

	return bytes.NewBuffer(m), nil
}
