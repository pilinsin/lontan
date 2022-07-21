package gui

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"time"

	mp3 "github.com/hajimehoshi/go-mp3"
	oto "github.com/hajimehoshi/oto/v2"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	gutil "github.com/pilinsin/lontan/gui/util"
	store "github.com/pilinsin/lontan/store"
	pb "github.com/pilinsin/lontan/store/pb"
	ipfs "github.com/pilinsin/p2p-verse/ipfs"
	proto "google.golang.org/protobuf/proto"
)

var otoCtx *oto.Context

func init() {
	ctx, readyChan, _ := oto.NewContext(store.SampleRate, store.NumChannel, store.ByteDepth)
	<-readyChan
	otoCtx = ctx
}

type audioDecoder struct {
	*mp3.Decoder
	readIdx int64
}

func newAudioDecoder(aData []byte) (*audioDecoder, error) {
	mp3Dec, err := mp3.NewDecoder(bytes.NewReader(aData))
	if err != nil {
		return nil, err
	}

	return &audioDecoder{
		Decoder: mp3Dec,
	}, nil
}
func (dec *audioDecoder) Read(buf []byte) (int, error) {
	n, err := dec.Decoder.Read(buf)
	if err != nil {
		return 0, err
	}
	dec.readIdx += int64(n)
	return n, nil
}
func (dec *audioDecoder) Seek(offset int64, whence int) (int64, error) {
	ofst, err := dec.Decoder.Seek(offset, whence)
	if err != nil {
		return 0, err
	}

	dec.readIdx = ofst
	return ofst, nil

}

type iAudioPlayer interface {
	gutil.IPlayer
}
type otoPlayer struct {
	oto.Player
	src *audioDecoder
}

//single r can be used by only single Player
func newOtoPlayer(dec *audioDecoder) (*otoPlayer, error) {
	player := otoCtx.NewPlayer(dec)
	player.(oto.BufferSizeSetter).SetBufferSize(store.ByteRate * 5)
	return &otoPlayer{
		Player: player,
		src:    dec,
	}, nil
}
func (op *otoPlayer) IsPausing() bool {
	return !op.IsPlaying()
}
func (op *otoPlayer) PlayedTime() (time.Duration, error) {
	playedSize := (op.src.readIdx + 1) - int64(op.UnplayedBufferSize())
	d := float64(playedSize) / float64(store.ByteRate)
	return time.ParseDuration(fmt.Sprintf("%vs", d))
}
func (op *otoPlayer) Wait(d time.Duration) {
	if d == 0 {
		return
	}

	op.Pause()
	time.Sleep(d)
	op.Play()
}
func (op *otoPlayer) Seek(offset int64, whence int) (int64, error) {
	seeker, ok := op.Player.(io.Seeker)
	if !ok {
		return -1, errors.New("io.Seeker is not implemented")
	}

	return seeker.Seek(offset, whence)
}

type audioPlayer struct {
	ctx    context.Context
	cancel func()

	cid string
	is  ipfs.Ipfs

	player  iAudioPlayer
	timeBar gutil.ITimeBar
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
	if err := proto.Unmarshal(m, pbAudio); err != nil {
		return err
	}

	dec, err := newAudioDecoder(pbAudio.GetData())
	if err != nil {
		return err
	}
	player, err := newOtoPlayer(dec)
	if err != nil {
		return err
	}
	ap.player = player

	ap.timeBar = gutil.NewTimeBar(pbAudio.GetSecond())

	ctx, cancel := context.WithCancel(context.Background())
	ap.ctx = ctx
	ap.cancel = cancel

	return nil
}
func (ap *audioPlayer) Close() error {
	if ap.IsPlaying() && !ap.IsPausing() {
		ap.Pause()
	}

	ap.cancel()
	err1 := ap.player.Close()
	err2 := ap.timeBar.Close()
	if err1 != nil {
		return err1
	} else {
		return err2
	}
}

func (ap *audioPlayer) IsPlaying() bool {
	return ap.player.IsPlaying() && ap.timeBar.IsPlaying()
}
func (ap *audioPlayer) IsPausing() bool {
	return ap.player.IsPausing() && ap.timeBar.IsPausing()
}
func (ap *audioPlayer) Play() {
	ap.timeBar.Play()
	ap.player.Play()
}
func (ap *audioPlayer) Pause() {
	ap.timeBar.Pause()
	ap.player.Pause()
}

func (ap *audioPlayer) SyncTime() {
	go func() {
		ticker := time.NewTicker(time.Second * 30)
		defer ticker.Stop()

		for {
			select {
			case <-ap.ctx.Done():
				return
			case <-ticker.C:
				if !ap.IsPlaying() && ap.IsPausing() {
					continue
				}

				aTime, err1 := ap.player.PlayedTime()
				tTime, err2 := ap.timeBar.PlayedTime()
				if err1 != nil || err2 != nil {
					return
				}
				min := gutil.MinTime(aTime, tTime)
				ap.player.Wait(aTime - min)
				ap.timeBar.Wait(tTime - min)
			}
		}

	}()
}

func (ap *audioPlayer) Render() (fyne.CanvasObject, gutil.Closer) {
	var btns fyne.CanvasObject
	var obj fyne.CanvasObject

	var playBtn *widget.Button
	playBtn = widget.NewButtonWithIcon("", theme.MediaPlayIcon(), func() {
		if ap.player.IsPlaying() {
			ap.Pause()
			playBtn.SetIcon(theme.MediaPlayIcon())
		} else {
			ap.Play()
			playBtn.SetIcon(theme.MediaPauseIcon())
		}
	})

	resetBtn := widget.NewButtonWithIcon("", theme.MediaSkipPreviousIcon(), func() {
		ap.Close()
		ap.init()
		obj.(*fyne.Container).Objects[0] = ap.timeBar.Render()
		obj.Refresh()
		ap.SyncTime()
		playBtn.SetIcon(theme.MediaPlayIcon())
	})

	ap.SyncTime()
	btns = container.NewHBox(playBtn, resetBtn)
	obj = container.NewBorder(nil, nil, btns, nil, ap.timeBar.Render())
	return obj, ap.Close
}
