package store

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"strings"

	"fyne.io/fyne/v2"
	"github.com/hajimehoshi/go-mp3"
	ffmpeg "github.com/u2takey/ffmpeg-go"

	pb "github.com/pilinsin/lontan/store/pb"
	ipfs "github.com/pilinsin/p2p-verse/ipfs"
	proto "google.golang.org/protobuf/proto"
)

func EncodeAudio(r fyne.URIReadCloser, is ipfs.Ipfs) (io.Reader, error) {
	tmpDir, err := os.MkdirTemp("", "audio_convert")
	if err != nil {
		return nil, err
	}
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
	if err != nil {
		return nil, err
	}
	cids := make([]string, len(files))
	var second int64
	for idx, file := range files {
		f, err := os.Open(filepath.Join(tmpDir, file.Name()))
		if err != nil {
			return nil, err
		}
		defer f.Close()

		mp3Dec, err := mp3.NewDecoder(f)
		if err != nil {
			return nil, err
		}
		second += mp3Dec.Length() / (int64(mp3Dec.SampleRate()) * 4)

		cid, err := is.AddReader(f)
		if err != nil {
			return nil, err
		}
		cids[idx] = cid
	}

	pbAudio := &pb.Audio{
		Cids:   cids,
		Second: second,
	}
	m, err := proto.Marshal(pbAudio)
	if err != nil {
		return nil, err
	}

	return bytes.NewBuffer(m), nil
}
