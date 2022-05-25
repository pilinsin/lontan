package store

import (
	//"errors"
	"io"
	//"os"
	"time"
	//"bytes"
	//proto "google.golang.org/protobuf/proto"
	//reisen "github.com/zergon321/reisen"

	pb "github.com/pilinsin/lontan/store/pb"
	ipfs "github.com/pilinsin/p2p-verse/ipfs"
)

//protoc --go_out=. *.proto

type audio struct {
	chunkCids  []string
	duration   time.Duration
	sampleRate int
}

func (a audio) ChunkCids() []string     { return a.chunkCids }
func (a audio) Duration() time.Duration { return a.duration }
func (a audio) SampleRate() int         { return a.sampleRate }

func encodeAudio(cids []string, duration time.Duration, sr int) *pb.Audio {
	return &pb.Audio{
		ChunkCids:  cids,
		Duration:   int64(duration),
		SampleRate: int64(sr),
	}
}
func DecodeAudio(pbm *pb.Audio) *audio {
	return &audio{
		chunkCids:  pbm.GetChunkCids(),
		duration:   time.Duration(pbm.GetDuration()),
		sampleRate: int(pbm.GetSampleRate()),
	}
}

type chunkAudio struct {
	samples [][2]float64
}

func (c chunkAudio) Samples() [][2]float64 { return c.samples }

func encodeChunkAudio(samples [][2]float64) *pb.ChunkAudio {
	pbss := make([]*pb.AudioSample, len(samples))
	for idx, sample := range samples {
		pbss[idx] = encodeAudioSample(sample)
	}

	return &pb.ChunkAudio{
		Samples: pbss,
	}
}
func DecodeChunkAudio(pbc *pb.ChunkAudio) *chunkAudio {
	pbss := pbc.GetSamples()
	samples := make([][2]float64, len(pbss))
	for idx, pbs := range pbss {
		samples[idx] = decodeAudioSample(pbs)
	}

	return &chunkAudio{samples}
}

/*
func readAudioOnly(media *reisen.Media) (<-chan [2]float64, chan error, error){
	sampleBuffer := make(chan [2]float64, sampleBufferSize)
	errs := make(chan error)

	if err := media.OpenDecode(); err != nil{
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
		audioStream.Close()
		media.CloseDecode()
		close(sampleBuffer)
		close(errs)
	}()

	return sampleBuffer, errs, nil
}
*/

func encodeAudioFile(fname string, is ipfs.Ipfs) ([]byte, error) {
	return nil, nil
	/*
		media, err := reisen.NewMedia(fname)
		if err != nil{return nil, err}
		defer media.Close()

		dur, err := media.Duration()
		if err != nil{return nil, err}
		ast := media.AudioStreams()
		if len(ast) == 0{return nil, errors.New("no stream")}
		sr := ast[0].SampleRate()
		if sr == 0{
			sr = 44100
		}

		sampleBuffer, errs, err := readAudioOnly(media)
		if err != nil{return nil, err}


		chunkCids := make([]string, 0)
		chunkSampleSize := sr*10
		for{
			if err := isErrs(errs); err != nil{return nil, err}

			samples := make([][2]float64, 0)
			for sample := range sampleBuffer{
				samples = append(samples, sample)
				if len(samples) >= chunkSampleSize{
					break
				}
			}

			mch, err := proto.Marshal(encodeChunkAudio(samples))
			if err != nil{return nil, err}
			cid, err := is.Add(mch)
			if err != nil{return nil, err}
			chunkCids = append(chunkCids, cid)
		}

		ma, err := proto.Marshal(encodeAudio(chunkCids, dur, sr))
		if err != nil{return nil, err}

		return ma, nil
	*/
}
func EncodeAudio(r io.Reader, is ipfs.Ipfs) (io.Reader, error) {
	return nil, nil
	/*
		tmpDir := "media_cache"
		tmpFile, err := os.CreateTemp(tmpDir, "*****")
		if err != nil{return nil, err}
		if _, err := tmpFile.ReadFrom(r); err != nil{return nil, err}
		fname := tmpFile.Name()
		tmpFile.Close()

		m, err := encodeAudioFile(fname, is)
		if err != nil{return nil, err}

		os.RemoveAll(tmpDir)
		return bytes.NewReader(m), nil
	*/
}
