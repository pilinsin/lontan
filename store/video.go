package store

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"

	"fyne.io/fyne/v2"
	gocv "gocv.io/x/gocv"

	pb "github.com/pilinsin/lontan/store/pb"
	ffmpeg "github.com/u2takey/ffmpeg-go"
	proto "google.golang.org/protobuf/proto"
)

const VideoW = 480
const VideoH = 270

func getVideoFpsAndLength(r fyne.URIReadCloser) (float64, int64, float64, error) {
	vc, err := gocv.VideoCaptureFile(r.URI().Path())
	if err != nil {
		return -1, -1, -1, err
	}
	fps := vc.Get(gocv.VideoCaptureFPS)
	nFrames := vc.Get(gocv.VideoCaptureFrameCount)
	sec := nFrames / fps
	vc.Close()

	return fps, int64(nFrames), sec, nil
}

func encodeVideo(r fyne.URIReadCloser) (string, error) {
	fileName := strings.TrimSuffix(r.URI().Name(), r.URI().Extension())
	f, err := os.CreateTemp(exeDir(), fileName+"_tmp_convert*"+r.URI().Extension())
	if err != nil {
		return "", err
	}
	f.Close()

	strm := ffmpeg.Input(r.URI().Path()).Video().
		Output(f.Name(), ffmpeg.KwArgs{
			"vf": fmt.Sprintf("scale=%dx%d:flags=lanczos", VideoW, VideoH),
		})
	if err := strm.OverWriteOutput().Run(); err != nil {
		return "", err
	}

	return f.Name(), nil
}

func EncodeVideo(r fyne.URIReadCloser) (io.Reader, error) {
	fps, nFrames, sec, err := getVideoFpsAndLength(r)
	if err != nil {
		return nil, err
	}

	encAudioName, err := encodeAudio(r)
	if err != nil {
		return nil, err
	}
	defer os.Remove(encAudioName)
	ma, err := os.ReadFile(encAudioName)
	if err != nil {
		return nil, err
	}

	encVideoName, err := encodeVideo(r)
	if err != nil {
		return nil, err
	}
	defer os.Remove(encVideoName)
	mv, err := os.ReadFile(encVideoName)
	if err != nil {
		return nil, err
	}

	pbVideo := &pb.Video{
		Audio:      ma,
		Video:      mv,
		FrameRate:  fps,
		FrameCount: nFrames,
		Second:     sec,
	}
	m, err := proto.Marshal(pbVideo)
	if err != nil {
		return nil, err
	}

	return bytes.NewBuffer(m), nil
}
