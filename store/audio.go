package store

import (
	"bytes"
	"io"
	"os"
	"strings"

	"fyne.io/fyne/v2"
	"github.com/hajimehoshi/go-mp3"
	ffmpeg "github.com/u2takey/ffmpeg-go"

	pb "github.com/pilinsin/lontan/store/pb"
	proto "google.golang.org/protobuf/proto"
)

const (
	NumChannel     = 2
	ByteDepth      = 2
	BytesPerSample = NumChannel * ByteDepth
	SampleRate     = 44100
	ByteRate       = SampleRate * BytesPerSample
)

func encodeAudio(r fyne.URIReadCloser) (string, error) {
	fileName := strings.TrimSuffix(r.URI().Name(), r.URI().Extension())
	f, err := os.CreateTemp(exeDir(), fileName+"_tmp_convert*.mp3")
	if err != nil {
		return "", err
	}
	f.Close()

	strm := ffmpeg.Input(r.URI().Path()).Audio().
		Output(f.Name(), ffmpeg.KwArgs{
			"c:a": "mp3",
			"ac":  NumChannel,
			"ar":  SampleRate,
		})
	if err := strm.OverWriteOutput().Run(); err != nil {
		return "", err
	}

	return f.Name(), nil
}
func EncodeAudio(r fyne.URIReadCloser) (io.Reader, error) {
	encodedName, err := encodeAudio(r)
	defer os.Remove(encodedName)
	if err != nil {
		return nil, err
	}

	f, err := os.Open(encodedName)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	mp3Dec, err := mp3.NewDecoder(f)
	if err != nil {
		return nil, err
	}
	second := float64(mp3Dec.Length()) / float64((mp3Dec.SampleRate())*4)

	data, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}

	pbAudio := &pb.Audio{
		Data:   data,
		Second: second,
	}
	m, err := proto.Marshal(pbAudio)
	if err != nil {
		return nil, err
	}

	return bytes.NewBuffer(m), nil
}
