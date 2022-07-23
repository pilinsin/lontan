package guiutil

import (
	"context"
	"fmt"
	"io"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func MinTime(d time.Duration, ds ...time.Duration) time.Duration {
	min := d
	for _, t := range ds {
		if t < min {
			min = t
		}
	}
	return min
}

func durationToDisplay(d int) string {
	h := d / 3600
	d %= 3600
	m := d / 60
	s := d % 60

	if h == 0 {
		return fmt.Sprintf("%02d:%02d", m, s)
	}
	return fmt.Sprintf("%d:%02d:%02d", h, m, s)
}

type IPlayer interface {
	Play()
	Pause()
	Wait(time.Duration)
	IsPlaying() bool
	IsPausing() bool
	PlayedTime() (time.Duration, error)
	Reset()
	io.Seeker
	io.Closer
}
type ITimeBar interface {
	IPlayer
	Render() fyne.CanvasObject
}

type timeBar struct {
	slider *widget.Slider
	ctx    context.Context
	cancel func()

	timeLabel *widget.Label
	totalStr  string
	last      float64
	isPlaying bool
	isPausing bool
	pauseCh   chan bool
	unPauseCh chan bool
}

//offsetRate: [0, 1]
func NewTimeBar(total float64, seek func(offsetRate float64)) *timeBar {
	ctx, cancel := context.WithCancel(context.Background())
	slider := widget.NewSlider(0, total)
	totalTime := durationToDisplay(int(total))
	timeLabel := widget.NewLabel("00:00/" + totalTime)
	tb := &timeBar{
		slider:    slider,
		ctx:       ctx,
		cancel:    cancel,
		timeLabel: timeLabel,
		totalStr:  totalTime,
		last:      slider.Value,
		pauseCh:   make(chan bool, 1),
		unPauseCh: make(chan bool, 1),
	}
	tb.slider.OnChanged = func(v float64) {
		offsetRate := v / total
		if seek == nil {
			return
		}
		seek(offsetRate)
		tb.Seek(int64(v), io.SeekStart)
	}

	return tb
}
func (tb *timeBar) Close() error {
	if tb.isPlaying && !tb.isPausing {
		tb.Pause()
	}
	tb.cancel()
	tb.unPauseCh <- true
	time.Sleep(time.Millisecond * 10)

	close(tb.pauseCh)
	close(tb.unPauseCh)
	return nil
}
func (tb *timeBar) Reset() {
	//if tb.Seek is implemented, tb.Pause() & tb.Seek(0, io.SeekStart)
	if tb.isPlaying && !tb.isPausing {
		tb.Pause()
	}

	tb.last = 0
	tb.slider.Value = 0
	tb.timeLabel.SetText("00:00/" + tb.totalStr)
}

func (tb *timeBar) Play() {
	if tb.isPausing {
		tb.isPausing = false
		tb.unPauseCh <- true
		return
	}

	if tb.isPlaying {
		return
	}
	tb.isPlaying = true
	go func() {
		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-tb.ctx.Done():
				tb.isPlaying = false
				tb.isPausing = false
				return
			default:
				if tb.isPausing {
					tb.pauseCh <- true
					<-tb.unPauseCh
				}

				tb.slider.Refresh()
				nowTime := durationToDisplay(int(tb.slider.Value))
				tb.timeLabel.SetText(nowTime + "/" + tb.totalStr)

				if tb.slider.Value >= tb.slider.Max {
					tb.isPlaying = false
					tb.isPausing = false
					return
				}

				tb.last = tb.slider.Value
				tb.slider.Value += 1
				<-ticker.C
			}
		}
	}()
}
func (tb *timeBar) Pause() {
	if !tb.isPlaying || tb.isPausing {
		return
	}
	tb.isPausing = true
	<-tb.pauseCh
}
func (tb *timeBar) Wait(d time.Duration) {
	if d == 0 {
		return
	}

	tb.Pause()
	time.Sleep(d)
	tb.Play()
}

func (tb *timeBar) IsPlaying() bool {
	return tb.isPlaying
}
func (tb *timeBar) IsPausing() bool {
	return tb.isPausing
}
func (tb *timeBar) PlayedTime() (time.Duration, error) {
	return time.ParseDuration(fmt.Sprintf("%vs", tb.slider.Value))
}

func (tb *timeBar) Seek(offset int64, _ int) (int64, error) {
	tb.slider.Value = tb.last
	tb.slider.Refresh()
	return offset, nil
	/*
		if tb.isPlaying && !tb.isPausing{
			tb.Pause()
		}

		ofst := float64(offset)
		tb.last = ofst
		tb.slider.Value = ofst
		nowTime := durationToDisplay(int(slider.Value))
		timeLabel.SetText(nowTime + "/" + totalTime)
		return offset, nil
	*/
}

func (tb *timeBar) Render() fyne.CanvasObject {
	return container.NewBorder(nil, nil, nil, tb.timeLabel, tb.slider)
}
