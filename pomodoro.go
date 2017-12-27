package main

import (
	"flag"
	"fmt"
	"time"

	ui "github.com/gizak/termui"
)

var workPeriod = flag.Duration("work", 10*time.Second, "work period")
var shortRestPeriod = flag.Duration("rest", 5*time.Second, "rest period")
var longRestPeriod = flag.Duration("long_rest", 6*time.Second, "long rest period")

type timerState int

const (
	started timerState = iota
	paused
) // controls whether timer is running

type pomodoroState int

const (
	work pomodoroState = iota
	shortRest
	longRest
) // controls duration of timer

type pomodoro struct {
	gauge       *ui.Gauge
	t           time.Duration
	maxDuration time.Duration
	tState      timerState
	pStateIdx   int // points into stateSeq
	stateSeq    [4]pomodoroState
}

var timeMap = map[pomodoroState]time.Duration{
	work:      *workPeriod,
	shortRest: *shortRestPeriod,
	longRest:  *longRestPeriod,
}

func (p *pomodoro) pStateTransition() { // mutates pStateIdx in place
	if p.pStateIdx < cap(p.stateSeq)-1 {
		p.pStateIdx++
	} else {
		p.pStateIdx = 0
	}
}

func makePomodoro() *pomodoro {
	g := ui.NewGauge()
	g.Percent = 100 // updated periodically
	g.Width = 50
	g.Height = 3
	g.Y = 11
	g.BorderLabel = "Duration || Elapsed: || pState: " // updated periodically
	g.BarColor = ui.ColorRed
	g.BorderFg = ui.ColorWhite
	g.BorderLabelFg = ui.ColorCyan
	var p = pomodoro{g, 0, *workPeriod, paused, 0, [4]pomodoroState{work, shortRest, work, longRest}}

	g.Handle("/timer/1s", func(e ui.Event) {
		if p.tState != started {
			return
		}
		_ = e.Data.(ui.EvtTimer)
		p.t += time.Duration(1) * time.Second
		g.Percent = int(float64(p.maxDuration-p.t) / float64(p.maxDuration) * 100)
		g.BorderLabel = fmt.Sprintf("Duration %s || Elapsed: %s || pState: %d", p.maxDuration, p.t.String(), p.stateSeq[p.pStateIdx])
		if p.t >= p.maxDuration {
			p.pStateTransition()
			p.tState = paused
			p.t = 0
			p.maxDuration = timeMap[p.stateSeq[p.pStateIdx]]
			return
		}

	})
	// g.Handle("/sys/kbd/s", func(ui.Event) {
	// 	p.tState = started
	// })
	// g.Handle("/sys/kbd/p", func(ui.Event) {
	// 	p.tState = paused
	// })
	// g.Handle("/sys/kbd/r", func(ui.Event) {
	// 	p.t = 0
	// 	g.Percent = int(float64(p.maxDuration-p.t) / float64(p.maxDuration) * 100)
	// 	g.BorderLabel = fmt.Sprintf("Time Elapsed: %s", p.t.String())
	// })
	return &p
}
