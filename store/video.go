package store

import (
	//"errors"
	"io"
	//"os"
	"encoding/binary"
	"image"
	"time"
	//"bytes"
	//proto "google.golang.org/protobuf/proto"
	//reisen "github.com/zergon321/reisen"

	pb "github.com/pilinsin/lontan/store/pb"
	ipfs "github.com/pilinsin/p2p-verse/ipfs"
)

//protoc --go_out=. *.proto

const (
	frameBufferSize = 1024

	channelCount     = 2
	bitDepth         = 8
	sampleBufferSize = 32 * channelCount * bitDepth * 1024
)

type video struct {
	chunkCids  []string
	duration   time.Duration
	frameRate  int
	sampleRate int
}

func (v video) ChunkCids() []string     { return v.chunkCids }
func (v video) Duration() time.Duration { return v.duration }
func (v video) FrameRate() int          { return v.frameRate }
func (v video) SampleRate() int         { return v.sampleRate }

func encodeVideo(cids []string, duration time.Duration, fps, sr int) *pb.Video {
	return &pb.Video{
		ChunkCids:  cids,
		Duration:   int64(duration),
		FrameRate:  int64(fps),
		SampleRate: int64(sr),
	}
}
func DecodeVideo(pbm *pb.Video) *video {
	return &video{
		chunkCids:  pbm.GetChunkCids(),
		duration:   time.Duration(pbm.GetDuration()),
		frameRate:  int(pbm.GetFrameRate()),
		sampleRate: int(pbm.GetSampleRate()),
	}
}

type chunkVideo struct {
	frames  []*image.RGBA
	samples [][2]float64
}

func (c chunkVideo) Frames() []*image.RGBA { return c.frames }
func (c chunkVideo) Samples() [][2]float64 { return c.samples }

func encodeChunkVideo(frames []*image.RGBA, samples [][2]float64) *pb.ChunkVideo {
	pbfs := make([]*pb.ImageRGBA, len(frames))
	for idx, frame := range frames {
		pbfs[idx] = encodeRGBA(frame)
	}

	pbss := make([]*pb.AudioSample, len(samples))
	for idx, sample := range samples {
		pbss[idx] = encodeAudioSample(sample)
	}

	return &pb.ChunkVideo{
		Frames:  pbfs,
		Samples: pbss,
	}
}
func DecodeChunkVideo(pbc *pb.ChunkVideo) *chunkVideo {
	pbfs := pbc.GetFrames()
	frames := make([]*image.RGBA, len(pbfs))
	for idx, pbf := range pbfs {
		frames[idx] = decodeRGBA(pbf)
	}

	pbss := pbc.GetSamples()
	samples := make([][2]float64, len(pbss))
	for idx, pbs := range pbss {
		samples[idx] = decodeAudioSample(pbs)
	}

	return &chunkVideo{frames, samples}
}

func u8ArrayToU32Array(u8s []uint8) []uint32 {
	u32s := make([]uint32, len(u8s))
	for idx := range u8s {
		u32s[idx] = uint32(u8s[idx])
	}
	return u32s
}
func u32ArrayToU8Array(u32s []uint32) []uint8 {
	u8s := make([]uint8, len(u32s))
	for idx := range u32s {
		u8s[idx] = uint8(u32s[idx])
	}
	return u8s
}

func encodeRGBA(img *image.RGBA) *pb.ImageRGBA {
	return &pb.ImageRGBA{
		Pix:    u8ArrayToU32Array(img.Pix),
		Stride: int64(img.Stride),
		Rect:   encodeRectangle(img.Rect),
	}
}
func encodeRectangle(rect image.Rectangle) *pb.Rectangle {
	return &pb.Rectangle{
		Min: encodePoint(rect.Min),
		Max: encodePoint(rect.Max),
	}
}
func encodePoint(pt image.Point) *pb.Point {
	return &pb.Point{
		X: int64(pt.X),
		Y: int64(pt.Y),
	}
}
func decodeRGBA(img *pb.ImageRGBA) *image.RGBA {
	return &image.RGBA{
		Pix:    u32ArrayToU8Array(img.GetPix()),
		Stride: int(img.GetStride()),
		Rect:   decodeRectangle(img.GetRect()),
	}
}
func decodeRectangle(rect *pb.Rectangle) image.Rectangle {
	return image.Rectangle{
		Min: decodePoint(rect.GetMin()),
		Max: decodePoint(rect.GetMax()),
	}
}
func decodePoint(pt *pb.Point) image.Point {
	return image.Point{
		X: int(pt.GetX()),
		Y: int(pt.GetY()),
	}
}

func encodeAudioSample(sample [2]float64) *pb.AudioSample {
	return &pb.AudioSample{
		Left:  sample[0],
		Right: sample[1],
	}
}
func decodeAudioSample(sample *pb.AudioSample) [2]float64 {
	return [2]float64{sample.GetLeft(), sample.GetRight()}
}

func addError(errs chan error, err error) {
	go func(err error) {
		errs <- err
	}(err)
}
func readAudioSample(r io.Reader, errs chan error) [2]float64 {
	sample := [2]float64{0, 0}
	var result float64

	if err := binary.Read(r, binary.LittleEndian, &result); err != nil {
		addError(errs, err)
	}
	sample[0] = result

	if err := binary.Read(r, binary.LittleEndian, &result); err != nil {
		addError(errs, err)
	}
	sample[1] = result

	return sample
}

/*
func read(media *reisen.Media) (<-chan *image.RGBA, <-chan [2]float64, chan error, error){
	frameBuffer := make(chan *image.RGBA, frameBufferSize)
	sampleBuffer := make(chan [2]float64, sampleBufferSize)
	errs := make(chan error)

	if err := media.OpenDecode(); err != nil{
		return nil, nil, nil, err
	}
	videoStream := media.VideoStreams()[0]
	if err := videoStream.Open(); err != nil{
		return nil, nil, nil, err
	}
	audioStream := media.AudioStreams()[0]
	if err := audioStream.Open(); err != nil{
		return nil, nil, nil, err
	}

	go func(){
		for{
			packet, ok, err := media.ReadPacket()
			if err != nil{
				addError(errs, err)
			}
			if !ok{break}

			switch packet.Type() {
			case reisen.StreamVideo:
				vs := media.Streams()[packet.StreamIndex()].(*reisen.VideoStream)
				vFrame, ok, err := vs.ReadVideoFrame()
				if err != nil{
					addError(errs, err)
				}
				if !ok{break}
				if vFrame == nil{continue}

				frameBuffer <- vFrame.Image()

			case reisen.StreamAudio:
				as := media.Streams()[packet.StreamIndex()].(*reisen.AudioStream)
				aFrame, ok, err := as.ReadAudioFrame()
				if err != nil{
					addError(errs, err)
				}
				if !ok{break}
				if aFrame == nil{continue}

				reader := bytes.NewReader(audioFrame.Data())
				for reader.Len() > 0{
					sampleBuffer <- readAudioSample(reader, errs)
				}
			}
		}
		videoStream.Close()
		audioStream.Close()
		media.CloseDecode()
		close(frameBuffer)
		close(sampleBuffer)
		close(errs)
	}()

	return frameBuffer, sampleBuffer, errs, nil
}
*/

func isErrs(errs chan error) error {
	select {
	case err, ok := <-errs:
		if ok {
			return err
		}
	default:
	}

	return nil
}
func encodeVideoFile(fname string, is ipfs.Ipfs) ([]byte, error) {
	return nil, nil
	/*
		media, err := reisen.NewMedia(fname)
		if err != nil{return nil, err}
		defer media.Close()

		dur, err := media.Duration()
		if err != nil{return nil, err}
		vst := media.VideoStreams()
		if len(vst) == 0{return nil, errors.New("no stream")}
		fps, _ := vst[0].FrameRate()
		if fps == 0{
			fps = 60
		}
		ast := media.AudioStreams()
		if len(ast) == 0{return nil, errors.New("no stream")}
		sr := ast[0].SampleRate()
		if sr == 0{
			sr = 44100
		}

		frameBuffer, sampleBuffer, errs, err := read(media)
		if err != nil{return nil, err}


		chunkCids := make([]string, 0)
		chunkFrameSize := fps*10
		chunkSampleSize := sr*10
		for{
			if err := isErrs(errs); err != nil{return nil, err}

			frames := make([]*image.RGBA, 0)
			for frame := range frameBuffer{
				frames = append(frames, frame)
				if len(frames) >= chunkFrameSize{
					break
				}
			}

			samples := make([][2]float64, 0)
			for sample := range sampleBuffer{
				samples = append(samples, sample)
				if len(samples) >= chunkSampleSize{
					break
				}
			}

			mch, err := proto.Marshal(encodeChunkVideo(frames, samples))
			if err != nil{return nil, err}
			cid, err := is.Add(mch)
			if err != nil{return nil, err}
			chunkCids = append(chunkCids, cid)
		}

		mv, err := proto.Marshal(encodeVideo(chunkCids, dur, fps, sr))
		if err != nil{return nil, err}

		return mv, nil
	*/
}
func EncodeVideo(r io.Reader, is ipfs.Ipfs) (io.Reader, error) {
	return nil, nil
	/*
		tmpDir := "media_cache"
		tmpFile, err := os.CreateTemp(tmpDir, "*****")
		if err != nil{return nil, err}
		if _, err := tmpFile.ReadFrom(r); err != nil{return nil, err}
		fname := tmpFile.Name()
		tmpFile.Close()

		m, err := encodeVideoFile(fname, is)
		if err != nil{return nil, err}

		os.RemoveAll(tmpDir)
		return bytes.NewReader(m), nil
	*/
}
