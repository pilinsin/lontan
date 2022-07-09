package store

import (
	"bytes"
	"io"

	"fyne.io/fyne/v2"
	bimg "github.com/h2non/bimg"
	proto "google.golang.org/protobuf/proto"

	pb "github.com/pilinsin/lontan/store/pb"
)

// convert to webp
func EncodeImage(r fyne.URIReadCloser) (io.Reader, error) {
	b, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}
	data, err := bimg.NewImage(b).Convert(bimg.WEBP)
	if err != nil {
		return nil, err
	}

	pbImage := &pb.Image{
		Name: r.URI().Name(),
		Data: data,
	}
	m, err := proto.Marshal(pbImage)
	if err != nil {
		return nil, err
	}

	return bytes.NewBuffer(m), nil
}
