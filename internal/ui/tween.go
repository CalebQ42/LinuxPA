package ui

import (
	"math"
	"time"
)

type Tween struct {
	startVal  int32
	startTime time.Time
	dif       int32
	dur       time.Duration
}

func NewTween(startVal, endVal int32, duration time.Duration) *Tween {
	return &Tween{
		startVal:  startVal,
		dif:       endVal - startVal,
		dur:       duration,
		startTime: time.Now().Round(time.Millisecond),
	}
}

func (t *Tween) CurVal() int32 {
	since := time.Now().Round(time.Millisecond).Sub(t.startTime)
	if since <= 0 {
		return t.startVal
	}
	if since >= t.dur {
		return t.startVal + t.dif
	}
	m := float32(since.Milliseconds()) / float32(t.dur.Milliseconds())
	m = -1 * float32(math.Log10(-.9*float64(m)+1))
	return t.startVal + int32(math.Round(float64(float32(t.dif)*m)))
}

func (t *Tween) Ended() bool {
	return time.Now().Round(time.Millisecond).Sub(t.startTime) >= t.dur
}
