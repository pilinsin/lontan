package gui

/*
import (
	"context"
	"fmt"
	"image"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/speaker"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	store "github.com/pilinsin/lontan/store"
	pb "github.com/pilinsin/lontan/store/pb"
	ipfs "github.com/pilinsin/p2p-verse/ipfs"
	proto "google.golang.org/protobuf/proto"
)

const (
	width           = 1280
	height          = 720
	frameBufferSize = 1024
)

type videoPlayer struct {
	is         ipfs.Ipfs
	chunkCids  []string
	duration   time.Duration
	frameRate  int
	sampleRate int

	ctx    context.Context
	cancel func()

	chunk      *chunk
	nextChunk  *chunk
	chunkIndex int

	ticker *time.Ticker

	sprite  *image.RGBA
	screen  *canvas.Image
	timeBar *widget.Slider

	playing bool
	paused  bool
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

	v := store.DecodeVideo(pbVideo)
	vp := &videoPlayer{
		is:         is,
		chunkCids:  v.ChunkCids(),
		duration:   v.Duration(),
		frameRate:  v.FrameRate(),
		sampleRate: v.SampleRate(),
	}
	if err := vp.init(); err != nil {
		return nil, err
	}

	return vp, nil
}

func (vp *videoPlayer) speakerInit() error {
	ssr := beep.SampleRate(vp.sampleRate)
	return speaker.Init(ssr, ssr.N(time.Millisecond*100))
}
func (vp *videoPlayer) init() error {
	if err := vp.speakerInit(); err != nil {
		return err
	}

	vp.ctx, vp.cancel = context.WithCancel(context.Background())
	ch, err := vp.loadChunk(0)
	if err != nil {
		return err
	}
	ch2, err := vp.loadChunk(1)
	vp.chunk = ch
	vp.nextChunk = ch2
	if err != nil {
		return err
	}

	tickDur, err := time.ParseDuration(fmt.Sprintf("%fs", 1.0/float64(vp.frameRate)))
	if err != nil {
		return err
	}
	vp.ticker = time.NewTicker(tickDur)

	vp.sprite = image.NewRGBA(image.Rect(0, 0, width, height))
	vp.screen = canvas.NewImageFromImage(vp.sprite)
	vp.screen.ScaleMode = canvas.ImageScaleFastest

	vp.timeBar = &widget.Slider{
		Max:  vp.duration.Seconds(),
		Step: time.Second.Seconds() * 10,
	}
	onChanged := func(val float64) {
		ratio := val / vp.timeBar.Max
		idx := int(ratio * float64(len(vp.chunkCids)))
		if idx != vp.chunkIndex {
			vp.Clear()
			if err := vp.speakerInit(); err != nil {
				return
			}
			vp.ctx, vp.cancel = context.WithCancel(context.Background())
			if err := vp.LoadChunk(idx); err != nil {
				return
			}
			vp.Play()
		}
	}
	vp.timeBar.OnChanged = onChanged
	vp.timeBar.ExtendBaseWidget(vp.timeBar)

	vp.paused = true
	return nil
}
func (vp *videoPlayer) Clear() {
	//stop player for re-init
	vp.paused = true
	vp.playing = false
	if vp.cancel != nil {
		vp.cancel()
	}
}
func (vp *videoPlayer) Close() {
	vp.paused = true
	vp.playing = false
	if vp.cancel != nil {
		vp.cancel()
	}
	speaker.Clear()
	speaker.Close()
	if vp.ticker != nil {
		vp.ticker.Stop()
	}
	vp.is = nil
}

type chunk struct {
	frameBuffer  <-chan *image.RGBA
	sampleBuffer <-chan [2]float64
}

func (vp *videoPlayer) loadChunk(idx int) (*chunk, error) {
	if idx >= len(vp.chunkCids) {
		return nil, nil
	}

	mpbc, err := vp.is.Get(vp.chunkCids[idx])
	if err != nil {
		return nil, err
	}
	pbc := &pb.ChunkVideo{}
	if err := proto.Unmarshal(mpbc, pbc); err != nil {
		return nil, err
	}
	chk := store.DecodeChunkVideo(pbc)

	frameBuffer := make(chan *image.RGBA, frameBufferSize)
	sampleBuffer := make(chan [2]float64, sampleBufferSize)
	go func() {
		defer close(frameBuffer)
		for _, frame := range chk.Frames() {
			select {
			case <-vp.ctx.Done():
				return
			default:
				frameBuffer <- frame
			}
		}
	}()
	go func() {
		defer close(sampleBuffer)
		for _, sample := range chk.Samples() {
			select {
			case <-vp.ctx.Done():
				return
			default:
				sampleBuffer <- sample
			}
		}
	}()

	return &chunk{frameBuffer, sampleBuffer}, nil
}
func (vp *videoPlayer) LoadChunk(idx int) error {
	if idx == vp.chunkIndex {
		return nil
	}
	if idx == vp.chunkIndex+1 {
		*vp.chunk = *vp.nextChunk
		nc, err := vp.loadChunk(idx + 1)
		vp.nextChunk = nc
		vp.chunkIndex++
		return err
	}

	ch, err := vp.loadChunk(idx)
	if err != nil {
		return err
	}
	vp.chunk = ch

	nch, err := vp.loadChunk(idx + 1)
	vp.nextChunk = nch
	vp.chunkIndex = idx
	return err
}

func (vp *videoPlayer) newBeepStreamer() beep.Streamer {
	return beep.StreamerFunc(func(samples [][2]float64) (int, bool) {
		numRead := 0
		if vp.chunk == nil {
			return numRead, false
		}

		for i := 0; i < len(samples); i++ {
			select {
			case <-vp.ctx.Done():
				break
			default:
			}

			if vp.paused {
				time.Sleep(time.Millisecond * 8)
				i--
				continue
			}

			sample, ok := <-vp.chunk.sampleBuffer
			if !ok {
				numRead = i + 1
				break
			}
			samples[i] = sample
			numRead++
		}

		if numRead < len(samples) {
			return numRead, false
		}
		return numRead, true
	})
}

func (vp *videoPlayer) updateSprite() error {
	if vp.chunk == nil {
		return fmt.Errorf("chunk is nil")
	}
	for {
		if vp.paused {
			time.Sleep(time.Millisecond * 8)
			continue
		}
		select {
		case <-vp.ticker.C:
			frame, ok := <-vp.chunk.frameBuffer
			if ok {
				vp.sprite.Pix = frame.Pix
				vp.screen.Refresh()
			} else {
				return nil
			}
		case <-vp.ctx.Done():
			return vp.ctx.Err()
		default:
		}
	}
}
func (vp *videoPlayer) updateTimeBar() {
	ticker := time.NewTicker(time.Second)
	go func() {
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				if !vp.paused {
					vp.timeBar.Value += time.Second.Seconds()
					vp.timeBar.Refresh()
				}
			case <-vp.ctx.Done():
				return
			}
		}
	}()
}
func (vp *videoPlayer) Play() {
	vp.paused = false
	vp.playing = true
	go func() {
		for {
			speaker.Play(vp.newBeepStreamer())
			vp.updateTimeBar()
			if err := vp.updateSprite(); err != nil {
				break
			} //vp.cancel() from outside
			vp.cancel()
			vp.ctx, vp.cancel = context.WithCancel(context.Background())
			if err := vp.LoadChunk(vp.chunkIndex + 1); err != nil {
				break
			}
		}
		vp.paused = true
		vp.playing = false
	}()
}

func (vp *videoPlayer) Pause() {
	vp.paused = !vp.paused
}

func (vp *videoPlayer) Render() fyne.CanvasObject {
	var playBtn *widget.Button
	playBtn = widget.NewButtonWithIcon("", theme.MediaPlayIcon(), func() {
		if !vp.playing {
			vp.Play()
		} else {
			vp.Pause()
		}

		if vp.paused {
			playBtn.SetIcon(theme.MediaPlayIcon())
		} else {
			playBtn.SetIcon(theme.MediaPauseIcon())
		}
	})

	info := container.NewBorder(nil, nil, playBtn, nil, vp.timeBar)
	return container.NewBorder(nil, info, nil, nil, vp.screen)
}
*/
