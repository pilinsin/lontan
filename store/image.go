package store

import (
	"io"
	"bytes"

	bimg "github.com/h2non/bimg"
	proto "google.golang.org/protobuf/proto"

	pb "github.com/pilinsin/lontan/store/pb"
)
// convert to webp
func EncodeImage(name string, r io.Reader) (io.Reader, error) {
	b, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}
	data, err := bimg.NewImage(b).Convert(bimg.WEBP)
	if err != nil {
		return nil, err
	}

	pbImage := &pb.Image{
		Name: name,
		Data: data,
	}
	m, err := proto.Marshal(pbImage)
	if err != nil {
		return nil, err
	}

	return bytes.NewBuffer(m), nil
}
