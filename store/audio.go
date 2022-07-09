package store

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"strings"

	"fyne.io/fyne/v2"
	ffmpeg "github.com/u2takey/ffmpeg-go"

	proto "google.golang.org/protobuf/proto"
	pb "github.com/pilinsin/lontan/store/pb"
	ipfs "github.com/pilinsin/p2p-verse/ipfs"
)


func EncodeAudio(r fyne.URIReadCloser, is ipfs.Ipfs) (io.Reader, error) {
	tmpDir, err := os.MkdirTemp("", "audio_convert")
	if err != nil{return nil, err}
	outName := filepath.Join(tmpDir, strings.TrimSuffix(r.URI().Name(), r.URI().Extension()))
	outSuffix := "%03d.mp3"
	strm := ffmpeg.Input(r.URI().Path()).Audio().
		Output(outName+outSuffix, ffmpeg.KwArgs{"c:a": "mp3", "f": "segment", "segment_time": 10})
	err = strm.OverWriteOutput().Run()
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(tmpDir)

	files, err := os.ReadDir(tmpDir)
	if err != nil{return nil, err}
	cids := make([]string, len(files))
	for idx, file := range files{
		f, err := os.Open(filepath.Join(tmpDir, file.Name()))
		if err != nil{return nil, err}
		defer f.Close()

		cid, err := is.AddReader(f)
		if err != nil{return nil, err}
		cids[idx] = cid
	}

	pbAudio := &pb.Audio{
		Cids: cids,
	}
	m, err := proto.Marshal(pbAudio)
	if err != nil{return nil, err}

	return bytes.NewBuffer(m), nil
}
