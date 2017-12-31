package main

import (
	"fmt"
	"time"

	ui "github.com/gizak/termui"
)

// pomodoro implements the cli component that displays and counts down for
// the current pomodoro.

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

func (p *pomodoro) pStateTransition(sd statsDisplay) { // mutates pStateIdx in place
	debugDisplay.updateText([]string{nameStateMap[p.stateSeq[p.pStateIdx]]})
	statsGlobal[p.stateSeq[p.pStateIdx]]++
	sd.refreshStatsDisplay()
	if p.pStateIdx >= cap(p.stateSeq)-1 {
		p.pStateIdx = 0
	} else {
		p.pStateIdx++
	}
}

func makePomodoro(sd statsDisplay) *pomodoro {
	g := ui.NewGauge()
	g.Percent = 100 // updated periodically
	g.Width = 50
	g.Height = 3
	g.Y = 7
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
			p.pStateTransition(sd)
			p.tState = paused
			p.t = 0
			p.maxDuration = *timeMap[p.stateSeq[p.pStateIdx]]
			return
		}

	})
	// TODO(iantay) create functions to call from main.go
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
