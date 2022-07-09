package gui

import (
	"fmt"
	"io"
	"errors"

	mp3 "github.com/hajimehoshi/go-mp3"
	oto "github.com/hajimehoshi/oto/v2"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	proto "google.golang.org/protobuf/proto"
	pb "github.com/pilinsin/lontan/store/pb"
	ipfs "github.com/pilinsin/p2p-verse/ipfs"
)

type multiChunkDecoder struct{
	is ipfs.Ipfs
	cids []string
	loadIdx int
	bufSize int
	ch chan chan byte
	sampleRate int

	playIdx int
	playBuffer chan byte
}
func newMultiChunkDecoder(cids []string, is ipfs.Ipfs) (*multiChunkDecoder, error){
	if len(cids) == 0{return nil, errors.New("no data input")}

	dec := &multiChunkDecoder{
		is: is,
		cids: cids,
		playIdx: -1,
	}
	r, err := is.GetReader(cids[0])
	if err != nil{return nil, err}
	
	mp3Dec, err := mp3.NewDecoder(r)
	if err != nil{return nil, err}
	dec.sampleRate = mp3Dec.SampleRate()

	//bufSize := sampleRate * nChannels(2) * bitDepth(2) * sec(10) + alpha
	dec.bufSize = dec.sampleRate * 50
	dec.ch = make(chan chan byte, 4)
	dec.playBuffer = make(chan byte, dec.bufSize * 2)
	dec.initLoad()

	return dec, nil
}

func (dec *multiChunkDecoder) initLoad() error{
	for{
		if len(dec.ch) > 2{
			return nil
		}

		ch := make(chan byte, dec.bufSize)
		dec.ch <- ch
			
		var buf []byte
		var err error
		for {
			if dec.isFullyLoaded(){
				return nil
			}
			buf, err = dec.loadAt()
			if err == nil{break}
			if err != io.EOF{
				return err
			}
		}

		if len(buf) > 0{
			for _, b := range buf{
				ch <- b
			}
		}
		close(ch)

		if dec.isFullyLoaded(){
			close(dec.ch)
			return nil
		}
	}
}
func (dec *multiChunkDecoder) isFullyLoaded() bool{
	return dec.loadIdx >= len(dec.cids)
}
func (dec *multiChunkDecoder) Load(){
	if dec.isFullyLoaded(){return}

	go func(){
		if len(dec.ch) < 3{
			ch := make(chan byte, dec.bufSize)
			dec.ch <- ch

			var buf []byte
			var err error
			for {
				if dec.isFullyLoaded(){return}

				buf, err = dec.loadAt()
				if err == nil{break}
				if err != io.EOF{
					return
				}
			}

			if len(buf) > 0{
				for _, b := range buf{
					ch <- b
				}
			}
			close(ch)

			if dec.isFullyLoaded(){
				close(dec.ch)
			}
		}
	}()
}
func (dec *multiChunkDecoder) loadAt() ([]byte, error){
	if dec.loadIdx < 0{
		return nil, errors.New("invalid idx")
	}
	if dec.isFullyLoaded(){return nil, nil}

	r, err := dec.is.GetReader(dec.cids[dec.loadIdx])
	if err != nil{return nil, err}
	mp3Dec, err := mp3.NewDecoder(r)
	if err != nil{return nil, err}

	buf := make([]byte, dec.bufSize)
	n, err := io.ReadFull(mp3Dec, buf)
	if err == io.ErrUnexpectedEOF{
		err = nil
	}
	fmt.Println(err)
	if err != nil && err != io.EOF{return nil, err}

	dec.loadIdx++
	return buf[:n], err
}

func (dec *multiChunkDecoder) readPlayBufferFromChan() error{
	var ok bool
	if dec.playBuffer == nil || len(dec.playBuffer) == 0{
		dec.playBuffer, ok = <-dec.ch
		if !ok{return io.EOF}

		dec.playIdx++
	}
	return nil
}
func (dec *multiChunkDecoder) Read(buf []byte) (int, error){
	for idx := 0; idx < len(buf); idx++{
		if err := dec.readPlayBufferFromChan(); err != nil{
			return idx, err
		}

		val, ok := <-dec.playBuffer
		if !ok{
			return idx, io.EOF
		}
		buf[idx] = val
	}

	dec.Load()
	return len(buf), nil
}
func (dec *multiChunkDecoder) Seek(offset int64, whence int) (int64, error){
	ofst := int(offset)
	if ofst < 0 || ofst >= len(dec.cids){
		return -1, errors.New("invalid offset")
	}
	switch whence {
	case io.SeekCurrent:
		ofst += dec.playIdx
	case io.SeekEnd:
		ofst = len(dec.cids) - ofst
	default:
		//io.SeekStart
	}

	if dec.loadIdx < len(dec.cids){
		close(dec.ch)
	}

	dec.playIdx = ofst - 2
	dec.loadIdx = ofst - 1
	dec.ch = make(chan chan byte, 4)
	dec.initLoad()

	return int64(ofst), nil
}

//single r can be used by only single Player
func newOtoPlayer(dec *multiChunkDecoder) (oto.Player, error) {
	otoCtx, readyChan, err := oto.NewContext(dec.sampleRate, 2, 2)
	if err != nil {
		return nil, err
	}
	<-readyChan

	return otoCtx.NewPlayer(dec), nil
}

type audioPlayer struct {
	cid    string
	is     ipfs.Ipfs
	player oto.Player
}

func NewAudioPlayer(cid string, is ipfs.Ipfs) (*audioPlayer, error) {
	ap := &audioPlayer{
		cid: cid,
		is:  is,
	}
	if err := ap.init(); err != nil {
		return nil, err
	}

	return ap, nil
}
func (ap *audioPlayer) init() error {
	m, err := ap.is.Get(ap.cid)
	if err != nil {
		return err
	}
	pbAudio := &pb.Audio{}
	if err := proto.Unmarshal(m, pbAudio); err != nil{return err}
	dec, err := newMultiChunkDecoder(pbAudio.GetCids(), ap.is)
	if err != nil{return err}

	player, err := newOtoPlayer(dec)
	if err != nil {
		return err
	}
	ap.player = player
	return nil
}

func (ap *audioPlayer) Render() fyne.CanvasObject {
	var playBtn *widget.Button
	playBtn = widget.NewButtonWithIcon("", theme.MediaPlayIcon(), func() {
		if ap.player.IsPlaying() {
			ap.player.Pause()
			playBtn.SetIcon(theme.MediaPlayIcon())
		} else {
			ap.player.Play()
			playBtn.SetIcon(theme.MediaPauseIcon())
		}
	})
	resetBtn := widget.NewButtonWithIcon("", theme.MediaSkipPreviousIcon(), func() {
		ap.player.Pause()
		ap.player.Close()
		ap.init()
		playBtn.SetIcon(theme.MediaPlayIcon())
	})

	return container.NewHBox(playBtn, resetBtn)
}
