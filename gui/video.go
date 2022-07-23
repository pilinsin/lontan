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

func videoDataToFile(vData []byte) (string, error) {
	f, err := os.CreateTemp("", ".*_tmp")
	if err != nil {
		return "", err
	}
	defer f.Close()
	if _, err := f.Write(vData); err != nil {
		return "", err
	}

	return f.Name(), nil
}

type videoDecoder struct {
	srcName    string
	frameCount int64
	readIdx    int64

	vc *gocv.VideoCapture
}

func newVideoDecoder(vTmpName string, frameCount int64) (*videoDecoder, error) {
	dec := &videoDecoder{
		srcName:    vTmpName,
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
	bufSize   int
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

	bufSize := int(fps+1) * 2
	v := &video{
		ctx:       ctx,
		cancel:    cancel,
		src:       vr,
		buf:       make(chan image.Image, bufSize),
		bufSize:   bufSize,
		duration:  dur,
		fps:       fps,
		pauseCh:   make(chan bool, 1),
		unPauseCh: make(chan bool, 1),
		screen:    screen,
	}
	v.init()

	return v, nil
}
func (v *video) init() {
	tmpBufSize := int((v.fps + 1) * 0.3)
	if len(v.buf) >= v.bufSize-tmpBufSize {
		return
	}

	buf := make([]image.Image, tmpBufSize)
	n, err := v.src.Read(buf)
	if err != nil {
		return
	}

	for _, b := range buf[:n] {
		v.buf <- b
	}
}

func (v *video) Close() error {
	if v.isPlaying && !v.isPausing {
		v.Pause()
	}
	v.cancel()
	v.unPauseCh <- true
	time.Sleep(time.Millisecond * 10)

	v.src.Close()
	close(v.pauseCh)
	close(v.unPauseCh)
	close(v.buf)
	return nil
}
func (v *video) Reset() {
	if v.isPlaying && !v.isPausing {
		v.Pause()
	}

	v.Seek(0, io.SeekStart)
}

func (v *video) PlayedTime() (time.Duration, error) {
	playedSize := (v.src.readIdx + 1) - int64(len(v.buf))
	d := float64(playedSize) / v.fps
	return time.ParseDuration(fmt.Sprintf("%vs", d))
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

				v.init()
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
	if !v.isPlaying || v.isPausing {
		return
	}
	v.isPausing = true
	<-v.pauseCh
}
func (v *video) Wait(d time.Duration) {
	if d == 0 {
		return
	}

	v.Pause()
	time.Sleep(d)
	v.Play()
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
		return 0, err
	}

	close(v.buf)
	v.buf = make(chan image.Image, v.bufSize)
	v.init()

	return n, nil
}

type videoPlayer struct {
	ctx    context.Context
	cancel func()

	aData      []byte
	vFileName  string
	frameRate  float64
	frameCount int64
	second     float64

	audio   iAudioPlayer
	video   iVideoPlayer
	timeBar gutil.ITimeBar

	isSyncing bool

	screen *canvas.Image
}

func NewVideoPlayer(cid string, is ipfs.Ipfs) (*videoPlayer, error) {
	m, err := is.Get(cid)
	if err != nil {
		return nil, err
	}
	pbVideo := &pb.Video{}
	if err := proto.Unmarshal(m, pbVideo); err != nil {
		return nil, err
	}

	vData, err := videoDataToFile(pbVideo.GetVideo())
	if err != nil {
		return nil, err
	}
	vp := &videoPlayer{
		aData:      pbVideo.GetAudio(),
		vFileName:  vData,
		frameRate:  pbVideo.GetFrameRate(),
		frameCount: pbVideo.GetFrameCount(),
		second:     pbVideo.GetSecond(),
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

	aDec, err := newAudioDecoder(vp.aData)
	if err != nil {
		return err
	}
	vp.audio, err = newOtoPlayer(aDec)
	if err != nil {
		return err
	}

	vDec, err := newVideoDecoder(vp.vFileName, vp.frameCount)
	if err != nil {
		return err
	}
	vp.video, err = newVideo(vDec, vp.frameRate, vp.screen)
	if err != nil {
		return err
	}

	vp.timeBar = gutil.NewTimeBar(vp.second, nil)

	ctx, cancel := context.WithCancel(context.Background())
	vp.ctx = ctx
	vp.cancel = cancel

	vp.isSyncing = false

	return nil
}
func (vp *videoPlayer) Close() error {
	if vp.IsPlaying() && !vp.IsPausing() {
		vp.Pause()
	}

	vp.cancel()

	err1 := vp.audio.Close()
	err2 := vp.timeBar.Close()
	err3 := vp.video.Close()

	os.Remove(vp.vFileName)

	if err1 != nil {
		return err1
	} else if err2 != nil {
		return err2
	} else {
		return err3
	}
}
func (vp *videoPlayer) Reset() {
	if vp.IsPlaying() && !vp.IsPausing() {
		vp.Pause()
	}

	vp.timeBar.Reset()
	vp.video.Reset()
	vp.audio.Reset()

	vp.screen.Image = image.NewGray(image.Rect(0, 0, store.VideoW, store.VideoH))
	vp.screen.Refresh()
}

func (vp *videoPlayer) Play() {
	vp.timeBar.Play()
	vp.video.Play()
	vp.audio.Play()
}
func (vp *videoPlayer) Pause() {
	vp.timeBar.Pause()
	vp.video.Pause()
	vp.audio.Pause()
}

func (vp *videoPlayer) IsPlaying() bool {
	return vp.video.IsPlaying() && vp.audio.IsPlaying() && vp.timeBar.IsPlaying()
}
func (vp *videoPlayer) IsPausing() bool {
	return vp.video.IsPausing() && vp.audio.IsPausing() && vp.timeBar.IsPausing()
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
				if !vp.IsPlaying() || vp.IsPausing() {
					continue
				}
				vp.isSyncing = true

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

				vp.isSyncing = false
			}
		}

	}()
}
func (vp *videoPlayer) Render() (fyne.CanvasObject, gutil.Closer) {
	var btns fyne.CanvasObject
	var screen fyne.CanvasObject
	var timeBar fyne.CanvasObject
	var obj fyne.CanvasObject

	var playBtn *widget.Button
	playBtn = widget.NewButtonWithIcon("", theme.MediaPlayIcon(), func() {
		if vp.isSyncing {
			return
		}
		if vp.IsPlaying() && !vp.IsPausing() {
			vp.Pause()
			playBtn.SetIcon(theme.MediaPlayIcon())
		} else {
			vp.Play()
			playBtn.SetIcon(theme.MediaPauseIcon())
		}
	})

	resetBtn := widget.NewButtonWithIcon("", theme.MediaSkipPreviousIcon(), func() {
		if vp.isSyncing {
			return
		}
		vp.Reset()

		screen.(*fyne.Container).Objects[0] = vp.screen
		screen.Refresh()
		timeBar.(*fyne.Container).Objects[0] = vp.timeBar.Render()
		timeBar.Refresh()
		obj.Refresh()
		playBtn.SetIcon(theme.MediaPlayIcon())
	})

	vp.SyncTime()
	screen = container.NewGridWrap(fyne.NewSize(800, 450), vp.screen)
	btns = container.NewHBox(playBtn, resetBtn)
	timeBar = container.NewBorder(nil, nil, btns, nil, vp.timeBar.Render())
	obj = container.NewBorder(screen, timeBar, nil, nil)
	return obj, vp.Close
}
