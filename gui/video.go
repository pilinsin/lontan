package gui

import (
	"context"
	"errors"
	"fmt"
	"image"
	"io"
	"os"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	gutil "github.com/pilinsin/lontan/gui/util"
	"github.com/pilinsin/lontan/store"
	pb "github.com/pilinsin/lontan/store/pb"
	ipfs "github.com/pilinsin/p2p-verse/ipfs"
	"gocv.io/x/gocv"
	proto "google.golang.org/protobuf/proto"
)

type videoDecoder struct {
	srcName    string
	frameCount int64
	readIdx    int64

	vc *gocv.VideoCapture
}

func newVideoDecoder(vData []byte, frameCount int64) (*videoDecoder, error) {
	f, err := os.CreateTemp("", ".*_tmp")
	if err != nil {
		return nil, err
	}
	defer f.Close()
	if _, err := f.Write(vData); err != nil {
		return nil, err
	}

	dec := &videoDecoder{
		srcName:    f.Name(),
		frameCount: frameCount,
	}

	if err := dec.init(); err != nil {
		return nil, err
	}
	return dec, nil
}
func (dec *videoDecoder) init() error {
	if dec.vc != nil {
		dec.vc.Close()
	}

	vc, err := gocv.VideoCaptureFile(dec.srcName)
	if err != nil {
		return err
	}
	dec.vc = vc

	dec.readIdx = 0
	return nil
}
func (dec *videoDecoder) Close() {
	dec.vc.Close()
	os.Remove(dec.srcName)
}

func (dec *videoDecoder) Read(buf []image.Image) (int, error) {
	mat := gocv.NewMat()
	for idx := 0; idx < len(buf); idx++ {
		if ok := dec.vc.Read(&mat); !ok {
			dec.readIdx += int64(idx)
			return idx, io.EOF
		}
		img, _ := mat.ToImage()
		buf[idx] = img
	}

	dec.readIdx += int64(len(buf))
	return len(buf), nil
}

func (dec *videoDecoder) seek(n int64) error {
	if n < 0 || n > dec.frameCount {
		return errors.New("invalid seek count n")
	}

	if n < dec.readIdx+1 {
		if err := dec.init(); err != nil {
			return err
		}
	} else {
		n -= (dec.readIdx + 1)
	}

	if n == 0 {
		return nil
	}

	mat := gocv.NewMat()
	for idx := 0; int64(idx) <= n; idx++ {
		if ok := dec.vc.Read(&mat); !ok {
			return io.EOF
		}
	}
	return nil
}
func (dec *videoDecoder) Seek(offset int64, whence int) (int64, error) {
	if offset == 0 && whence == io.SeekStart {
		return dec.readIdx, nil
	}

	switch whence {
	case io.SeekCurrent:
		offset += dec.readIdx
	case io.SeekEnd:
		offset += dec.frameCount
	default:
		//io.SeekStart
	}

	err := dec.seek(offset)
	dec.readIdx = offset
	return offset, err
}

type iVideoPlayer interface {
	gutil.IPlayer
}

type video struct {
	ctx    context.Context
	cancel func()

	src       *videoDecoder
	buf       chan image.Image
	duration  time.Duration
	fps       float64
	isPlaying bool
	isPausing bool
	pauseCh   chan bool
	unPauseCh chan bool

	screen *canvas.Image
}

func newVideo(vr *videoDecoder, fps float64, screen *canvas.Image) (*video, error) {
	dur, err := time.ParseDuration(fmt.Sprintf("%vs", 1.0/fps))
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithCancel(context.Background())

	v := &video{
		ctx:       ctx,
		cancel:    cancel,
		src:       vr,
		buf:       make(chan image.Image, int(fps+1)*60),
		duration:  dur,
		fps:       fps,
		pauseCh:   make(chan bool, 1),
		unPauseCh: make(chan bool, 1),
		screen:    screen,
	}
	v.init(fps)

	return v, nil
}
func (v *video) init(fps float64) {
	go func() {
		ticker := time.NewTicker(time.Millisecond)
		defer ticker.Stop()

		defer close(v.buf)
		for {
			select {
			case <-v.ctx.Done():
				return
			case <-ticker.C:
				buf := make([]image.Image, int(fps+1)*5)
				n, err := v.src.Read(buf)
				if err != nil {
					return
				}

				if n > 0 {
					for _, b := range buf {
						v.buf <- b
					}
				}
			}
		}
	}()
}
func (v *video) Close() error {
	<-v.pauseCh

	v.cancel()
	v.src.Close()
	close(v.pauseCh)
	close(v.unPauseCh)
	return nil
}

func (v *video) PlayedTime() (time.Duration, error) {
	playedSize := (v.src.readIdx + 1) - int64(len(v.buf))
	d := float64(playedSize) / v.fps
	return time.ParseDuration(fmt.Sprintf("%vs", d))
}
func (v *video) Wait(d time.Duration) {
	if d == 0 {
		return
	}

	v.Pause()
	time.Sleep(d)
	v.Play()
}

func (v *video) Play() {
	if v.isPausing {
		v.isPausing = false
		v.unPauseCh <- true
		return
	}

	if v.isPlaying {
		return
	}
	v.isPlaying = true
	go func() {
		ticker := time.NewTicker(v.duration)
		defer ticker.Stop()
		for {
			select {
			case <-v.ctx.Done():
				v.isPlaying = false
				v.isPausing = false
				return
			default:
				if v.isPausing {
					v.pauseCh <- true
					<-v.unPauseCh
				}

				res, ok := <-v.buf
				if !ok {
					v.isPlaying = false
					v.isPausing = false
					return
				}
				v.screen.Image = res
				v.screen.Refresh()
				<-ticker.C
			}
		}
	}()
}
func (v *video) Pause() {
	if !v.isPlaying {
		return
	}
	v.isPausing = true
	<-v.pauseCh
}
func (v *video) IsPlaying() bool {
	return v.isPlaying
}
func (v *video) IsPausing() bool {
	return v.isPausing
}
func (v *video) Seek(offset int64, whence int) (int64, error) {
	v.Pause()

	n, err := v.src.Seek(offset, whence)
	if err != nil {
		return -1, err
	}

	return n, nil
}

type videoPlayer struct {
	ctx    context.Context
	cancel func()

	is  ipfs.Ipfs
	cid string

	audio   iAudioPlayer
	video   iVideoPlayer
	timeBar gutil.ITimeBar

	screen *canvas.Image
}

func NewVideoPlayer(cid string, is ipfs.Ipfs) (*videoPlayer, error) {
	vp := &videoPlayer{
		cid: cid,
		is:  is,
	}
	if err := vp.init(); err != nil {
		return nil, err
	}

	return vp, nil
}

func (vp *videoPlayer) init() error {
	img := image.NewGray(image.Rect(0, 0, store.VideoW, store.VideoH))
	vp.screen = canvas.NewImageFromImage(img)
	vp.screen.FillMode = canvas.ImageFillContain
	vp.screen.Refresh()

	m, err := vp.is.Get(vp.cid)
	if err != nil {
		return err
	}
	pbVideo := &pb.Video{}
	if err := proto.Unmarshal(m, pbVideo); err != nil {
		return err
	}

	aDec, err := newAudioDecoder(pbVideo.GetAudio())
	if err != nil {
		return err
	}
	audio, err := newOtoPlayer(aDec)
	if err != nil {
		return err
	}
	vp.audio = audio

	vDec, err := newVideoDecoder(pbVideo.GetVideo(), pbVideo.GetFrameCount())
	if err != nil {
		return err
	}
	video, err := newVideo(vDec, pbVideo.GetFrameRate(), vp.screen)
	if err != nil {
		return err
	}
	vp.video = video

	vp.timeBar = gutil.NewTimeBar(pbVideo.GetSecond())

	ctx, cancel := context.WithCancel(context.Background())
	vp.ctx = ctx
	vp.cancel = cancel

	return nil
}
func (vp *videoPlayer) Close() error {
	vp.audio.Pause()
	vp.video.Pause()
	vp.timeBar.Pause()

	vp.cancel()
	time.Sleep(time.Second)

	err1 := vp.video.Close()
	err2 := vp.audio.Close()
	err3 := vp.timeBar.Close()
	if err1 != nil {
		return err1
	} else if err2 != nil {
		return err2
	} else {
		return err3
	}
}

func (vp *videoPlayer) SyncTime() {
	go func() {
		ticker := time.NewTicker(time.Second * 30)
		defer ticker.Stop()

		for {
			select {
			case <-vp.ctx.Done():
				return
			case <-ticker.C:
				aPlaying := vp.audio.IsPlaying()
				vPlaying := vp.video.IsPlaying()
				tPlaying := vp.timeBar.IsPlaying()
				if !aPlaying || !vPlaying || !tPlaying {
					continue
				}
				aPausing := vp.audio.IsPausing()
				vPausing := vp.video.IsPausing()
				tPausing := vp.timeBar.IsPausing()
				if aPausing || vPausing || tPausing {
					continue
				}

				aTime, err1 := vp.audio.PlayedTime()
				vTime, err2 := vp.video.PlayedTime()
				tTime, err3 := vp.timeBar.PlayedTime()
				if err1 != nil || err2 != nil || err3 != nil {
					return
				}
				min := gutil.MinTime(aTime, vTime, tTime)
				vp.audio.Wait(aTime - min)
				vp.video.Wait(vTime - min)
				vp.timeBar.Wait(tTime - min)

				fmt.Println("min", min)
				fmt.Println("audio", aTime)
				fmt.Println("video", vTime)
				fmt.Println("tbar", tTime)
			}
		}

	}()
}
func (vp *videoPlayer) Render() (fyne.CanvasObject, gutil.Closer) {
	var playBtn *widget.Button
	playBtn = widget.NewButtonWithIcon("", theme.MediaPlayIcon(), func() {
		if vp.audio.IsPlaying() {
			vp.audio.Pause()
			vp.video.Pause()
			vp.timeBar.Pause()
			playBtn.SetIcon(theme.MediaPlayIcon())
		} else {
			vp.audio.Play()
			vp.video.Play()
			vp.timeBar.Play()
			playBtn.SetIcon(theme.MediaPauseIcon())
		}
	})

	resetBtn := widget.NewButtonWithIcon("", theme.MediaSkipPreviousIcon(), func() {
		vp.audio.Pause()
		vp.video.Pause()
		vp.timeBar.Pause()

		vp.Close()
		vp.init()
		vp.SyncTime()
		playBtn.SetIcon(theme.MediaPlayIcon())
	})

	vp.SyncTime()
	screen := container.NewGridWrap(fyne.NewSize(800, 450), vp.screen)

	btns := container.NewHBox(playBtn, resetBtn)
	obj := container.NewBorder(screen, vp.timeBar.Render(btns), nil, nil)
	return obj, vp.Close
}
