package gui

import (
	"context"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/speaker"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	store "github.com/pilinsin/lontan/store"
	pb "github.com/pilinsin/lontan/store/pb"
	ipfs "github.com/pilinsin/p2p-verse/ipfs"
	proto "google.golang.org/protobuf/proto"
)

const (
	channelCount     = 2
	bitDepth         = 8
	sampleBufferSize = 32 * channelCount * bitDepth * 1024
)

type audioPlayer struct {
	is         ipfs.Ipfs
	chunkCids  []string
	duration   time.Duration
	sampleRate int

	ctx    context.Context
	cancel func()

	chunk      *audioChunk
	nextChunk  *audioChunk
	chunkIndex int

	timeBar *widget.Slider

	playing bool
	paused  bool
}

func NewAudioPlayer(cid string, is ipfs.Ipfs) (*audioPlayer, error) {
	m, err := is.Get(cid)
	if err != nil {
		return nil, err
	}
	pbAudio := &pb.Audio{}
	if err := proto.Unmarshal(m, pbAudio); err != nil {
		return nil, err
	}

	a := store.DecodeAudio(pbAudio)
	ap := &audioPlayer{
		is:         is,
		chunkCids:  a.ChunkCids(),
		duration:   a.Duration(),
		sampleRate: a.SampleRate(),
	}
	if err := ap.init(); err != nil {
		return nil, err
	}

	return ap, nil
}

func (ap *audioPlayer) speakerInit() error {
	ssr := beep.SampleRate(ap.sampleRate)
	return speaker.Init(ssr, ssr.N(time.Millisecond*100))
}
func (ap *audioPlayer) init() error {
	if err := ap.speakerInit(); err != nil {
		return err
	}

	ap.ctx, ap.cancel = context.WithCancel(context.Background())
	ch, err := ap.loadChunk(0)
	if err != nil {
		return err
	}
	ch2, err := ap.loadChunk(1)
	ap.chunk = ch
	ap.nextChunk = ch2
	if err != nil {
		return err
	}

	ap.timeBar = &widget.Slider{
		Max:  ap.duration.Seconds(),
		Step: time.Second.Seconds() * 10,
	}
	onChanged := func(val float64) {
		ratio := val / ap.timeBar.Max
		idx := int(ratio * float64(len(ap.chunkCids)))
		if idx != ap.chunkIndex {
			ap.Clear()
			if err := ap.speakerInit(); err != nil {
				return
			}
			ap.ctx, ap.cancel = context.WithCancel(context.Background())
			if err := ap.LoadChunk(idx); err != nil {
				return
			}
			ap.Play()
		}
	}
	ap.timeBar.OnChanged = onChanged
	ap.timeBar.ExtendBaseWidget(ap.timeBar)

	ap.paused = true
	return nil
}
func (ap *audioPlayer) Clear() {
	//stop player for re-init
	ap.paused = true
	ap.playing = false
	if ap.cancel != nil {
		ap.cancel()
	}
}
func (ap *audioPlayer) Close() {
	ap.paused = true
	ap.playing = false
	if ap.cancel != nil {
		ap.cancel()
	}
	speaker.Clear()
	speaker.Close()
	ap.is = nil
}

type audioChunk struct {
	sampleBuffer <-chan [2]float64
}

func (ap *audioPlayer) loadChunk(idx int) (*audioChunk, error) {
	if idx >= len(ap.chunkCids) {
		return nil, nil
	}

	mpbc, err := ap.is.Get(ap.chunkCids[idx])
	if err != nil {
		return nil, err
	}
	pbc := &pb.ChunkAudio{}
	if err := proto.Unmarshal(mpbc, pbc); err != nil {
		return nil, err
	}
	chunk := store.DecodeChunkAudio(pbc)

	sampleBuffer := make(chan [2]float64, sampleBufferSize)
	go func() {
		defer close(sampleBuffer)
		for _, sample := range chunk.Samples() {
			select {
			case <-ap.ctx.Done():
				return
			default:
				sampleBuffer <- sample
			}
		}
	}()

	return &audioChunk{sampleBuffer}, nil
}
func (ap *audioPlayer) LoadChunk(idx int) error {
	if idx == ap.chunkIndex {
		return nil
	}
	if idx == ap.chunkIndex+1 {
		*ap.chunk = *ap.nextChunk
		nc, err := ap.loadChunk(idx + 1)
		ap.nextChunk = nc
		ap.chunkIndex++
		return err
	}

	ch, err := ap.loadChunk(idx)
	if err != nil {
		return err
	}
	ap.chunk = ch

	nch, err := ap.loadChunk(idx + 1)
	ap.nextChunk = nch
	ap.chunkIndex = idx
	return err
}

func (ap *audioPlayer) newBeepStreamer() beep.Streamer {
	return beep.StreamerFunc(func(samples [][2]float64) (int, bool) {
		numRead := 0
		if ap.chunk == nil {
			return numRead, false
		}

		for i := 0; i < len(samples); i++ {
			select {
			case <-ap.ctx.Done():
				break
			default:
			}

			if ap.paused {
				time.Sleep(time.Millisecond * 8)
				i--
				continue
			}

			sample, ok := <-ap.chunk.sampleBuffer
			if !ok {
				numRead = i + 1
				ap.cancel()
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

func (ap *audioPlayer) updateTimeBar() {
	ticker := time.NewTicker(time.Second)
	go func() {
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				if !ap.paused {
					ap.timeBar.Value += time.Second.Seconds()
					ap.timeBar.Refresh()
				}
			case <-ap.ctx.Done():
				return
			}
		}
	}()
}
func (ap *audioPlayer) Play() {
	ap.paused = false
	ap.playing = true
	speaker.Play(ap.newBeepStreamer())
	ap.updateTimeBar()
	go func() {
		for {
			select {
			case <-ap.ctx.Done():
				if ap.playing { //when sampleBuffer is exhausted, ctx.Done() && playing == true
					ap.ctx, ap.cancel = context.WithCancel(context.Background())
					if err := ap.LoadChunk(ap.chunkIndex + 1); err != nil {
						break
					}

					speaker.Play(ap.newBeepStreamer())
					ap.updateTimeBar()
				}
			}
		}
		ap.paused = true
		ap.playing = false
	}()
}

func (ap *audioPlayer) Pause() {
	ap.paused = !ap.paused
}

func (ap *audioPlayer) Render() fyne.CanvasObject {
	var playBtn *widget.Button
	playBtn = widget.NewButtonWithIcon("", theme.MediaPlayIcon(), func() {
		if !ap.playing {
			ap.Play()
		} else {
			ap.Pause()
		}

		if ap.paused {
			playBtn.SetIcon(theme.MediaPlayIcon())
		} else {
			playBtn.SetIcon(theme.MediaPauseIcon())
		}
	})

	return container.NewBorder(nil, nil, playBtn, nil, ap.timeBar)
}
