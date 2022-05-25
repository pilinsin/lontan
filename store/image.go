package store

import (
	"bytes"
	proto "google.golang.org/protobuf/proto"
	"io"

	pb "github.com/pilinsin/lontan/store/pb"
)

func EncodeImage(name string, r io.Reader) (io.Reader, error) {
	data, err := io.ReadAll(r)
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
